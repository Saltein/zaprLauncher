package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// find AppData Path
func GetAppDataPath(appName string) string {

	var basePath string
	switch runtime.GOOS {

	case "windows":
		basePath = os.Getenv("APPDATA")
	case "darwin":
		basePath = filepath.Join(os.Getenv("HOME"), "Library", "Aplication Support")
	default:
		basePath = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	return filepath.Join(basePath, appName)

}
