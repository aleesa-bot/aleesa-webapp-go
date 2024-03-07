package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func theCatAPIClient() (string, error) {
	var (
		err error
		c   = http.Client{
			Timeout: 10 * time.Second,
		}
		URL = "https://api.thecatapi.com/v1/images/search"
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
			log.Errorf("Unable to close response body for request to %s: %s", URL, err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("request to %s failed: %s", URL, resp.Status)

		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var cat theCatAPI
	err = json.Unmarshal(respBody, &cat)

	if err != nil {
		return "", err
	}

	return cat[0].Url, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
