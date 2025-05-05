package hosts

import (
	"bufio"
	"bytes"
	"docky/config"
	"docky/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func PushToLocalHosts(domain string) error {
	var file *os.File
	var err error
	filePath := config.GetLocalHostsFilePath()

	if !utils.FileIsExists(filePath) {
		file, err = os.Create(filePath)
	} else {
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("127.0.0.1 " + domain + "\n")
	return err
}

func PushToHosts() error {
	hostsFile := config.GetLocalHostsFilePath()
	if !utils.FileIsExists(hostsFile) {
		return fmt.Errorf("файл %s не найден", hostsFile)
	}

	isWSL, _, hostFileWindows, _, err := detectWindowsHostsPath()
	if err != nil {
		return err
	}

	file, err := os.Open(hostsFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла hosts: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hostEntry := strings.TrimSpace(scanner.Text())
		if hostEntry == "" {
			continue
		}
		fmt.Printf("Добавление записи в hosts: %s\n", hostEntry)

		if isWSL {
			err := addToWindowsHosts(hostFileWindows, hostEntry)
			if err != nil {
				return fmt.Errorf("ошибка при добавлении в Windows hosts: %w", err)
			}
		} else {
			if !lineInFile("/etc/hosts", hostEntry) {
				err := addToLinuxHosts(hostEntry)
				if err != nil {
					return fmt.Errorf("ошибка добавления в /etc/hosts: %w", err)
				}
				fmt.Println("Запись добавлена в /etc/hosts на Linux.")
			} else {
				fmt.Println("Запись уже существует в /etc/hosts.")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка сканирования файла hosts: %w", err)
	}

	return nil
}

func detectWindowsHostsPath() (bool, string, string, string, error) {
	isWSL := false
	systemDrive := "C"
	var hostFileWindowsWSL, hostFileWindows string

	if _, err := exec.LookPath("powershell.exe"); err == nil {
		if content, err := os.ReadFile("/proc/version"); err == nil && strings.Contains(strings.ToLower(string(content)), "microsoft") {
			isWSL = true
		} else if out, err := exec.Command("systemd-detect-virt").Output(); err == nil && strings.TrimSpace(string(out)) == "wsl" {
			isWSL = true
		}
	}

	if isWSL {
		findHostsFile := func(drive string) bool {
			driveLower := strings.ToLower(drive)
			hostFileWindowsWSL = fmt.Sprintf("/mnt/%s/Windows/System32/drivers/etc/hosts", driveLower)
			hostFileWindows = fmt.Sprintf("%s:\\Windows\\System32\\drivers\\etc\\hosts", drive)
			_, err := os.Stat(hostFileWindowsWSL)
			return err == nil
		}

		if !findHostsFile(systemDrive) {
			out, err := exec.Command("powershell.exe", "-NoProfile", "-Command", "[System.Environment]::SystemDirectory.Substring(0,1)").Output()
			if err != nil {
				return false, "", "", "", fmt.Errorf("не удалось получить системный диск: %w", err)
			}
			systemDrive = strings.TrimSpace(string(out))
			systemDrive = strings.Trim(systemDrive, "\r\n")
			if systemDrive == "" || !findHostsFile(systemDrive) {
				return false, "", "", "", errors.New("не удалось найти файл hosts на Windows")
			}
		}
	} else {
		if _, err := os.Stat("/etc/hosts"); os.IsNotExist(err) {
			return false, "", "", "", errors.New("/etc/hosts не найден")
		}
	}

	return isWSL, systemDrive, hostFileWindows, hostFileWindowsWSL, nil
}

func addToWindowsHosts(hostFile, hostEntry string) error {
	fmt.Println(hostEntry)
	if err := runAsAdminPowerShell(hostFile, hostEntry); err != nil {
		return err
	}

	fmt.Println("Запись добавлена или уже существует в hosts Windows.")
	return nil
}

func runAsAdminPowerShell(hostFile, hostEntry string) error {
	script := fmt.Sprintf(`$path = '%s'; $entry = '%s'; if (-not (Get-Content -Path $path | Where-Object { $_.Trim() -eq $entry })) { Add-Content -Path $path -Value "`+"\n"+`$entry" }`,
		hostFile,
		hostEntry,
	)

	utf16le := utf16LE(script)
	encoded := base64.StdEncoding.EncodeToString(utf16le)

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
	utf16 := []uint16{}
	for _, r := range s {
		utf16 = append(utf16, uint16(r))
	}

	result := make([]byte, len(utf16)*2)
	for i, v := range utf16 {
		result[i*2] = byte(v)
		result[i*2+1] = byte(v >> 8)
	}
	return result
}

func addToLinuxHosts(entry string) error {
	cmd := exec.Command("sudo", "tee", "-a", "/etc/hosts")
	cmd.Stdin = bytes.NewBufferString(entry + "\n")
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lineInFile(filename, line string) bool {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false
	}
	lines := strings.Split(string(content), "\n")
	for _, l := range lines {
		if strings.Contains(l, line) {
			return true
		}
	}
	return false
}
