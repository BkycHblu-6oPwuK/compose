package internal

import (
	"docky/config"
	"docky/utils"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

//go:embed files/*
var files embed.FS

const rootDir = "files"

func cleanCacheDir(targetDir string, MaxAgeCacheDir time.Duration) error {
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > MaxAgeCacheDir {
		return os.RemoveAll(targetDir)
	}

	return nil
}

func extractAllFiles(targetDir string, root string) error {
	return fs.WalkDir(files, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(targetDir, relPath)

		data, err := files.ReadFile(path)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(dstPath), 0755)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, 0644)
	})
}

func CleanCacheDir() error {
	return cleanCacheDir(config.GetScriptCacheDir(), 0)
}

func ExtractFilesInCache() error {
	targetDir := config.GetScriptCacheDir()
	err := cleanCacheDir(targetDir, 7*24*time.Hour)
	if err != nil {
		return err
	}
	if fileExists, _ := utils.FileIsExists(targetDir); fileExists {
		return nil
	}

	err = extractAllFiles(targetDir, rootDir)
	if err != nil {
		return err
	}
	return nil
}

func PublishFiles() error {
	targetDir := config.GetDockerFilesDirPath()
	var err error = nil
	if fileExists, _ := utils.FileIsExists(targetDir); fileExists {
		err = os.Rename(targetDir, targetDir+config.Timestamp)
		if err != nil {
			return err
		}
	}

	err = extractAllFiles(targetDir, filepath.Join(rootDir, config.DockerFilesDirName, config.GetCurFramework()))
	if err != nil {
		return err
	}
	return err
}

func PublishFile(filePath, targetPath string, isDir bool) error {
	if !isDir {
		data, err := files.ReadFile(filepath.Join(rootDir, filePath))
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}

		err = os.WriteFile(targetPath, data, 0644)
		if err != nil {
			return err
		}
	} else {
		err := extractAllFiles(targetPath, filepath.Join(rootDir, filePath))
		if err != nil {
			return err
		}
	}

	return nil
}
