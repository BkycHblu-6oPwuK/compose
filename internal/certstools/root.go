package certstools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
)

var (
	dockerCachePath = config.GetDockerFilesDirPathInCache()
	confPath        = config.GetConfFilesDirPath()
	partialsPath    = filepath.Join(config.GetScriptCacheDir(), "partials", config.GetCurFramework())

	nginxCachePath      = filepath.Join(dockerCachePath, composefiletools.Nginx)
	nginxCacheCertsPath = filepath.Join(nginxCachePath, "certs")
	nginxPath           = filepath.Join(confPath, composefiletools.Nginx)
	nginxPathConf       = filepath.Join(nginxPath, composefiletools.ConfDir)
	nginxPathSnippets   = filepath.Join(nginxPathConf, "snippets")
	appNginxPath        = filepath.Join(confPath, composefiletools.App, composefiletools.Nginx)
	nginxCertsPath      = filepath.Join(nginxPath, "certs")

	rootCrt = filepath.Join(nginxCacheCertsPath, "rootCA.crt")
	rootKey = filepath.Join(nginxCacheCertsPath, "rootCA.key")

	stubDomainCertsFile = filepath.Join(partialsPath, "certs_domain_conf.stub")
	stubNginxCertsFile  = filepath.Join(partialsPath, "certs_nginx_conf.stub")
	stubNginxFile       = filepath.Join(partialsPath, "nginx_conf.stub")
	stubNginxPartFile   = filepath.Join(partialsPath, "nginx_part_conf.stub")
	stubNginxAppFile    = filepath.Join(partialsPath, "nginx_app_conf.stub")
)

type certContext struct {
	Domain   string
	RootPath string
}

func CreateCerts(domainName, rootPath string, isCreateSite bool) error {
	ctx := certContext{
		Domain:   domainName,
		RootPath: rootPath,
	}

	nginxCertsDomainPath := filepath.Join(nginxCertsPath, domainName)
	if err := filetools.InitDirs(nginxCertsDomainPath, nginxPathSnippets, appNginxPath); err != nil {
		return err
	}

	paths := map[string]string{
		"domainKey":    filepath.Join(nginxCertsDomainPath, domainName+".key"),
		"domainCsr":    filepath.Join(nginxCertsDomainPath, domainName+".csr"),
		"domainExt":    filepath.Join(nginxCertsDomainPath, domainName+".ext"),
		"domainCrt":    filepath.Join(nginxCertsDomainPath, domainName+".crt"),
		"nginxKey":     filepath.Join(nginxCertsDomainPath, "nginx.key"),
		"nginxCsr":     filepath.Join(nginxCertsDomainPath, "nginx.csr"),
		"nginxExt":     filepath.Join(nginxCertsDomainPath, "nginx.ext"),
		"nginxCrt":     filepath.Join(nginxCertsDomainPath, "nginx.crt"),
		"nginxConf":    filepath.Join(nginxPathConf, domainName+".conf"),
		"nginxPart":    filepath.Join(nginxPathSnippets, domainName+".conf"),
		"nginxAppConf": filepath.Join(appNginxPath, domainName+".conf"),
	}

	if err := renderStub(stubDomainCertsFile, paths["domainExt"], ctx); err != nil {
		return err
	}
	if err := renderStub(stubNginxCertsFile, paths["nginxExt"], ctx); err != nil {
		return err
	}
	if err := renderStub(stubNginxPartFile, paths["nginxPart"], ctx); err != nil {
		return err
	}
	if err := renderStub(stubNginxFile, paths["nginxConf"], ctx); err != nil {
		return err
	}
	if err := renderStub(stubNginxAppFile, paths["nginxAppConf"], ctx); err != nil {
		return err
	}

	if err := execOpenSSL(paths["nginxKey"], paths["nginxCsr"], "nginx", paths["nginxCrt"], paths["nginxExt"]); err != nil {
		return err
	}
	if err := execOpenSSL(paths["domainKey"], paths["domainCsr"], domainName, paths["domainCrt"], paths["domainExt"]); err != nil {
		return err
	}

	err := createVolumes(domainName, isCreateSite)

	if err != nil {
		return err
	}

	return nil
}

func createVolumes(domainName string, isCreateSite bool) error {
	serviceNames := []string{
		composefiletools.App,
		composefiletools.Nginx,
	}
	composeFile, err := composefile.Load(config.GetDockerComposeFilePath())
	if err != nil {
		return fmt.Errorf("ошибка при загрузке docker-compose.yml: %v", err)
	}

	nginxVolumeExists := false
	nginxService, exists := composeFile.Services.Get(composefiletools.Nginx)
	if exists {
		for _, vol := range nginxService.Volumes {
			if strings.HasSuffix(vol, composefiletools.GetNginxConfPathInContainer()) {
				nginxVolumeExists = true
				break
			}
		}
	}

	certificateConfPath := composefiletools.GetCertificateConfVolumePath(domainName)
	domainNameConf := domainName + ".conf"
	volumes := map[string][]string{
		composefiletools.App: {
			composefiletools.GetAppNginxConfVolumePath(domainNameConf),
			certificateConfPath,
		},
		composefiletools.Nginx: {
			certificateConfPath,
		},
	}
	if !nginxVolumeExists {
		volumes[composefiletools.Nginx] = append(volumes[composefiletools.Nginx], composefiletools.GetNginxConfVolumePath(domainNameConf))
		volumes[composefiletools.Nginx] = append(volumes[composefiletools.Nginx], composefiletools.GetNginxSnippetsConfVolumePath(domainNameConf))
	}
	if isCreateSite {
		volumes[composefiletools.App] = append(volumes[composefiletools.App], composefiletools.GetSymlinksConfVolumePath())
	}
	return composefiletools.PublishVolumes(serviceNames, volumes, nil)
}

func renderStub(stubPath, outPath string, ctx certContext) error {
	content, err := os.ReadFile(stubPath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении шаблона: %w", err)
	}

	tmpl, err := template.New(filepath.Base(stubPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("ошибка при разборе шаблона: %w", err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer outFile.Close()

	if err := tmpl.Execute(outFile, ctx); err != nil {
		return fmt.Errorf("ошибка при применении шаблона: %w", err)
	}

	return nil
}

func execOpenSSL(keyPath, csrPath, cn, crtPath, extPath string) error {
	// openssl genrsa -out keyPath 2048
	if err := runOpenSsl("genrsa", "-out", keyPath, "2048"); err != nil {
		return fmt.Errorf("ошибка openssl genrsa: %w", err)
	}

	// openssl req -new -key keyPath -out csrPath -subj "/C=RU/ST=Omsk Oblast/L=Omsk/O=company/CN=..."
	subj := fmt.Sprintf("/C=RU/ST=Omsk Oblast/L=Omsk/O=company/CN=%s", cn)
	if err := runOpenSsl("req", "-new", "-key", keyPath, "-out", csrPath, "-subj", subj); err != nil {
		return fmt.Errorf("ошибка запроса openssl: %w", err)
	}

	// openssl x509 -req -in csrPath -CA rootCrt -CAkey rootKey -CAcreateserial -out crtPath -days 9999 -sha256 -extfile extPath
	if err := runOpenSsl("x509", "-req", "-in", csrPath, "-CA", rootCrt, "-CAkey", rootKey, "-CAcreateserial", "-out", crtPath, "-days", "9999", "-sha256", "-extfile", extPath); err != nil {
		return fmt.Errorf("ошибка openssl x509: %w", err)
	}

	return nil
}

func runOpenSsl(args ...string) error {
	c := exec.Command("openssl", args...)
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
