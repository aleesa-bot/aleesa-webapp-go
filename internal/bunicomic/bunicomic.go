package bunicomic

import (
	"aleesa-webapp-go/internal/config"
	"context"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strings"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
)

func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		body      string
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		url       = "https://www.bunicomic.com/?random&nocache=1"
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToString(&body).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	// Парсим body ответа от сервера.
	doc, err := html.Parse(strings.NewReader(body))

	if err != nil {
		return "", err
	}

	var (
		metaTags []string
		meta     func(*html.Node)
	)

	meta = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, meta := range n.Attr {
				if meta.Key == "content" {
					// adds a new link entry when the attribute matches
					metaTags = append(metaTags, meta.Val)
				}
			}
		}

		// traverses the HTML of the webpage from the first child node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			meta(c)
		}
	}
	meta(doc)

	// Возвращаем ссылку на картинку buni comic.
	if len(metaTags) != 0 {
		// А теперь среди полученнного мусора найдём ссылку на урл.
		for _, buni := range metaTags {
			match, err := regexp.MatchString(`wp-content/uploads`, buni) //nolint: staticcheck

			if err != nil {
				return "", err
			}

			if match {
				return buni, nil
			}
		}
	}

	// Картинки не выпарсилось, это точно ошибка.
	err = fmt.Errorf("unable to get link to buni comic strip image")

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
