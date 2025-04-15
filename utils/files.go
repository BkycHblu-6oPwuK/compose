package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func fileIsExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Проверяет существует ли файл по указанному пути или директория
func FileIsExists(path string) bool {
	return fileIsExists(path)
}

func findFileUpwards(startDir, fileName string) (string, error) {
	dir := startDir
	for {
		path := filepath.Join(dir, fileName)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("файл %s не найден начиная с %s", fileName, startDir)
}

// Находит путь к директории с указанным файлом
func FindFileUpwards(startDir, fileName string) (string, error) {
	return findFileUpwards(startDir, fileName)
}
