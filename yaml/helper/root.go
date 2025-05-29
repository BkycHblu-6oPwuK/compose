package helper
//@todo убрать из volume конф.файлы, копировать все в докерфайлах. Тогда не будет проблем с кастомизацией. publish --file php.ini в _conf/php.ini + add volume
import (
	"docky/config"
	"docky/yaml"
)

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
	PhpMyAdmin    string = "phpmyadmin"
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

func GetAvailableVersions(service string, yamlConfig *config.YamlConfig) []string {
	switch service {
	case App:
		if yamlConfig.FrameworkName == config.Laravel {
			return []string{"8.2", "8.3", "8.4"}
		}
		return []string{"7.4", "8.2", "8.3", "8.4"}
	case Mysql:
		if yamlConfig.FrameworkName == config.Laravel {
			return []string{"8.0", "latest"}
		}
		return []string{"5.7", "8.0", "latest"}
	case Postgres:
		return []string{"17", "latest"}
	default:
		return nil
	}
}

func GetCurrentDbType() (string, error) {
	compose, err := yaml.Load()
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
	compose, err := yaml.Load()
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

func BuildYaml(yamlConfig *config.YamlConfig) *yaml.ComposeFile {
	fileBuilder := yaml.NewComposeFileBuilder().AddDefaultNetwork().AddService(Nginx, buildNginxService()).
		AddService(App, buildAppService(yamlConfig))

	switch yamlConfig.FrameworkName {
	case config.Laravel:
		buildLaravelYaml(fileBuilder, yamlConfig)
	case config.Bitrix:
		buildBitrixYaml(fileBuilder)
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

func buildLaravelYaml(fileBuilder *yaml.ComposeFileBuilder, yamlConfig *config.YamlConfig) {
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
}

func buildBitrixYaml(fileBuilder *yaml.ComposeFileBuilder) {
	fileBuilder.AddService(Mysql, buildMysqlService()).AddVolume(Mysql_data, buildBaseVolume())
}
