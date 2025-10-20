package monkeyuser

import (
	"aleesa-webapp-go/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
)

// MonkeyUsers описывает структуру json-массива ответа www.monkeyusers.com.
type MonkeyUsers []struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// APIClient клиент сервиса monkeyuser.com.
func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		respBody  bytes.Buffer
		respBytes []byte
		body      string
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		baseURL   = "https://www.monkeyuser.com"
		indexURL  = baseURL + "/index.json"
	)

	err := requests.
		URL(indexURL).
		UserAgent(userAgent).
		ToBytesBuffer(&respBody).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", indexURL, err)
	}

	// На всякий случай валидируем utf-8.
	respBytes = bytes.ToValidUTF8(respBody.Bytes(), []byte{0xef, 0xbf, 0xbd})

	var monkeyusers MonkeyUsers

	if err := json.Unmarshal(respBytes, &monkeyusers); err != nil {
		err = fmt.Errorf("unable to to parse response from %s: %w", indexURL, err)

		return "", err
	}

	amountOfUsers := len(monkeyusers)

	if amountOfUsers == 0 {
		err = fmt.Errorf("no links found in %s", indexURL)

		return "", err
	}

	ofChoice := rand.IntN(amountOfUsers)
	pageURL := fmt.Sprintf("%s%s", baseURL, monkeyusers[ofChoice].URL)

	// Делаем второй запрос к серверу.
	err = requests.
		URL(pageURL).
		UserAgent(userAgent).
		ToString(&body).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", pageURL, err)
	}

	// Парсим body ответа от сервера
	doc, err := html.Parse(strings.NewReader(body))

	if err != nil {
		return "", err
	}

	var (
		hyperlink string
		link      func(*html.Node)
	)

	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a_href := range n.Attr { //nolint: revive
				// Если у нас class == "lazyload small-image", то после него ключ data-src должен содержать
				// искомый relative url
				if a_href.Key == "class" && a_href.Val == "logo" {
					continue
				}

				var v = a_href.Val
				if a_href.Key == "src" && v != "/images/logo.png" {
					hyperlink = v

					break
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link(c)
		}
	}

	link(doc)

	// Возвращаем ссылку на картинку buni comic
	if hyperlink != "" {
		return fmt.Sprintf("%s%s", baseURL, hyperlink), nil
	}

	err = errors.New("no links found on monkeyusers.com")

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
