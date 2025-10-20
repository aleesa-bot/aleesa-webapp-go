package randomfox

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"

	"github.com/carlmjohnson/requests"
)

// RandomFox структура, возвращаемая на запрос в randomfox.ca API.
type RandomFox struct {
	Image string `json:"image,omitempty"`
	Link  string `json:"link"`
}

// APIClient клиент сервиса randomfox.ca.
func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		url       = "https://randomfox.ca/floof/"
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		fox       RandomFox
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToJSON(&fox).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	var re = regexp.MustCompile(`\\`)

	answer := re.ReplaceAll([]byte(fox.Link), []byte(""))

	return string(answer), nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
