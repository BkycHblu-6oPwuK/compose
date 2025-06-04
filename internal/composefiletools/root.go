package composefiletools

import (
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile"
)

const (
	Dockerfile    string = "Dockerfile"
	Nginx         string = "nginx"
	ConfDir       string = "conf.d"
	App           string = "app"
	Mysql         string = "mysql"
	Postgres      string = "postgres"
	Sqlite        string = "sqlite"
	Memcached     string = "memcached"
	Redis         string = "redis"
	Mailhog       string = "mailhog"
	PhpMyAdmin    string = "phpmyadmin"
	Node          string = "node"
	Sphinx        string = "sphinx"
	Bin           string = "bin"
	Mysql_data    string = "mysql_data"
	Postgres_data string = "postgres_data"
	Redis_data    string = "redis_data"
	Sphinx_data   string = "sphinx_data"
)

var (
	AvailableFramework = [4]string{
		config.Bitrix,
		config.Laravel,
		config.Symfony,
		config.Vanilla,
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

func GetAvailableVersions(service string, yamlConfig *config.YamlConfig) []string {
	switch service {
	case App:
		switch yamlConfig.FrameworkName {
		case config.Laravel, config.Symfony, config.Vanilla:
			return []string{"8.2", "8.3", "8.4"}
		default:
			return []string{"7.4", "8.2", "8.3", "8.4"}
		}
	case Mysql:
		switch yamlConfig.FrameworkName {
		case config.Laravel, config.Symfony, config.Vanilla:
			return []string{"8.0", "latest"}
		default:
			return []string{"5.7", "8.0", "latest"}
		}
	case Postgres:
		return []string{"17", "latest"}
	default:
		return nil
	}
}

func GetCurrentDbType() (string, error) {
	compose, err := composefile.Load(config.GetDockerComposeFilePath())
	if err != nil {
		return "", err
	}
	switch true {
	case compose.Services.Has(Mysql):
		return Mysql, nil
	case compose.Services.Has(Postgres):
		return Postgres, nil
	default:
		return Sqlite, nil
	}
}

func GetCurrentServerCache() (string, error) {
	compose, err := composefile.Load(config.GetDockerComposeFilePath())
	if err != nil {
		return "", err
	}
	switch true {
	case compose.Services.Has(Redis):
		return Redis, nil
	case compose.Services.Has(Memcached):
		return Memcached, nil
	default:
		return "", nil
	}
}

func BuildYaml(yamlConfig *config.YamlConfig) *composefile.ComposeFile {
	fileBuilder := composefile.NewComposeFileBuilder().AddDefaultNetwork().AddService(Nginx, buildNginxService()).
		AddService(App, buildAppService(yamlConfig))

	switch yamlConfig.DbType {
	case Postgres:
		fileBuilder.AddService(Postgres, buildPostgresService()).AddVolume(Postgres_data, buildBaseVolume())
	case Mysql:
		fileBuilder.AddService(Mysql, buildMysqlService()).AddVolume(Mysql_data, buildBaseVolume())
	}

	switch yamlConfig.ServerCache {
	case Memcached:
		fileBuilder.AddService(Memcached, buildMemcachedService())
	case Redis:
		fileBuilder.AddService(Redis, buildRedisService()).AddVolume(Redis_data, buildBaseVolume())
	}

	if yamlConfig.CreateNode {
		fileBuilder.AddService(Node, buildNodeService())
	}
	if yamlConfig.CreateSphinx {
		fileBuilder.AddService(Sphinx, buildSphinxService()).AddVolume(Sphinx_data, buildBaseVolume())
	}
	fileBuilder.AddService(Mailhog, buildMailHogService())
	file := fileBuilder.Build()
	return &file
}
