package thecatapi

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/carlmjohnson/requests"
)

// TheCatAPIinnerStruct структурка, описывающая элемент массива, возвращемый Тhe cat api.
type TheCatAPIinnerStruct struct {
	ID     string `json:"id,omitempty"`
	URL    string `json:"url"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// TheCatAPI структурка, описывающая массив, возвращаемый на запрос урла картинки из The Cat API.
type TheCatAPI []TheCatAPIinnerStruct

// APIClient клиент сервиса thecatapi.com.
func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		url       = "https://api.thecatapi.com/v1/images/search"
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		cat       TheCatAPI
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&cat).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	return cat[0].URL, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
