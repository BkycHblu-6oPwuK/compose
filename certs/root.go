package certs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"docky/config"
	"docky/yaml/helper"
)

var (
	dockerCachePath = config.GetDockerFilesDirPathInCache()
	dockerPath      = config.GetDockerFilesDirPath()
	partialsPath    = filepath.Join(config.GetScriptCacheDir(), "partials")

	nginxCachePath      = filepath.Join(dockerCachePath, helper.Nginx)
	nginxCacheCertsPath = filepath.Join(nginxCachePath, "certs")
	nginxPath           = filepath.Join(dockerPath, helper.Nginx)
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

func Create(domainName, rootPath string) error {
	ctx := certContext{
		Domain:   domainName,
		RootPath: rootPath,
	}

	nginxCertsDomainPath := filepath.Join(nginxCertsPath, domainName)
	if err := os.MkdirAll(nginxCertsDomainPath, 0755); err != nil {
		return fmt.Errorf("failed to create certs directory: %w", err)
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
		"nginxConf":    filepath.Join(nginxPath, "conf.d", domainName+".conf"),
		"nginxPart":    filepath.Join(nginxPath, "conf.d", "snippets", domainName+".conf"),
		"nginxAppConf": filepath.Join(dockerPath, helper.App, helper.Nginx, domainName+".conf"),
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

	return nil
}

func renderStub(stubPath, outPath string, ctx certContext) error {
	content, err := os.ReadFile(stubPath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении шаблона: %w", err)
	}

	tmpl, err := template.New(filepath.Base(stubPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("ошибка при парсинге шаблона: %w", err)
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
	if err := run("openssl", "genrsa", "-out", keyPath, "2048"); err != nil {
		return fmt.Errorf("openssl genrsa error: %w", err)
	}

	// openssl req -new -key keyPath -out csrPath -subj "/C=RU/ST=Omsk Oblast/L=Omsk/O=company/CN=..."
	subj := fmt.Sprintf("/C=RU/ST=Omsk Oblast/L=Omsk/O=company/CN=%s", cn)
	if err := run("openssl", "req", "-new", "-key", keyPath, "-out", csrPath, "-subj", subj); err != nil {
		return fmt.Errorf("openssl req error: %w", err)
	}

	// openssl x509 -req -in csrPath -CA rootCrt -CAkey rootKey -CAcreateserial -out crtPath -days 9999 -sha256 -extfile extPath
	if err := run("openssl", "x509", "-req", "-in", csrPath, "-CA", rootCrt, "-CAkey", rootKey, "-CAcreateserial", "-out", crtPath, "-days", "9999", "-sha256", "-extfile", extPath); err != nil {
		return fmt.Errorf("openssl x509 error: %w", err)
	}

	return nil
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	output, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
