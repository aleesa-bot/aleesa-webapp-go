package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func monkeyUserClient() (string, error) {
	var (
		err error
		c   = http.Client{
			Timeout: 10 * time.Second,
		}

		baseURL   = "https://www.monkeyuser.com"
		indexURL  = fmt.Sprintf("%s/index.json", baseURL)
		userAgent = userAgents[rand.Intn(len(userAgents))]
		resp      *http.Response
	)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(http.MethodGet, indexURL, nil) //nolint: noctx

	if err != nil {
		return "", err
	}

	// Притворяемся браузером.
	req.Header.Set("User-Agent", userAgent)

	resp, err = c.Do(req) //nolint: bodyclose

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Errorf("Unable to close response body for monkeyuser request: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("request to %s failed: %s", indexURL, resp.Status)

		return "", err
	}

	text, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	// Time to fix UTF-8, just in case
	text = bytes.ToValidUTF8(text, []byte{0xef, 0xbf, 0xbd})

	var monkeyusers monkeyUsers

	if err := json.Unmarshal(text, &monkeyusers); err != nil {
		err = fmt.Errorf("unable to to parse response from %s: %w", indexURL, err)

		return "", err
	}

	amountOfUsers := len(monkeyusers)

	if amountOfUsers == 0 {
		err = fmt.Errorf("no links found in %s", indexURL)

		return "", err
	}

	ofChoice := rand.Intn(amountOfUsers)
	pageURL := fmt.Sprintf("%s%s", baseURL, monkeyusers[ofChoice].URL)

	// Делаем второй запрос к серверу.
	req, err = http.NewRequest(http.MethodGet, pageURL, nil) //nolint: noctx

	if err != nil {
		return "", err
	}

	// Притворяемся тем же браузером.
	req.Header.Set("User-Agent", userAgent)

	resp, err = c.Do(req) //nolint: bodyclose

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
		err = fmt.Errorf("request to %s failed: %s", indexURL, resp.Status)

		return "", err
	}

	// Парсим body ответа от сервера
	doc, err := html.Parse(resp.Body)

	if err != nil {
		return "", err
	}

	var (
		hyperlink string
		link      func(*html.Node)
	)

	link = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a_href := range n.Attr { //nolint: revive,stylecheck
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

	err = fmt.Errorf("no links found on monkeyusers.com")

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
