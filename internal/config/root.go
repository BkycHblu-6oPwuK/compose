package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
)

const (
	ScriptName            string = "docky"
	ScriptVersion         string = "2.2.1"
	SiteDirName           string = "site"
	DockerFilesDirName    string = "_docker"
	ConfFilesDirName      string = "_conf"
	EnvFile               string = ".env"
	LocalHostsFileName    string = "hosts"
	DockerComposeFileName string = "docker-compose.yml"
)

var (
	scriptCacheDir string
	curDirPath     string // директория из которой запускается команда
	workDirPath    string // директория с docker-compose.yml
	timeStamp      string
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

func GetDockerFilesDirPath() string {
	return filepath.Join(GetWorkDirPath(), DockerFilesDirName)
}

func GetDockerFilesDirPathInCache() string {
	return filepath.Join(GetScriptCacheDir(), DockerFilesDirName, GetCurFramework().String())
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

func GetTimeStamp() string {
	if timeStamp == "" {
		timeStamp = strconv.Itoa(int(time.Now().Unix()))
	}
	return timeStamp
}
