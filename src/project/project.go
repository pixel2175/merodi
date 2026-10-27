package project

import (
	. "github.com/pixel2175/errs"
	"merodi/src/utils/path"
	"os"
	"path/filepath"
)

func GetProjectDir(dir string) (string, error) {
	if dir == "" {
		dir = path.Getwd()
	}

	return filepath.Abs(dir)
}

func SetProjectDir(dir string) (string, error) {
	projectPath := TryV(GetProjectDir(dir))

	for projectPath != "/" {
		if IsProjectExists(projectPath) {
			return projectPath, nil
		}

		projectPath = filepath.Dir(projectPath)
	}

	return "", os.ErrNotExist
}

func IsProjectExists(path string) bool {
	configPath := filepath.Join(path, "config.toml")

	_, err := os.Stat(configPath)
	return err == nil
}
