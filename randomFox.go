package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

func randomFoxClient() (string, error) {
	var (
		err error
		c   = http.Client{
			Timeout: 10 * time.Second,
		}
		URL = "https://randomfox.ca/floof/"
	)

	req, err := http.NewRequest(http.MethodGet, URL, nil) //nolint: noctx
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	var resp *http.Response
	resp, err = c.Do(req) //nolint: bodyclose

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Errorf("Unable to close response body for randomfox.ca api request: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var fox randomFox
	err = json.Unmarshal(respBody, &fox)

	if err != nil {
		return "", err
	}

	var re = regexp.MustCompile(`\\`)
	url := re.ReplaceAll([]byte(fox.Link), []byte(""))

	return string(url), nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
