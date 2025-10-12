package webapp

import (
	"fmt"
	"os"
	"path/filepath"
)

func DefaultConfigFileLocations() ([]string, error) {
	var locations []string

	executablePath, err := os.Executable()

	if err != nil {
		return locations, fmt.Errorf("unable to get my executable path: %w", err)
	}

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	locations = []string{
		"~/.aleesa-webapp-go.json",
		"~/aleesa-webapp-go.json",
		"/etc/aleesa-webapp-go.json",
		configJSONPath,
	}

	return locations, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
