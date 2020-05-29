package appconfig

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	systemAppDir = "/var/lib/dodo"
)

func GetAppDir() string {
	dir := filepath.FromSlash(systemAppDir)
	if user, err := user.Current(); err == nil && user.HomeDir != "" {
		dir = filepath.Join(user.HomeDir, ".dodo")
	}

	os.MkdirAll(dir, 0700)

	return dir
}

func GetPluginDir() string {
	return filepath.Join(GetAppDir(), "plugins")
}
