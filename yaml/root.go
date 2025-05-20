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
	Config   *config.YamlConfig                 `yaml:"-"`
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
	Command       interface{}       `yaml:"command,omitempty"`
	ExtraHosts    []string          `yaml:"extra_hosts,omitempty"`
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
	Dockerfile    string = "Dockerfile"
	Nginx         string = "nginx"
	App           string = "app"
	Mysql         string = "mysql"
	Postgres      string = "postgres"
	Sqlite        string = "sqlite"
	Memcached     string = "memcached"
	Redis         string = "redis"
	Mailhog       string = "mailhog"
	Node          string = "node"
	Sphinx        string = "sphinx"
	Bin           string = "bin"
	Compose       string = "compose"
	Mysql_data    string = "mysql_data"
	Postgres_data string = "postgres_data"
	Redis_data    string = "redis_data"
	Sphinx_data   string = "sphinx_data"
)

var (
	AvailableFramework = [2]string{
		config.Bitrix,
		config.Laravel,
	}
	AvailableDb = [3]string{
		Mysql,
		Postgres,
		Sqlite,
	}
	AvailableServerCache = [2]string{
		Memcached,
		Redis,
	}
)

func NewYamlFile(cfg *config.YamlConfig) *ComposeFile {
	return &ComposeFile{
		Services: utils.NewOrderedMap[string, Service](),
		Volumes:  make(map[string]Volume),
		Networks: map[string]Network{
			Compose: {
				Driver: "bridge",
			},
		},
		Config: cfg,
	}
}

func (c *ComposeFile) GetAvailableVersions(service string) []string {
	switch service {
	case App:
		if c.Config.FrameworkName == config.Laravel {
			return []string{"8.2", "8.3", "8.4"}
		}
		return []string{"7.4", "8.2", "8.3", "8.4"}
	case Mysql:
		if c.Config.FrameworkName == config.Laravel {
			return []string{"8.0", "latest"}
		}
		return []string{"5.7", "8.0", "latest"}
	case Postgres:
		return []string{"17", "latest"}
	default:
		return nil
	}
}

func (c *ComposeFile) addService(name string, service Service) *ComposeFile {
	if c.Services == nil {
		c.Services = utils.NewOrderedMap[string, Service]()
	}
	c.Services.Set(name, service)
	return c
}
func (c *ComposeFile) addVolume(name string, volume Volume) {
	if c.Volumes == nil {
		c.Volumes = make(map[string]Volume)
	}
	c.Volumes[name] = volume
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
		Build: Build{Context: "${" + config.DockerPathVarName + "}", Dockerfile: "${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/" + Dockerfile + "", Args: getBaseArgsBuild()},
		Ports: []string{
			"9000:9000",
			"6001:6001",
		},
		Volumes: []string{
			"${" + config.SitePathVarName + "}:" + config.SitePathInContainer,
			"${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/php.ini:/usr/local/etc/php/conf.d/php.ini",
			"${" + config.DockerPathVarName + "}/" + App + "/php-" + phpVersionVarName + "/xdebug.ini:/usr/local/etc/php/conf.d/xdebug.ini",
			"${" + config.DockerPathVarName + "}/" + App + "/php-fpm.conf:/usr/local/etc/php-fpm.d/zzzzwww.conf",
			"${" + config.DockerPathVarName + "}/" + App + "/nginx:/etc/nginx/conf.d",
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Networks:   []string{Compose},
		Environment: map[string]string{
			"XDEBUG_TRIGGER": "testTrig",
			"PHP_IDE_CONFIG": "serverName=xdebugServer",
		},
		ContainerName: App,
	}
	if c.Config.DbType != Sqlite {
		service.Dependencies = []string{c.Config.DbType}
	}
	return c.addService(App, service)
}

