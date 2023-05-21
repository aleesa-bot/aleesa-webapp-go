package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func randomFoxClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://randomfox.ca/floof/", nil)
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
			log.Errorf("Unable to close response body for randomfox.ca api request: %s", err)
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
