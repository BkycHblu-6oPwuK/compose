package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func dirIsExists(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}
func DirIsExists(dir string) bool {
	return dirIsExists(dir)
}

func findFileUpwards(startDir, fileName string) (string, error) {
	dir := startDir
	fmt.Println(dir)
	for {
		path := filepath.Join(dir, fileName)
		fmt.Println("Проверяем файл:", dir)
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
func FindFileUpwards(startDir, fileName string) (string, error) {
	return findFileUpwards(startDir, fileName)
}