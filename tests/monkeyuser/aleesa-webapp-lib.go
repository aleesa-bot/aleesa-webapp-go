package main

import (
	"bufio"
	"os"

	log "github.com/sirupsen/logrus"
)

// Читает даденный файл построчно в массив строк.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Unable to close file %s:%s", path, err)
		}
	}(file)

	var (
		lines   []string
		scanner = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
