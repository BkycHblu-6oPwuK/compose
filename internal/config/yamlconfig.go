package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config/framework"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
	"github.com/joho/godotenv"
)

type YamlConfig struct {
	FrameworkName   framework.Framework // bitrix, laravel и т.д.
	DbType          string              // mysql, postgres, sqlite
	PhpVersion      string
	MysqlVersion    string
	MariadbVersion  string
	PostgresVersion string
	NodeVersion     string
	NodePath        string
	CreateNode      bool
	CreateSphinx    bool
	ServerCache     string // memcached, redis
	UserGroup       string
}

const (
	UserGroupVarName       string = "USERGROUP"
	DockyFrameworkVarName  string = "DOCKY_FRAMEWORK"
	DockerPathVarName      string = "DOCKER_PATH"
	ConfPathVarName        string = "CONF_PATH"
	PhpVersionVarName      string = "PHP_VERSION"
	MysqlVersionVarName    string = "MYSQL_VERSION"
	MariadbVersionVarName  string = "MARIADB_VERSION"
	PostgresVersionVarName string = "POSTGRES_VERSION"
	NodeVersionVarName     string = "NODE_VERSION"
	SitePathVarName        string = "SITE_PATH"
	NodePathVarName        string = "NODE_PATH"
)

var (
	cfg  *YamlConfig
	once sync.Once
)

func GetYamlConfig() *YamlConfig {
	once.Do(func() {
		if err := loadEnvFile(); err != nil {
			panic("Ошибка при загрузке файла .env")
		}
		var frameworkValue framework.Framework
		if frameworkEnv := os.Getenv(DockyFrameworkVarName); frameworkEnv == "" {
			frameworkValue = ""
		} else {
			frameworkValue = framework.ParseFramework(frameworkEnv)
		}
		nodeVersion := os.Getenv(NodeVersionVarName)
		if nodeVersion == "" {
			nodeVersion = "23"
		}
		cfg = &YamlConfig{
			FrameworkName:   frameworkValue,
			PhpVersion:      os.Getenv(PhpVersionVarName),
			MysqlVersion:    os.Getenv(MysqlVersionVarName),
			MariadbVersion:  os.Getenv(MariadbVersionVarName),
			PostgresVersion: os.Getenv(PostgresVersionVarName),
			NodeVersion:     nodeVersion,
			NodePath:        os.Getenv(NodePathVarName),
			UserGroup:       os.Getenv(UserGroupVarName),
		}
	})
	return cfg
}

func loadEnvFile() error {
	envPath := GetEnvFilePath()
	if fileExists, _ := filetools.FileIsExists(envPath); fileExists {
		return godotenv.Load(envPath)
	}
	return nil
}

func GetUserGroup() string {
	userGroup := GetYamlConfig().UserGroup
	if userGroup == "" {
		userGroup = strconv.Itoa(os.Getegid())
		if userGroup == "0" {
			userGroup = "1000"
		}
	}
	return userGroup
}

func GetCurFramework() framework.Framework {
	curFramework := GetYamlConfig().FrameworkName
	if curFramework == "" {
		curFramework = framework.Bitrix
	}
	return curFramework
}
