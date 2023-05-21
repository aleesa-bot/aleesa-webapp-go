package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func xkcdruClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
		// Возвращаем ответ сразу же, не переходя по редиректам
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(http.MethodGet, "https://xkcd.ru/random/", nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	var resp *http.Response
	resp, err = c.Do(req)

	if err != nil {
		return "", err
	}

	// Response body нам не интересен, нам интересен статус и нам интересен заголовок Location
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")

		if len(location) <= 3 {
			err = errors.New("Location header too short")
			return "", err
		}

		location = location[1:(len(location) - 1)]

		return fmt.Sprintf("https://xkcd.ru/i/%s_v1.png", location), nil
	} else {
		err = errors.New(
			"resp.StatusCode: " +
				strconv.Itoa(resp.StatusCode))
		return "", err
	}

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
