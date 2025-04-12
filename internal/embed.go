package internal

import (
	"docky/config"
	"docky/utils"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

//go:embed files/*
var files embed.FS

var MaxAgeCacheDir = 1 * 24 * time.Hour // 1 день по умолчанию

func cleanCacheDir(targetDir string) error {
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > MaxAgeCacheDir {
		fmt.Println("Удаление кэш директории:", targetDir)
		return os.RemoveAll(targetDir)
	}

	return nil
}

func extractAllFiles(targetDir string) error {
	return fs.WalkDir(files, "files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("files", path)
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

		fmt.Println("Запись:", dstPath)
		return os.WriteFile(dstPath, data, 0644)
	})
}
func ExtractFilesInCache() {
	targetDir := config.GetScriptCacheDir()
	err := cleanCacheDir(targetDir)
	if err != nil {
		log.Println("Ошибка при очистке кэш директории:", err)
	}
	if utils.FileIsExists(targetDir) {
		return
	}

	err = extractAllFiles(targetDir)
	if err != nil {
		log.Fatalf("Ошибка распаковки исходных файлов: %v", err)
	}

	fmt.Println("Распаковка завершена:", targetDir)
}
