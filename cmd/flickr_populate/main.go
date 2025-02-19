package main

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/flickr"
	"fmt"
	"os"
)

func main() {
	if c, err := config.ReadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		if err := flickr.Populate(c); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
