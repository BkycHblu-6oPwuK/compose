package yaml

import (
	"docky/config"
	"docky/utils"
	"os"

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
	Command       string            `yaml:"command,omitempty"`
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
	Dockerfile  = "Dockerfile"
	Nginx       = "nginx"
	App         = "app"
	Mysql       = "mysql"
	Node        = "node"
	Sphinx      = "sphinx"
	Bin         = "bin"
	Compose     = "compose"
	Mysql_data  = "mysql_data"
	Sphinx_data = "sphinx_data"
)

var (
	CreateNode           = false
	CreateSphinx         = false
	NodePath             = ""
	SitePath             = ""
	AvailablePhpVersions = [4]string{
		"7.4",
		"8.2",
		"8.3",
		"8.4",
	}
	AvailableMysqlVersions = [2]string{
		"5.7",
		"8.0",
	}
	PhpVersion   = "8.2"
	MysqlVersion = "8.0"
	NodeVersion  = "23"
)

func (c *ComposeFile) addService(name string, service Service) *ComposeFile {
	c.Services.Set(name, service)
	return c
}

func (c *ComposeFile) addNginxService() *ComposeFile {
	service := Service{
		Build: Build{Context: "${" + config.DockerPathVarName + "}", Dockerfile: "${" + config.DockerPathVarName + "}/" + Nginx + "/" + Dockerfile + "", Args: getBaseArgsBuild()},
		Volumes: []string{
			"${" + config.SitePathVarName + "}:" + config.SitePathInContainer,
			"${" + config.DockerPathVarName + "}/" + Nginx + "/conf.d:/etc/nginx/conf.d",
		},
		Ports:         []string{"80:80", "443:443"},
		Dependencies:  []string{App},
		Networks:      []string{Compose},
		ContainerName: Nginx,
	}

	return c.addService(Nginx, service)
}

func (c *ComposeFile) addAppService() *ComposeFile {
	phpVersionVarName := "${" + config.PhpVersionVarName + "}"
	service := Service{
		Build:        Build{Context: "${" + config.DockerPathVarName + "}", Dockerfile: "${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/" + Dockerfile + "", Args: getBaseArgsBuild()},
		Ports:        []string{"9000:9000"},
		Dependencies: []string{Mysql},
		Volumes: []string{
			"${" + config.SitePathVarName + "}:" + config.SitePathInContainer,
			"${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/php.ini:/usr/local/etc/php/conf.d/php.ini",
			"${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/xdebug.ini:/usr/local/etc/php/conf.d/xdebug.ini",
			"${" + config.DockerPathVarName + "}/" + App + "/php-fpm.conf:/usr/local/etc/php-fpm.d/zzzzwww.conf",
			"${" + config.DockerPathVarName + "}/" + App + "/nginx:/etc/nginx/conf.d",
		},
		ExstraHosts:   []string{"host.docker.internal:host-gateway"},
		Networks:      []string{Compose},
		ContainerName: App,
	}
	return c.addService(App, service)
}

func (c *ComposeFile) addMysqlSerice() *ComposeFile {
	service := Service{
		Image:   Mysql + ":${" + config.MysqlVersionVarName + "}",
		Restart: "always",
		Ports:   []string{"8102:3306"},
		Volumes: []string{
			Mysql_data + ":" + config.SitePathInContainer,
			"${" + config.DockerPathVarName + "}/" + Mysql + "/my.cnf:/etc/mysql/conf.d/my.cnf",
		},
		Networks: []string{Compose},
		Environment: map[string]string{
			"MYSQL_DATABASE":      "site",
			"MYSQL_ROOT_PASSWORD": "root",
		},
		ContainerName: Mysql,
	}
	return c.addService(Mysql, service)
}

func (c *ComposeFile) addNodeSerice() *ComposeFile {
	args := getBaseArgsBuild()
	args["NODE_VERSION"] = "${" + config.NodeVersionVarName + "}"
	args["NODE_PATH"] = "${" + config.NodePathVarName + "}"
	service := Service{
		Build: Build{Context: "${" + config.DockerPathVarName + "}", Dockerfile: "${" + config.DockerPathVarName + "}/" + Node + "/" + Dockerfile + "", Args: args},
		Ports: []string{"5173:5173", "5174:5174"},
		Volumes: []string{
			"${" + config.SitePathVarName + "}:" + config.SitePathInContainer,
		},
		Dependencies:  []string{App},
		Networks:      []string{Compose},
		Command:       "tail -f /dev/null",
		ContainerName: Node,
	}
	return c.addService(Node, service)
}

func (c *ComposeFile) addSphinxSerice() *ComposeFile {
	service := Service{
		Build:   Build{Context: "${" + config.DockerPathVarName + "}", Dockerfile: "${" + config.DockerPathVarName + "}/" + Sphinx + "/" + Dockerfile + ""},
		Restart: "always",
		Ports:   []string{"9312:9312", "9306:9306"},
		Volumes: []string{
			"${" + config.DockerPathVarName + "}/" + Sphinx + "/sphinx.conf:/usr/local/etc/sphinx.conf",
			Sphinx_data + ":/var/lib/sphinx/data",
		},
		Networks:      []string{Compose},
		ContainerName: Sphinx,
	}
	return c.addService(Sphinx, service)
}

func getBaseArgsBuild() map[string]string {
	return map[string]string{
		config.UserGroupVarName: "${" + config.UserGroupVarName + "}",
	}
}

func Create() error {
	volumes := map[string]Volume{
		Mysql_data: {},
	}
	if CreateSphinx {
		volumes[Sphinx_data] = Volume{}
	}
	networks := map[string]Network{
		Compose: {
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
	if CreateNode {
		file.addNodeSerice()
	}
	if CreateSphinx {
		file.addSphinxSerice()
	}
	return file.Save()
}

func (c *ComposeFile) Save() error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(config.GetDockerComposeFilePath(), out, 0644)
}

func Load() (*ComposeFile, error) {
	data, err := os.ReadFile(config.GetDockerComposeFilePath())
	if err != nil {
		return nil, err
	}

	compose := &ComposeFile{
		Services: utils.NewOrderedMap[string, Service](),
	}

	err = yaml.Unmarshal(data, compose)
	if err != nil {
		return nil, err
	}

	return compose, nil
}

func PublishNodeService() error {
	compose, err := Load()
	if err != nil {
		return err
	}
	if compose.Services.Has(Node) {
		return nil
	}
	compose.addNodeSerice()

	return compose.Save()
}

func PublishSphinxService() error {
	compose, err := Load()
	if err != nil {
		return err
	}
	if compose.Services.Has(Sphinx) {
		return nil
	}
	compose.addSphinxSerice()

	return compose.Save()
}
