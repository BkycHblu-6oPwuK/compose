package models

import (
	"docky/config"
	"docky/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ComposeFile struct {
	Services *utils.OrderedMap[string, Service] `yaml:"services"`
	Volumes  map[string]Volume                  `yaml:"volumes,omitempty"`
	Networks map[string]Network                 `yaml:"networks,omitempty"`
	Secrets  map[string]Secret                  `yaml:"secrets,omitempty"`
}

type Service struct {
	Image         string            `yaml:"image,omitempty"`
	Build         Build             `yaml:"build,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty"`
	Dependencies  []string          `yaml:"depends_on,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	ExstraHosts   []string          `yaml:"extra_hosts,omitempty"`
	Secrets       []string          `yaml:"secrets,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
}

type Build struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
}

type Network struct {
	Driver string `yaml:"driver,omitempty"`
}

type Volume struct {
	Driver string `yaml:"driver,omitempty"`
}

type Secret struct {
	File string `yaml:"file,omitempty"`
}

const (
	dockerfile = "Dockerfile"
	nginx      = "nginx"
	app        = "app"
	mysql      = "mysql"
	node       = "node"
	sphinx     = "sphinx"
	bin        = "bin"
	compose    = "compose"
	mysql_data = "mysql_data"
)

func (c *ComposeFile) addService(name string, service Service) *ComposeFile {
	c.Services.Set(name, service)
	return c
}

func (c *ComposeFile) addNginxService() *ComposeFile {
	dockerDirPath := config.GetCurrentDockerFileDirPath()
	nginxDockerPath := filepath.Join(dockerDirPath, nginx)
	service := Service{
		Image:   "",
		Build:   Build{Context: dockerDirPath, Dockerfile: filepath.Join(nginxDockerPath, dockerfile), Args: getBaseArgsBuild()},
		Restart: "",
		Volumes: []string{
			config.GetSiteDirPath() + ":/var/www",
			filepath.Join(nginxDockerPath, "conf.d") + ":/var/www",
		},
		Ports:         []string{"80:80", "443:443"},
		Dependencies:  []string{app},
		Networks:      []string{compose},
		ContainerName: nginx,
	}

	return c.addService(nginx, service)
}

func (c *ComposeFile) addAppService() *ComposeFile {
	dockerDirPath := config.GetCurrentDockerFileDirPath()
	appDockerPath := filepath.Join(dockerDirPath, app)
	phpDockerPath := filepath.Join(appDockerPath, "php-8.2")
	service := Service{
		Image:        "",
		Build:        Build{Context: dockerDirPath, Dockerfile: filepath.Join(phpDockerPath, dockerfile), Args: getBaseArgsBuild()},
		Restart:      "",
		Ports:        []string{"9000:9000"},
		Dependencies: []string{mysql},
		Volumes: []string{
			config.GetSiteDirPath() + ":/var/www",
			filepath.Join(phpDockerPath, "php.ini") + ":/usr/local/etc/php/conf.d/php.ini",
			filepath.Join(appDockerPath, "php-fpm.conf") + ":/usr/local/etc/php-fpm.d/zzzzwww.conf",
			filepath.Join(appDockerPath, "nginx") + ":/etc/nginx/conf.d",
		},
		ExstraHosts:   []string{"host.docker.internal:host-gateway"},
		Networks:      []string{compose},
		ContainerName: app,
	}
	return c.addService(app, service)
}

func (c *ComposeFile) addMysqlSerice() *ComposeFile {
	service := Service{
		Image:   mysql + ":8.0",
		Restart: "always",
		Ports:   []string{"8102:3306"},
		Volumes: []string{
			mysql_data + ":/var/www",
			filepath.Join(config.GetCurrentDockerFileDirPath(), mysql, "my.cnf") + ":/etc/mysql/conf.d/my.cnf",
		},
		Networks:      []string{compose},
		ContainerName: mysql,
	}
	return c.addService(mysql, service)
}

func getBaseArgsBuild() map[string]string {
	return map[string]string{
		"USERGROUP":   "${" + config.UserGroupVarName + "}",
		"DOCKER_PATH": config.GetCurrentDockerFileDirPath(),
	}
}

func CreateYmlFile() {
	volumes := map[string]Volume{
		"mysql_data": {},
	}

	networks := map[string]Network{
		"compose": {
			Driver: "bridge",
		},
	}

	secrets := map[string]Secret{}
	file := &ComposeFile{
		Services: utils.NewOrderedMap[string, Service](),
		Volumes:  volumes,
		Networks: networks,
		Secrets:  secrets,
	}
	file.addNginxService().addAppService().addMysqlSerice()
	out, err := yaml.Marshal(file)
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("docker-compose.yml", out, 0644)
}
