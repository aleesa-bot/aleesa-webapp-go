package xkcdru

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"
)

func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		err error
		url = "https://xkcd.ru/random/"
		c   = http.Client{
			Timeout: 10 * time.Second,
			// Возвращаем ответ сразу же, не переходя по редиректам.
			CheckRedirect: func(req *http.Request, via []*http.Request) error { //nolint: revive // Ну, не используется req и что из этого?
				return http.ErrUseLastResponse
			},
		}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil) //nolint: noctx

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", cfg.UserAgents[rand.IntN(len(cfg.UserAgents))])

	var resp *http.Response
	resp, err = c.Do(req) //nolint: bodyclose // Посмотри ниже, блять, тело запроса закрывается.

	if err != nil {
		return "", err
	}

	// Типа, надо закрывать Body в любом случае, как рекомендуют в документации https://pkg.go.dev/net/http .
	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Errorf("Unable to close response body for request to %s: %s", url, err)
		}
	}(resp.Body)

	// Response body нам не интересен, нам интересен статус и нам интересен заголовок Location
	if resp.StatusCode == 302 {
		location := resp.Header.Get("Location")

		if len(location) <= 3 {
			err = fmt.Errorf("location header too short")

			return "", err
		}

		location = location[1:(len(location) - 1)]

		return fmt.Sprintf("https://xkcd.ru/i/%s_v1.png", location), nil
	}

	err = fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
