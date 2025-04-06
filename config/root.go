package config

import (
	"docky/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ScriptName            = "docky"
	siteDirName           = "site" // директория с проектом
	dockerFilesDirName    = "_docker"
	localHostsFileName    = "hosts"
	dockerComposeFileName = "docker-compose.yml"
)

var (
	scriptCacheDir string
	curDirPath     string // директория из которой запускается команда
	workDirPath    string // директория с docker-compose.yml
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
	path, err := utils.FindFileUpwards(GetCurDirPath(), dockerComposeFileName)
	if err != nil {
		log.Fatalf("docker-compose.yml не найден: %v", err)
	}
	workDirPath = strings.TrimSuffix(path, "/"+dockerComposeFileName)
	return workDirPath
}
func GetWorkDirPath() string {
	if workDirPath != "" {
		return workDirPath
	}
	return getWorkDirPath()
}

func GetSiteDirPath() string {
	return filepath.Join(GetWorkDirPath(), siteDirName)
}
func GetDockerFilesDirPath() string {
	return filepath.Join(GetWorkDirPath(), dockerFilesDirName)
}
func GetLocalHostsFilePath() string {
	return filepath.Join(GetWorkDirPath(), localHostsFileName)
}
func GetDockerComposeFilePath() string {
	return filepath.Join(GetWorkDirPath(), dockerComposeFileName)
}
