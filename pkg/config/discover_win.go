//go:build windows
// +build windows

package config

import (
	"os"
	"os/user"
	"path/filepath"
)

const (
	envProgramData = "PROGRAMDATA"
	envAppData     = "APPDATA"
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
	}
}

func xdgDirectories() []string {
	return []string{}
}

func systemDirectories() []string {
	var directories []string

	if programData := os.Getenv(envProgramData); programData != "" {
		direcories = append(directories, programData)
	}

	if appData := os.Getenv(envAppData); appData != "" {
		directories = append(directories, appData)
	}

	return directories
}
