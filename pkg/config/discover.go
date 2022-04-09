package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envXDGHome = "XDG_CONFIG_HOME"
	envXDGDirs = "XDG_CONFIG_DIRS"

	xdgDefaultDir = "/etc/etc/xdg"
)

func discoverConfigFiles() []string {
	directories := []string{}
	directories = append(directories, cwdDirectories()...)
	directories = append(directories, userDirectories()...)
	directories = append(directories, xdgDirectories()...)
	directories = append(directories, systemDirectories()...)
	directories = uniqueStrings(directories)

	patterns := []string{}
	for _, ext := range supportedExtensions() {
		patterns = append(patterns, globPatternsForExtension(ext)...)
	}

	files := []string{}

	for _, directory := range directories {
		for _, pattern := range patterns {
			matches, err := filepath.Glob(filepath.Join(directory, pattern))
			if err != nil {
				continue
			}

			for _, filename := range matches {
				path, err := filepath.Abs(filename)
				if err != nil {
					continue
				}

				if _, err := os.Stat(path); os.IsNotExist(err) {
					continue
				}

				files = append(files, path)
			}
		}
	}

	return files
}

func cwdDirectories() []string {
	workingDir, err := os.Getwd()
	if err != nil {
		return []string{}
	}

	directories := []string{}
	for directory := workingDir; !isFSRoot(directory); directory = filepath.Dir(directory) {
		directories = append(directories, directory)
	}

	return directories
}

func supportedExtensions() []string {
	return []string{"yaml", "yml", "json"}
}

func globPatternsForExtension(ext string) []string {
	return []string{
		fmt.Sprintf("%s.%s", Name, ext),
		fmt.Sprintf(".%s.%s", Name, ext),
		filepath.Join(Name, fmt.Sprintf("config.%s", ext)),
		filepath.Join(fmt.Sprintf(".%s", Name), fmt.Sprintf("config.%s", ext)),
		filepath.Join(fmt.Sprintf("%s.%s.d", Name, ext), "*"),
		filepath.Join(fmt.Sprintf(".%s.%s.d", Name, ext), "*"),
	}
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	index := 0

	for _, item := range values {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = true
		values[index] = item
		index++
	}

	return values[:index]
}

func isFSRoot(path string) bool {
	return strings.HasSuffix(filepath.Clean(path), string(filepath.Separator))
}