func (c *ComposeFile) addMysqlService() *ComposeFile {
	service := Service{
		Image:   Mysql + ":${" + config.MysqlVersionVarName + "}",
		Restart: "always",
		Ports:   []string{"8102:3306"},
		Volumes: []string{
			Mysql_data + ":/var/lib/mysql",
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

func (c *ComposeFile) addPostgresService() *ComposeFile {
	service := Service{
		Image:   Postgres + ":${" + config.PostgresVersionVarName + "}",
		Restart: "always",
		Ports:   []string{"5432:5432"},
		Volumes: []string{
			Postgres_data + ":/var/lib/postgresql/data",
			"${" + config.DockerPathVarName + "}/" + Postgres + "/postgresql.conf:/etc/postgresql/postgresql.conf",
		},
		Networks: []string{Compose},
		Environment: map[string]string{
			"POSTGRES_DB":       "site",
			"POSTGRES_PASSWORD": "root",
			"POSTGRES_USER":     Postgres,
		},
		ContainerName: Postgres,
		Command:       []string{"-c", "config_file=/etc/postgresql/postgresql.conf"},
	}
	return c.addService(Postgres, service)
}

func (c *ComposeFile) addNodeService() *ComposeFile {
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

func (c *ComposeFile) addSphinxService() *ComposeFile {
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

func (c *ComposeFile) addMemcachedService() *ComposeFile {
	service := Service{
		Image:         Memcached,
		Ports:         []string{"11211:11211"},
		Networks:      []string{Compose},
		ContainerName: Memcached,
		Command:       []string{"--conn-limit=1024", "--memory-limit=64", "--threads=4"},
	}
	return c.addService(Memcached, service)
}

func (c *ComposeFile) addRedisService() *ComposeFile {
	service := Service{
		Image: Redis,
		Ports: []string{"6379:6379"},
		Volumes: []string{
			Redis_data + ":/data",
		},
		Networks:      []string{Compose},
		ContainerName: Redis,
		Command:       []string{"redis-server", "--appendonly", "yes"},
	}
	return c.addService(Redis, service)
}

func (c *ComposeFile) addMailHogService() *ComposeFile {
	service := Service{
		Image: "mailhog/mailhog",
		Ports: []string{
			"1025:1025",
			"8025:8025",
		},
		Networks:      []string{Compose},
		ContainerName: Mailhog,
	}
	return c.addService(Mailhog, service)
}

func getBaseArgsBuild() map[string]string {
	return map[string]string{
		config.UserGroupVarName: "${" + config.UserGroupVarName + "}",
	}
}

func (c *ComposeFile) Create() error {
	c.addNginxService().addAppService()

	switch c.Config.FrameworkName {
	case config.Laravel:
		switch c.Config.DbType {
		case Postgres:
			c.addPostgresService()
			c.addVolume(Postgres_data, Volume{})
		case Mysql:
			c.addMysqlService()
			c.addVolume(Mysql_data, Volume{})
		}

		switch c.Config.ServerCache {
		case Memcached:
			c.addMemcachedService()
		case Redis:
			c.addRedisService()
			c.addVolume(Redis_data, Volume{})
		}
		c.addMailHogService()
	default:
		c.addMysqlService()
		c.addVolume(Mysql_data, Volume{})
	}

	if c.Config.CreateNode {
		c.addNodeService()
	}
	if c.Config.CreateSphinx {
		c.addSphinxService()
		c.addVolume(Sphinx_data, Volume{})
	}
	return c.Save()
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
		Volumes:  make(map[string]Volume),
		Networks: make(map[string]Network),
		Config:   nil,
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
	compose.addNodeService()

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
	compose.addSphinxService().addVolume(Sphinx_data, Volume{})
	return compose.Save()
}

func PublisRedisService() error {
	compose, err := Load()
	if err != nil {
		return err
	}
	if compose.Services.Has(Redis) {
		return nil
	}
	compose.addRedisService().addVolume(Redis_data, Volume{})
	return compose.Save()
}

func PublishMemcachedService() error {
	compose, err := Load()
	if err != nil {
		return err
	}
	if compose.Services.Has(Memcached) {
		return nil
	}
	compose.addMemcachedService()

	return compose.Save()
}
