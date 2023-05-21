package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func oboobsClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "http://api.oboobs.ru/boobs/0/1/random", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	var resp *http.Response
	resp, err = c.Do(req)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Errorf("Unable to close response body for thecatapi request: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = errors.New(
			"resp.StatusCode: " +
				strconv.Itoa(resp.StatusCode))
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
		err = errors.New("Empty json array returned from api.oboobs.ru")
		return "", err
	}

	answer := fmt.Sprintf("https://media.oboobs.ru/%s", boobs[0].Preview)
	return answer, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
