package prazdnikisegodnyaru

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"aleesa-webapp-go/internal/pcachedb"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"golang.org/x/net/html"
)

// PsrClient обёртка для клиента сервиса prazdniki-segodnya.ru, работающая с кэшем. Запросы надо делать именно через неё.
func PsrClient(cfg *config.MyConfig) (string, error) {
	tsNow := time.Now()

	// Вычислим expiration time для MSK+0 в формате unix timestamp
	loc, _ := time.LoadLocation("Europe/Moscow")
	year, month, day := tsNow.In(loc).Date()
	eTime := time.Date(year, month, day, 0, 0, 0, 0, loc).AddDate(0, 0, +1).Unix()

	tsStr := pcachedb.GetValue(cfg, "cache", "prazdnikisegodhyaru_timestamp")

	// Кэш пустой, запросов пока не было
	if tsStr == "" {
		log.Debugf("PSR no cache entry for prazdnikisegodhyaru")

		answer, err := PsrAPIClient(cfg)

		if err != nil {
			return "", err
		}

		if err := UpdatePsrCache(cfg, eTime, answer); err != nil {
			log.Errorf("%s", err)
		}

		return answer, nil
	}

	tsCache, err := strconv.ParseInt(tsStr, 10, 64)

	if err != nil {
		log.Debug("PSR cache miss")
		log.Warnf("PSR-Cache: unable to parse string as int64 (sha256 collision?): %s", err)
		answer, err := PsrAPIClient(cfg)

		// Бинго! у нас ещё и ошибка при обращении к "api" prazdnikisegodhyaru, здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := UpdatePsrCache(cfg, eTime, answer); err != nil {
			log.Errorf("%s", err)
		}

		return answer, nil
	}

	// Текущее время в Unix timestamp
	tsNowUnix := tsNow.Unix()

	// Кэш просрочен
	if tsCache < tsNowUnix {
		log.Debug("PSR cache miss")

		answer, err := PsrAPIClient(cfg)

		// Здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := UpdatePsrCache(cfg, eTime, answer); err != nil {
			log.Errorf("%s", err)
		}

		return answer, nil
	}

	// TODO: как-то научиться определять ситуацию с отсутствующим значением в кэше
	// TODO: Научиться определять ситуацию, когда значением является unix timestamp
	log.Debug("PSR cache hit")

	return pcachedb.GetValue(cfg, "cache", "prazdnikisegodhyaru_value"), nil
}

// PsrAPIClient клиент сервиса prazdniki-segodnya.ru.
func PsrAPIClient(cfg *config.MyConfig) (string, error) {
	var (
		ctx      = context.Background()
		url      = "https://prazdniki-segodnya.ru/"
		respBody string
	)

	// Формальная попытка притвориться браузером
	_, _, day := time.Now().Date()
	userAgent := cfg.UserAgents[day]

	if err := requests.URL(url).UserAgent(userAgent).Header("User-Agent", "ru-RU").
		Header("Accept-Charset", "utf-8").ToString(&respBody).Fetch(ctx); err != nil {
		return "", fmt.Errorf("unable to GET %s: %w", url, err)
	}

	var holidays []string

	// Парсим body ответа от сервера
	tokenizer := html.NewTokenizer(strings.NewReader(respBody))

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()

			if errors.Is(err, io.EOF) {
				// end of the file, break out of the loop
				break
			}

			return "", err
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
		return "", errors.New("unable to parse response from prazdniki-segodnya.ru: no holidays found")
	}

	return "* " + strings.Join(holidays, "\n* "), nil
}

// UpdatePsrCache сохраняет в кэш ответ от сервиса prazdniki-segodnya.ru и временной штамп.
func UpdatePsrCache(cfg *config.MyConfig, eTime int64, value string) error {
	key := "prazdnikisegodhyaru_timestamp"
	timestamp := strconv.FormatInt(eTime, 10)

	if err := pcachedb.SaveKeyWithValue(cfg, "cache", key, timestamp); err != nil {
		return fmt.Errorf("PSR-Cache: unable to save timestamp to cache: %w", err)
	}

	key = "prazdnikisegodhyaru_value"

	if err := pcachedb.SaveKeyWithValue(cfg, "cache", key, value); err != nil {
		return fmt.Errorf("PSR-Cache: unable to save answer to cache: %w", err)
	}

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
