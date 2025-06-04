package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BkycHblu-6oPwuK/docky/pkg/filetools"

	"github.com/joho/godotenv"
)

const (
	ScriptName             string = "docky"
	SiteDirName            string = "site"
	DockerFilesDirName     string = "_docker"
	ConfFilesDirName       string = "_conf"
	Bitrix                 string = "bitrix"
	Laravel                string = "laravel"
	Vanilla                string = "vanilla"
	Symfony                string = "symfony"
	LocalHostsFileName     string = "hosts"
	DockerComposeFileName  string = "docker-compose.yml"
	UserGroupVarName       string = "USERGROUP"
	DockyFrameworkVarName  string = "DOCKY_FRAMEWORK"
	DockerPathVarName      string = "DOCKER_PATH"
	ConfPathVarName        string = "CONF_PATH"
	PhpVersionVarName      string = "PHP_VERSION"
	MysqlVersionVarName    string = "MYSQL_VERSION"
	PostgresVersionVarName string = "POSTGRES_VERSION"
	NodeVersionVarName     string = "NODE_VERSION"
	SitePathVarName        string = "SITE_PATH"
	NodePathVarName        string = "NODE_PATH"
	SitePathInContainer    string = "/var/www"
	EnvFile                string = ".env"
)

var (
	scriptCacheDir string
	curDirPath     string // директория из которой запускается команда
	workDirPath    string // директория с docker-compose.yml
	Timestamp      = strconv.Itoa(int(time.Now().Unix()))
)

func getScriptCacheDir() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("Не найдена кэш директория: %v", err)
	}
	scriptCacheDir = filepath.Join(cacheDir, ScriptName)
	return scriptCacheDir
}
func GetScriptCacheDir() string {
	if scriptCacheDir != "" {
		return scriptCacheDir
	}
	return getScriptCacheDir()
}

func getCurDirPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Не найдена текущая директория: %v", err)
	}
	curDirPath = cwd
	return curDirPath
}
func GetCurDirPath() string {
	if curDirPath != "" {
		return curDirPath
	}
	return getCurDirPath()
}

func getWorkDirPath() string {
	path, err := filetools.FindFileUpwards(GetCurDirPath(), DockerComposeFileName)
	if err != nil {
		path = GetCurDirPath()
	}
	workDirPath = strings.TrimSuffix(path, "/"+DockerComposeFileName)
	return workDirPath
}
func GetWorkDirPath() string {
	if workDirPath != "" {
		return workDirPath
	}
	return getWorkDirPath()
}

func GetSiteDirPath() string {
	return filepath.Join(GetWorkDirPath(), SiteDirName)
}

func GetUserGroup() string {
	userGroup := GetYamlConfig().UserGroup
	if userGroup == "" {
		userGroup = strconv.Itoa(os.Getegid())
	}
	return userGroup
}
func GetCurFramework() string {
	frameworkName := os.Getenv(DockyFrameworkVarName)
	if frameworkName != "" {
		return frameworkName
	}

	return Bitrix
}
func GetDockerFilesDirPath() string {
	return filepath.Join(GetWorkDirPath(), DockerFilesDirName)
}
func GetDockerFilesDirPathInCache() string {
	return filepath.Join(GetScriptCacheDir(), DockerFilesDirName, GetCurFramework())
}
func GetCurrentDockerFileDirPath() string {
	path := GetDockerFilesDirPath()
	if fileExists, _ := filetools.FileIsExists(path); fileExists {
		return path
	}
	path = GetDockerFilesDirPathInCache()
	return path
}
func GetConfFilesDirPath() string {
	return filepath.Join(GetWorkDirPath(), ConfFilesDirName)
}
func GetLocalHostsFilePath() string {
	return filepath.Join(GetConfFilesDirPath(), LocalHostsFileName)
}
func GetDockerComposeFilePath() string {
	return filepath.Join(GetWorkDirPath(), DockerComposeFileName)
}
func GetEnvFilePath() string {
	return filepath.Join(GetWorkDirPath(), EnvFile)
}
func loadEnvFile() error {
	envPath := GetEnvFilePath()
	if fileExists, _ := filetools.FileIsExists(envPath); fileExists {
		return godotenv.Load(envPath)
	}
	return nil
}

func init() {
	loadEnvFile()
}
