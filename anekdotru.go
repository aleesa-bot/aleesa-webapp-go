package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

func anekdotruClient() (string, error) {
	var (
		err error
		c   = http.Client{
			Timeout: 10 * time.Second,
		}
		URL = "https://www.anekdot.ru/rss/randomu.html"
	)

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])

	var resp *http.Response
	resp, err = c.Do(req) //nolint: bodyclose

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		if err != nil {
			log.Errorf("Unable to close response body to %s request: %s", URL, err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response code to %s is not 200: %d", URL, resp.StatusCode)

		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	// Time to fix UTF-8, just in case
	respBody = bytes.ToValidUTF8(respBody, []byte{0xef, 0xbf, 0xbd})

	re := regexp.MustCompile(`JSON.parse\('`)
	aneks := re.Split(string(respBody), 2)

	if (len(aneks) >= 2) && (aneks[1] != "") {
		re = regexp.MustCompile(`'\);`)
		a := re.Split(aneks[1], 2)

		if (len(a) >= 2) && (a[0] != "") {
			anek := []byte(a[0])

			re = regexp.MustCompile(`\\"`)
			anek = re.ReplaceAll(anek, []byte(`"`))

			re = regexp.MustCompile(`\\"`)
			anekJSON := re.ReplaceAll(anek, []byte(`"`))
			anekJSON = bytes.ToValidUTF8(anekJSON, []byte{0xef, 0xbf, 0xbd})

			var anekMap []string
			if err := json.Unmarshal(anekJSON, &anekMap); err != nil {
				return fmt.Sprintf("RespBody:%s\n\nAnekJSON: %s", respBody, string(anekJSON)), err
			}

			re = regexp.MustCompile(`<br>`)
			answer := re.ReplaceAll([]byte(anekMap[0]), []byte("\n"))

			return string(answer), nil
		}

		fmt.Printf("a length: %d\nvalue: %v", len(a), a)
	} else {
		fmt.Printf("aneks length: %d\nvalue: %v", len(aneks), aneks)
	}

	err = errors.New("Unable to parse response from www.anekdot.ru: " + string(respBody))

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
