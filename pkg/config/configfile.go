package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	configFileNames = []string{
		"dodo.yaml",
		"dodo.yml",
		"dodo.json",
		".dodo.yaml",
		".dodo.yml",
		".dodo.json",
	}
)

// LoadConfiguration tries to find a backdrop configuration by name in any of
// the supported locations. If given, will only look in the supplied config
// file.
func LoadConfiguration(
	backdrop string, configfile string,
) (*BackdropConfig, error) {
	if configfile != "" {
		return FindConfigInFile(backdrop, configfile)
	}
	config, err := FindConfigAnywhere(backdrop)
	if err == nil {
		return config, nil
	}
	log.WithFields(log.Fields{
		"name":   backdrop,
		"reason": err,
	}).Debug("No valid config found anywhere")
	return FallbackConfig(backdrop)
}

// FindConfigAnywhere tries to find a backdrop configuration by name in any of
// the supported locations.
func FindConfigAnywhere(backdrop string) (*BackdropConfig, error) {
	directories, err := FindConfigDirectories()
	if err != nil {
		return nil, err
	}

	for _, directory := range directories {
		config, err := FindConfigInDirectory(backdrop, directory)
		if err == nil {
			return config, err
		}
		log.WithFields(log.Fields{
			"name":      backdrop,
			"directory": directory,
			"reason":    err,
		}).Debug("No valid config found in directory")
	}
	return nil, fmt.Errorf("Could not find configuration for backdrop '%s' in any configuration file", backdrop)
}

// FindConfigInDirectory tries to find a backdrop configuration by name in
// any of the default files in a specified directory.
func FindConfigInDirectory(
	backdrop string, directory string,
) (*BackdropConfig, error) {
	for _, filename := range configFileNames {
		path, err := filepath.Abs(filepath.Join(directory, filename))
		if err != nil {
			log.Error(err)
		}
		config, err := FindConfigInFile(backdrop, path)
		if err == nil {
			return config, err
		}
		log.WithFields(log.Fields{
			"name":   backdrop,
			"file":   path,
			"reason": err,
		}).Debug("No valid config found in file")
	}
	return nil, fmt.Errorf("Could not find configuration for backdrop '%s' in directory '%s'", backdrop, directory)
}

// FindConfigInFile tries to find a backdrop configuration by name in a specific
// file.
func FindConfigInFile(
	backdrop string, filename string,
) (*BackdropConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Could not read file '%s'", filename)
	}

	config, err := ParseConfiguration(filename, bytes)
	if err != nil {
		return nil, fmt.Errorf("Could not parse config from '%s': %s", filename, err)
	}

	if config.Backdrops == nil && config.Includes == nil {
		return nil, fmt.Errorf("File '%s' does not contain any backdrop configurations", filename)
	}

	if backdropConfig, ok := config.Backdrops[backdrop]; ok {
		return &backdropConfig, nil
	}

	directory := filepath.Dir(filename)
	for _, include := range config.Includes {
		if include.File != "" {
			path, err := filepath.Abs(filepath.Join(directory, include.File))
			if err != nil {
				log.Error(err)
			}
			// TODO: this may create infinite loops
			config, err := FindConfigInFile(backdrop, path)
			if err == nil {
				return config, err
			}
			log.WithFields(log.Fields{
				"name":   backdrop,
				"file":   path,
				"reason": err,
			}).Debug("No valid config found in file")
		} else if include.Text != "" {
			config, err := ParseConfiguration(filename, []byte(include.Text))
			if err != nil {
				return nil, fmt.Errorf("Could not parse config from include in '%s': %s", filename, err)
			}
			if config.Backdrops == nil {
				return nil, fmt.Errorf("Include in file '%s' does not contain any backdrop configurations", filename)
			}
			if backdropConfig, ok := config.Backdrops[backdrop]; ok {
				return &backdropConfig, nil
			}
			log.WithFields(log.Fields{
				"name":   backdrop,
				"file":   filename,
				"reason": err,
			}).Debug("No valid config in include of file")
		}
	}

	return nil, fmt.Errorf("File '%s' does not contain configuration for backdrop '%s'", filename, backdrop)
}

// FallbackConfig guesses a general-purpose backdrop configuration based
// on the name, that can be used in case no better configuration was found.
func FallbackConfig(backdrop string) (*BackdropConfig, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	config := &BackdropConfig{
		Image:      backdrop,
		Script:     fmt.Sprintf("%s $@", backdrop),
		User:       fmt.Sprintf("%s:%s", user.Uid, user.Gid),
		WorkingDir: workingDir,
		Volumes: []Volume{Volume{
			Source: workingDir,
			Target: workingDir,
		}},
	}

	return config, nil
}

// ListConfigurations prints out all available backdrop names and the file
// it was found in.
func ListConfigurations() {
	result := map[string]string{}
	directories, err := FindConfigDirectories()
	if err != nil {
		return
	}

	for _, directory := range directories {
		for _, filename := range configFileNames {
			path, err := filepath.Abs(filepath.Join(directory, filename))
			if err != nil {
				continue
			}

			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			config, err := ParseConfiguration(path, bytes)
			if err != nil {
				continue
			}

			for name := range config.Backdrops {
				if result[name] == "" {
					log.WithFields(log.Fields{
						"file": path,
					}).Info(name)
				}
			}
		}
	}
}
