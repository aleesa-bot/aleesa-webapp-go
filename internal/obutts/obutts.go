package obutts

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/carlmjohnson/requests"
)

// Obutts структура, возвращаемая на запрос в obutts.ru.
type Obutts struct {
	ID      int    `json:"id,omitempty"`
	Author  string `json:"author,omitempty"`
	Rank    int    `json:"rank,omitempty"`
	Model   string `json:"model,omitempty"`
	Preview string `json:"preview"`
}

// APIClient клиент сервиса obutts.ru.
func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		url       = "http://api.obutts.ru/butts/0/1/random"
		butts     []Obutts
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&butts).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	if len(butts) == 0 {
		err = fmt.Errorf("empty json array returned from %s", url)

		return "", err
	}

	answer := "http://media.obutts.ru/" + butts[0].Preview

	return answer, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
