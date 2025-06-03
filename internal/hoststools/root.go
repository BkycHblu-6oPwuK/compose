package hoststools

import (
	"bufio"
	"bytes"
	"docky/internal/certstools"
	"docky/internal/config"
	"docky/internal/symlinkstools"
	"docky/pkg/filetools"
	"docky/pkg/readertools"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf16"
)

func PushToLocalHosts(domain string) error {
	filePath := config.GetLocalHostsFilePath()
	if err := filetools.InitDirs(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("ошибка инициализации директорий: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return fmt.Errorf("ошибка при открытии файла %s: %w", filePath, err)
	}
	defer file.Close()

	entry := "127.0.0.1 " + domain
	if lineInFile(filePath, entry) {
		return nil
	}

	_, err = file.WriteString(entry + "\n")
	return err
}

func PushToHosts() error {
	hostsFile := config.GetLocalHostsFilePath()
	if exists, _ := filetools.FileIsExists(hostsFile); !exists {
		return fmt.Errorf("файл %s не найден", hostsFile)
	}

	isWSL, _, hostFile, hostFileWSL, err := detectHostsPath()
	if err != nil {
		return err
	}

	file, err := os.Open(hostsFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := strings.TrimSpace(scanner.Text())
		if entry == "" {
			continue
		}

		fmt.Printf("добавление записи: %s\n", entry)

		if isWSL {
			if lineInFile(hostFileWSL, entry) {
				fmt.Println("запись уже есть в Windows hosts")
				return nil
			}
			if err := addToWindowsHosts(hostFile, entry); err != nil {
				return fmt.Errorf("ошибка при добавлении в hosts Windows: %w", err)
			}
		} else {
			if !lineInFile("/etc/hosts", entry) {
				if err := addToLinuxHosts(entry); err != nil {
					return fmt.Errorf("ошибка добавления в /etc/hosts: %w", err)
				}
				fmt.Println("добавлено в /etc/hosts.")
			} else {
				fmt.Println("запись уже существует")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка сканирования файла: %w", err)
	}

	return nil
}

func CreateDomain() error {
	domain, err := prepareDomain()
	if err != nil {
		return err
	}

	if err := certstools.CreateCerts(domain, config.SitePathInContainer, false); err != nil {
		return err
	}

	return PushToHosts()
}

func CreateSite() error {
	domain, err := prepareDomain()
	if err != nil {
		return err
	}

	if err := pushToSymlinks(domain); err != nil {
		return err
	}

	dirPath := filepath.Join(config.GetSiteDirPath(), domain)
	if err := filetools.InitDirs(dirPath); err != nil {
		return fmt.Errorf("ошибка инициализации директории сайта: %w", err)
	}

	if err := certstools.CreateCerts(domain, filepath.Join(config.SitePathInContainer, domain), true); err != nil {
		return err
	}

	return PushToHosts()
}

func pushToSymlinks(domain string) error {
	return symlinkstools.PushToSymlinks(map[string]string{
		"bitrix": domain + "/bitrix",
		"local":  domain + "/local",
		"upload": domain + "/upload",
	})
}

func prepareDomain() (string, error) {
	domain := readDomain()
	if err := PushToLocalHosts(domain); err != nil {
		return "", err
	}
	return domain, nil
}

func readDomain() string {
	return readertools.ReadLine("Введите название сайта (доменное имя): ")
}

func detectHostsPath() (isWSL bool, drive, hostFile, hostFileWSL string, err error) {
	drive = "C"

	if _, err := exec.LookPath("powershell.exe"); err == nil {
		if content, err := os.ReadFile("/proc/version"); err == nil && strings.Contains(strings.ToLower(string(content)), "microsoft") {
			isWSL = true
		}
	}

	if isWSL {
		getPaths := func(d string) (string, string) {
			return fmt.Sprintf("%s:\\Windows\\System32\\drivers\\etc\\hosts", d),
				fmt.Sprintf("/mnt/%s/Windows/System32/drivers/etc/hosts", strings.ToLower(d))
		}

		hostFile, hostFileWSL = getPaths(drive)

		if _, err := os.Stat(hostFileWSL); os.IsNotExist(err) {
			out, err := exec.Command("powershell.exe", "-NoProfile", "-Command", "[System.Environment]::SystemDirectory.Substring(0,1)").Output()
			if err != nil {
				return false, "", "", "", fmt.Errorf("не удалось получить диск: %w", err)
			}
			drive = strings.TrimSpace(string(out))
			hostFile, hostFileWSL = getPaths(drive)

			if _, err := os.Stat(hostFileWSL); os.IsNotExist(err) {
				return false, "", "", "", errors.New("файл hosts не был найден в Windows")
			}
		}
	} else {
		if _, err := os.Stat("/etc/hosts"); os.IsNotExist(err) {
			return false, "", "", "", errors.New("/etc/hosts не найден")
		}
	}

	return isWSL, drive, hostFile, hostFileWSL, nil
}

func addToWindowsHosts(hostFile, entry string) error {
	if err := runAsAdminPowerShell(hostFile, entry); err != nil {
		return err
	}

	fmt.Println("добавлено в hosts Windows")
	return nil
}

func runAsAdminPowerShell(hostFile, entry string) error {
	script := fmt.Sprintf(`$path = '%s'; $entry = '%s'; if (-not (Get-Content -Path $path | Where-Object { $_.Trim() -eq $entry })) { Add-Content -Path $path -Value "`+"\n"+`$entry" }`,
		hostFile,
		entry,
	)

	encoded := base64.StdEncoding.EncodeToString(utf16LE(script))
	cmd := exec.Command("powershell.exe",
		"-NoProfile",
		"-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-NoProfile','-ExecutionPolicy','Bypass','-EncodedCommand','%s' -Verb RunAs`, encoded),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func utf16LE(s string) []byte {
	runes := utf16.Encode([]rune(s))
	buf := new(bytes.Buffer)
	for _, r := range runes {
		buf.WriteByte(byte(r))
		buf.WriteByte(byte(r >> 8))
	}
	return buf.Bytes()
}

func addToLinuxHosts(entry string) error {
	cmd := exec.Command("tee", "-a", "/etc/hosts")
	cmd.Stdin = bytes.NewBufferString(entry + "\n")
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lineInFile(path, needle string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == needle {
			return true
		}
	}
	return false
}
