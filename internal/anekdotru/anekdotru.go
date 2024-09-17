package anekdotru

import (
	"aleesa-webapp-go/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"

	"github.com/carlmjohnson/requests"
)

func APIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx       = context.Background()
		respBody  bytes.Buffer
		respBytes []byte
		userAgent = cfg.UserAgents[rand.IntN(len(cfg.UserAgents))]
		url       = "https://www.anekdot.ru/rss/randomu.html"
	)

	err := requests.
		URL(url).
		UserAgent(userAgent).
		ToBytesBuffer(&respBody).
		Fetch(ctx)

	if err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	// Time to fix UTF-8, just in case.
	respBytes = bytes.ToValidUTF8(respBody.Bytes(), []byte{0xef, 0xbf, 0xbd})

	// Cut-off json.
	re := regexp.MustCompile(`JSON.parse\('`)
	aneks := re.Split(string(respBytes), 2)

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
				return fmt.Sprintf("RespBytes:%s\n\nAnekJSON: %s", respBytes, string(anekJSON)), err
			}

			re = regexp.MustCompile(`<br>`)
			answer := re.ReplaceAll([]byte(anekMap[0]), []byte("\n"))

			return string(answer), nil
		}

		fmt.Printf("a length: %d\nvalue: %v", len(a), a)
	} else {
		fmt.Printf("aneks length: %d\nvalue: %v", len(aneks), aneks)
	}

	err = errors.New("Unable to parse response from www.anekdot.ru: " + string(respBytes))

	return "", err
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
