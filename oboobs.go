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

func oboobsClient() (string, error) {
	var (
		err error
		c   = http.Client{
			Timeout: 10 * time.Second,
		}

		URL = "http://api.oboobs.ru/boobs/0/1/random"
	)

	req, err := http.NewRequest(http.MethodGet, URL, nil)

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

	var boobs []oboobs
	err = json.Unmarshal(respBody, &boobs)

	if err != nil {
		return "", err
	}

	if len(boobs) == 0 {
		err = fmt.Errorf("empty json array returned from %s", URL)

		return "", err
	}

	answer := fmt.Sprintf("https://media.oboobs.ru/%s", boobs[0].Preview)

	return answer, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
