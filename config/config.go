package config

import (
	"os"
	"path/filepath"
)

var (
	configDir = ""
	dirName   = "kakemoti"
)

func init() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err.Error())
	}

	configDir = filepath.Join(dir, dirName)
	if _, err := os.Stat(configDir); err == nil {
		return
	}

	if err := os.Mkdir(configDir, os.ModePerm); err != nil {
		panic(err.Error())
	}
}

func ConfigDir() string {
	return configDir
}
