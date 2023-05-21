package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func monkeyUserClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.monkeyuser.com/toc/", nil)
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
			log.Errorf("Unable to close response body for thecatapi request: %s", err)
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

	var hyperlinks []string
	var link func(*html.Node)

	link = func(n *html.Node) {
		myFlag := false
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a_href := range n.Attr {
				// Если у нас class == "lazyload small-image", то после него ключ data-src должен содержать
				// искомый relative url
				if a_href.Key == "class" && a_href.Val == "lazyload small-image" {
					myFlag = true
					continue
				}
				if myFlag == true {
					if a_href.Key == "data-src" {
						myFlag = false
						hyperlinks = append(hyperlinks, a_href.Val)
					}
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
	if hyperlinks != nil {
		monkeyuser := hyperlinks[rand.Intn(len(hyperlinks))]
		monkeyuser = fmt.Sprintf("https://www.monkeyuser.com%s", monkeyuser)
		return monkeyuser, nil
	}

	err = errors.New("No links found on monkeyusers.com")
	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
