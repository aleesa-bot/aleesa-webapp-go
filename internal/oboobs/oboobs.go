package oboobs

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/carlmjohnson/requests"
)

// Oboobs структура, возвращаемая на запрос в oboobs.ru.
type Oboobs struct {
	ID      int    `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Rank    int    `json:"rank,omitempty"`
	Model   string `json:"model,omitempty"`
	Preview string `json:"preview"`
}

func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		url       = "http://api.oboobs.ru/boobs/0/1/random"
		boobs     []Oboobs
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&boobs).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	if len(boobs) == 0 {
		err = fmt.Errorf("empty json array returned from %s", url)

		return "", err
	}

	answer := fmt.Sprintf("https://media.oboobs.ru/%s", boobs[0].Preview)

	return answer, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
