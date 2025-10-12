package main

import (
	"aleesa-webapp-go/internal/flickr"
	"aleesa-webapp-go/internal/webapp"
	"fmt"
	"os"
)

func main() {
	if err := webapp.ReadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		if err := flickr.APIClientInit(webapp.Config); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
