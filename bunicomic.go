package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func buniComicClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "http://www.bunicomic.com/?random&nocache=1", nil)

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
			log.Errorf("Unable to close response body for bunicomic request: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = errors.New(
			"resp.StatusCode: " +
				strconv.Itoa(resp.StatusCode))
		return "", err
	}

	// Парсим body ответа от сервера
	doc, err := html.Parse(resp.Body)

	if err != nil {
		return "", err
	}

	var metaTags []string
	var meta func(*html.Node)

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

	// Возвращаем ссылку на картинку buni comic
	if metaTags != nil {
		// А теперь среди полученнного мусора найдём ссылку на урл
		for _, buni := range metaTags {
			match, err := regexp.MatchString(`wp\-content/uploads`, buni)

			if err != nil {
				return "", err
			}

			if match {
				return buni, nil
			}
		}
	}

	// Картинки не выпарсилось, это точно ошибка
	err = errors.New("Unable to get link to buni comic strip image")
	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
