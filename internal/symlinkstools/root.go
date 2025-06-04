package symlinkstools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/composefiletools"
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/filetools"
)

const (
	FileName = "symlinks"
)

func PushTosymlinks(symlinks map[string]string) error {
	filePath := filepath.Join(config.GetConfFilesDirPath(), composefiletools.App, FileName)

	if err := filetools.InitDirs(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("ошибка инициализации директории для симлинков: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла %s: %w", filePath, err)
	}
	defer file.Close()

	base := composefiletools.SitePathInContainer

	for src, dst := range symlinks {
		line := fmt.Sprintf("%s/%s %s/%s\n", base, src, base, dst)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("ошибка записи в файл: %w", err)
		}
	}

	return nil
}
