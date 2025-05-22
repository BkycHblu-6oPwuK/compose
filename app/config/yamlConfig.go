package config

import (
	"os"
	"sync"
)

type YamlConfig struct {
	FrameworkName   string // bitrix, laravel и т.д.
	DbType          string // mysql, postgres, sqlite
	PhpVersion      string
	MysqlVersion    string
	PostgresVersion string
	NodeVersion     string
	NodePath        string
	SitePath        string
	CreateNode      bool
	CreateSphinx    bool
	ServerCache     string // memcached, redis
	UserGroup       string
}

var (
	cfg  *YamlConfig
	once sync.Once
)

func GetYamlConfig() *YamlConfig {
	once.Do(func() {
		nodeVersion := os.Getenv(NodeVersionVarName)
		if nodeVersion == "" {
			nodeVersion = "23"
		}
		cfg = &YamlConfig{
			FrameworkName:   os.Getenv(DockyFrameworkVarName),
			PhpVersion:      os.Getenv(PhpVersionVarName),
			MysqlVersion:    os.Getenv(MysqlVersionVarName),
			PostgresVersion: os.Getenv(PostgresVersionVarName),
			NodeVersion:     nodeVersion,
			NodePath:        os.Getenv(NodePathVarName),
			SitePath:        os.Getenv(SitePathVarName),
			UserGroup:       os.Getenv(UserGroupVarName),
		}
	})
	return cfg
}
