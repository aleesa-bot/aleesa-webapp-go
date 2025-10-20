package webapp

import (
	"fmt"
	"os"
	"path/filepath"
)

// DefaultConfigFileLocations выдаёт варианты путей, по которым может находиться конфиг-файл.
func DefaultConfigFileLocations() ([]string, error) {
	var locations []string

	executablePath, err := os.Executable()

	if err != nil {
		return locations, fmt.Errorf("unable to get my executable path: %w", err)
	}

	configJSONPath := filepath.Dir(executablePath) + "/data/config.json"

	locations = []string{
		"~/.aleesa-webapp-go.json",
		"~/aleesa-webapp-go.json",
		"/etc/aleesa-webapp-go.json",
		configJSONPath,
	}

	return locations, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
