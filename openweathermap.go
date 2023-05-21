package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	log "github.com/sirupsen/logrus"
)

// Обёртка с кэшированием
func owmClient(city string) (string, error) {
	var answer string

	if city == "" {
		return "Мне нужно ИМЯ города.", nil
	}

	if len(city) > 80 {
		return "Длинновато для имени города", nil
	}

	// "Нормализуем" название города.
	city = strings.TrimSpace(city)
	firstLetter := strings.SplitN(city, "", 2)[0]
	city = fmt.Sprintf("%s%s", strings.ToUpper(firstLetter), city[1:])
	city = strings.ToValidUTF8(city, "")

	switch city {
	case "Msk", "Default", "Dc", "Мск", "Dc-Universe":
		city = "Москва"
	case "Spb", "Спб":
		city = "Санкт-Петербург"
	case "Ект", "Ёбург", "Ебург", "Екат", "Ekt", "Eburg", "Ekat":
		city = "Екатеринбург"
	}

	// А вот тут мы реализовываем механику кэширования
	tsNow := time.Now().Unix()
	key := city + "+timestamp"

	tsCacheString := getValue("cache", key)

	// Для этого города в кэше пока ничего нету, вероятно это первый запрос
	if tsCacheString == "" {
		log.Debugf("OWM no cache entry for city %s", city)
		answer, err := owmAPIClient(city)

		// Здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updateOwmCache(tsNow, city, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	tsCache, err := strconv.ParseInt(tsCacheString, 10, 64)

	// Если конверсия не удалась, значит значение в кэше хреновое, надо делать запрос в api и обновлять кэш
	// Первое, что приходит в голову - коллизия хэшей sha256
	if err != nil {
		log.Warnf("OWM-Cache: unable to parse string as int64 (sha256 collision?): %s", err)
		answer, err := owmAPIClient(city)

		// Бинго! у нас ещё и ошибка при обращении к api openwethermap, здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updateOwmCache(tsNow, city, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	// Кэш просрочен
	if tsCache < tsNow {
		log.Debugf("OWM cache miss for city: %s", city)
		answer, err := owmAPIClient(city)

		// Здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updateOwmCache(tsNow, city, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	key = city + "+value"
	answer = getValue("cache", key)

	// Какие-то проблемы с кэшом - ts есть, но самого кэша нету. Более подробная инфа по идее уже попала в лог.
	// Или в кэше что-то не то - длина ответа *точно* не может быть меньше 140 символов.
	// Попробуем достать данные из api и сохранить в тот самый многострадальный кэш
	if len(answer) < 140 {
		log.Debugf("OWM cache miss for city: %s", city)

		if len(answer) > 0 {
			log.Errorf("OWM-Cache: record for city %s in cache is too short, less than 140 chars", city)
		}

		answer, err = owmAPIClient(city)

		// Здесь мы ничего сделать не можем
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш
		if err := updateOwmCache(tsNow, city, answer); err != nil {
			log.Error(err)
		}

		return answer, nil
	}

	log.Debugf("OWM cache hit for city: %s", city)

	return answer, nil
}

// Делает запрос непосредственно в openweathermap.com api
func owmAPIClient(city string) (string, error) {
	var answer string

	w, err := owm.NewCurrent("C", "ru", config.OpenWeatherMap.Appid)

	if err != nil {
		err := errors.New(fmt.Sprintf("OWM-Api: %s", err))
		return answer, err
	}

	if err = w.CurrentByName(city); err != nil {
		// TODO: extend error
		return answer, err
	}
	// тута парсер - направление ветра, градусы итд

	wind := "разный"

	switch {
	case w.Wind.Deg == 0:
		wind = "северный"
	case w.Wind.Deg > 0 && w.Wind.Deg <= 30:
		wind = "северо-северо-восточный"
	case w.Wind.Deg > 30 && w.Wind.Deg <= 60:
		wind = "северо-восточный"
	case w.Wind.Deg > 60 && w.Wind.Deg < 90:
		wind = "восточно-северо-восточный"
	case w.Wind.Deg == 90:
		wind = "восточный"
	case w.Wind.Deg > 90 && w.Wind.Deg <= 120:
		wind = "восточно-юго-восточный"
	case w.Wind.Deg > 120 && w.Wind.Deg <= 150:
		wind = "юговосточный"
	case w.Wind.Deg > 150 && w.Wind.Deg < 180:
		wind = "юго-юго-восточный"
	case w.Wind.Deg == 180:
		wind = "южный"
	case w.Wind.Deg > 180 && w.Wind.Deg <= 210:
		wind = "юго-юго-западный"
	case w.Wind.Deg > 210 && w.Wind.Deg <= 240:
		wind = "юго-западный"
	case w.Wind.Deg > 240 && w.Wind.Deg < 270:
		wind = "западно-юго-западный"
	case w.Wind.Deg == 270:
		wind = "западный"
	case w.Wind.Deg > 270 && w.Wind.Deg <= 300:
		wind = "западно-северо-западный"
	case w.Wind.Deg > 300 && w.Wind.Deg <= 330:
		wind = "северо-западный"
	case w.Wind.Deg > 330 && w.Wind.Deg < 360:
		wind = "северо-северо-западный"
	case w.Wind.Deg == 360:
		wind = "северный"
	}

	if w.Main.TempMin == w.Main.TempMax {
		if config.OpenWeatherMap.Country {
			answer = fmt.Sprintf(
				"Погода в городе %s, %s:\\n%s, ветер %s %.1f м/c, температура %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				city,
				w.Sys.Country,
				w.Weather[0].Description,
				wind,
				w.Wind.Speed,
				int(w.Main.TempMin),
				int(w.Main.FeelsLike),
				w.Main.Humidity,
				int(w.Main.Pressure*0.75006375541921),
			)
		} else {
			answer = fmt.Sprintf(
				"Погода в городе %s, (ш:%f, д:%f):\n%s, ветер %s %.1f м/c, температура %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				city,
				w.GeoPos.Latitude,
				w.GeoPos.Longitude,
				w.Weather[0].Description,
				wind,
				w.Wind.Speed,
				int(w.Main.TempMin),
				int(w.Main.FeelsLike),
				w.Main.Humidity,
				int(w.Main.Pressure*0.75006375541921),
			)
		}
	} else {
		if config.OpenWeatherMap.Country {
			answer = fmt.Sprintf(
				"Погода в городе %s, %s:\n%s, ветер %s %.1f м/c, температура от %d до %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				city,
				w.Sys.Country,
				w.Weather[0].Description,
				wind,
				w.Wind.Speed,
				int(w.Main.TempMin),
				int(w.Main.TempMax),
				int(w.Main.FeelsLike),
				w.Main.Humidity,
				int(w.Main.Pressure*0.75006375541921),
			)
		} else {
			answer = fmt.Sprintf(
				"Погода в городе %s, (ш:%f, д:%f):\n%s, ветер %s %.1f м/c, температура от %d до %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				city,
				w.GeoPos.Latitude,
				w.GeoPos.Longitude,
				w.Weather[0].Description,
				wind,
				w.Wind.Speed,
				int(w.Main.TempMin),
				int(w.Main.TempMax),
				int(w.Main.FeelsLike),
				w.Main.Humidity,
				int(w.Main.Pressure*0.75006375541921),
			)
		}
	}

	return answer, err
}

// Обновляет кэш
func updateOwmCache(tsNowUnix int64, city string, value string) error {
	key := fmt.Sprintf("%s+timestamp", city)

	// "Кэшируем" на 3 часа
	ts := tsNowUnix + int64((3 * time.Hour).Seconds())
	timestamp := fmt.Sprintf("%d", ts)

	if err := saveKeyWithValue("cache", key, timestamp); err != nil {
		return errors.New(fmt.Sprintf("OWM-Cache: unable to save timestamp to cache: %s", err))
	}

	key = fmt.Sprintf("%s+value", city)

	if err := saveKeyWithValue("cache", key, value); err != nil {
		return errors.New(fmt.Sprintf("OWM-Cache: unable to save answer to cache: %s", err))
	}

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
