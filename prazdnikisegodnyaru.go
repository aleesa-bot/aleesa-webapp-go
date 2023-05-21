package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

func psrClient() (string, error) {
	tsNow := time.Now()
	// Вычислим expiration time для MSK+0 в формате unix timestamp
	loc, _ := time.LoadLocation("Europe/Moscow")
	year, month, day := tsNow.In(loc).Date()
	eTime := time.Date(year, month, day, 0, 0, 0, 0, loc).AddDate(0, 0, +1).Unix()

	tsStr := getValue("cache", "prazdnikisegodhyaru_timestamp")

	// Кэш пустой, запросов пока не было
	if tsStr == "" {
		log.Debugf("PSR no cache entry for prazdnikisegodhyaru")
		answer, err := psrAPIClient()

		if err != nil {
			return "", err
		}

		if err := updatePsrCache(eTime, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	tsCache, err := strconv.ParseInt(tsStr, 10, 64)

	if err != nil {
		log.Debug("PSR cache miss")
		log.Warnf("PSR-Cache: unable to parse string as int64 (sha256 collision?): %s", err)
		answer, err := psrAPIClient()

		// Бинго! у нас ещё и ошибка при обращении к "api" prazdnikisegodhyaru, здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updatePsrCache(eTime, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	// Текущее время в Unix timestamp
	tsNowUnix := tsNow.Unix()

	// Кэш просрочен
	if tsCache < tsNowUnix {
		log.Debug("PSR cache miss")
		answer, err := psrAPIClient()

		// Здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updatePsrCache(eTime, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	// TODO: как-то научиться определять ситуацию с отсутствующим значением в кэше
	// TODO: Научиться определять ситуацию, когда значением является unix timestamp
	log.Debug("PSR cache hit")
	return getValue("cache", "prazdnikisegodhyaru_value"), nil
}

func psrAPIClient() (string, error) {
	var err error

	var c = http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, "https://prazdniki-segodnya.ru/", nil)
	if err != nil {
		return "", err
	}

	// Формальная попытка притвориться браузером
	_, _, day := time.Now().Date()
	req.Header.Set("User-Agent", userAgents[day])
	req.Header.Set("Accept-Language", "ru-RU")
	req.Header.Set("Accept-Charset", "utf-8")

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

	var holidays []string
	// Парсим body ответа от сервера
	tokenizer := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				// end of the file, break out of the loop
				break
			} else {
				return "", err
			}
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "div" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && attribute.Val == "list-group-item text-monospace" {
						tokenType = tokenizer.Next()

						if tokenType == html.TextToken {
							str := tokenizer.Token().Data
							holidays = append(holidays, strings.TrimSpace(str))
						}
					}
				}
			}
		}
	}

	if len(holidays) == 0 {
		err = errors.New("unable to parse response from prazdniki-segodnya.ru: no holidays found")
		return "", err
	}

	return "* " + strings.Join(holidays, "\n* "), nil
}

func updatePsrCache(eTime int64, value string) error {
	key := "prazdnikisegodhyaru_timestamp"
	timestamp := fmt.Sprintf("%d", eTime)

	if err := saveKeyWithValue("cache", key, timestamp); err != nil {
		return errors.New(fmt.Sprintf("PSR-Cache: unable to save timestamp to cache: %s", err))
	}

	key = "prazdnikisegodhyaru_value"

	if err := saveKeyWithValue("cache", key, value); err != nil {
		return errors.New(fmt.Sprintf("PSR-Cache: unable to save answer to cache: %s", err))
	}

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
