//go:build darwin
// +build darwin

package config

import (
	"os"
	"os/user"
	"path/filepath"
)

func userDirectories() []string {
	user, err := user.Current()
	if err != nil {
		return []string{}
	}

	if user.HomeDir == "" {
		return []string{}
	}

	return []string{
		user.HomeDir,
		filepath.Join(user.HomeDir, "."+Name),
		filepath.Join(user.HomeDir, "Library", "Application Support"),
	}
}

func xdgDirectories() []string {
	var directories []string

	if xdgHome := os.Getenv(envXDGHome); xdgHome != "" {
		directories = append(directories, filepath.Join(xdgHome, Name))
	} else if user, err := user.Current(); err == nil && user.HomeDir != "" {
		directories = append(directories, filepath.Join(user.HomeDir, ".config", Name))
	}

	if xdgDirs := os.Getenv(envXDGDirs); xdgDirs != "" {
		for _, dir := range filepath.SplitList(xdgDirs) {
			directories = append(directories, filepath.Join(dir, Name))
		}
	} else if xdgDefaultDir != "" {
		directories = append(directories, filepath.Join(xdgDefaultDir, Name))
	}

	return directories
}

func systemDirectories() []string {
	return []string{
		filepath.Join("/", "etc"),
		filepath.Join("/", "Library", "Application Support"),
	}
}
