package filetools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func fileIsExists(path string) (fileExists bool, isDir bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false, false
	}
	return true, info.IsDir()
}

// Проверяет существует ли файл по указанному пути или директория
func FileIsExists(path string) (fileExists bool, isDir bool) {
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

func MoveDirContents(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			err = os.Rename(srcPath, dstPath)
			if err != nil {
				err = CopyDir(srcPath, dstPath)
				if err != nil {
					return err
				}
				os.RemoveAll(srcPath)
			}
		} else {
			err = os.Rename(srcPath, dstPath)
			if err != nil {
				err = CopyFile(srcPath, dstPath)
				if err != nil {
					return err
				}
				os.Remove(srcPath)
			}
		}
	}

	return os.RemoveAll(srcDir)
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	err = InitDirs(filepath.Dir(dst))
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func CopyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = InitDirs(dst)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
		} else {
			err = CopyFile(srcPath, dstPath)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func IsDirEmpty(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdir(1)
	return err == io.EOF
}

func InitDirs(paths ...string) error {
	for _, path := range paths {
		if fileExists, _ := fileIsExists(path); !fileExists {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("error initialize directory %s: %v", path, err)
			}
		}
	}
	return nil
}