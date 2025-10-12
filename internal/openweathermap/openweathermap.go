package openweathermap

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"aleesa-webapp-go/internal/pcachedb"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
)

// Структурка для передачи в обрабтчик запросов в openweathermap api.
type OwmItem struct {
	City     string
	FakeName string
	Where    string
	Lon      float64
	Lat      float64
}

// OwmClient обёртка с кэшированием для openweatermap.org.
func OwmClient(cfg *config.MyConfig, city string, days uint32) (string, error) {
	var (
		answer string
		err    error
		item   OwmItem
	)

	if city == "" {
		return "Мне нужно ИМЯ города.", err
	}

	if len(city) > 80 {
		return "Длинновато для имени города", err
	}

	// "Нормализуем" название города.
	city = strings.TrimSpace(city)
	firstLetter := strings.SplitN(city, "", 2)[0]
	city = fmt.Sprintf("%s%s", strings.ToUpper(firstLetter), city[1:])
	city = strings.ToValidUTF8(city, "")

	switch city {
	case "Msk", "Default", "Dc", "Мск", "Dc-Universe":
		item.City = "Москва"
	case "Spb", "Спб":
		item.City = "Санкт-Петербург"
	case "Ект", "Ёбург", "Ебург", "Екат", "Ekt", "Eburg", "Ekat":
		item.City = "Екатеринбург"
	case "Кудыкина гора":
		// Координаты Килиманжаро (-3.066524619045895, 37.35588473583672).
		item.Lat = -3.066524619045895
		item.Lon = 37.35588473583672
		item.FakeName = "Килиманжаро"
		item.Where = "на горе"
	default:
		item.City = city
		item.Where = "в городе"
	}

	if days == 0 {
		answer, err = QueryOwmCache(cfg, item)
	}

	if err != nil {
		return "", err
	}

	return answer, nil
}

// OwmAPIClient делает запрос непосредственно в openweathermap.com api.
func OwmAPIClient(cfg *config.MyConfig, item OwmItem) (string, error) {
	var (
		answer   string
		cityName string
		client   = http.DefaultClient
	)

	w, err := owm.NewCurrent(
		"C",
		"ru",
		cfg.OpenWeatherMap.Appid,
		owm.WithHttpClient(client),
	)

	if err != nil {
		err := fmt.Errorf("OWM-Api: %w", err)

		return answer, err
	}

	if item.City != "" {
		if err = w.CurrentByName(item.City); err != nil {
			// TODO: extend error
			return answer, err
		}

		cityName = item.City
	} else {
		coord := owm.Coordinates{Longitude: item.Lon, Latitude: item.Lat}

		if err = w.CurrentByCoordinates(&coord); err != nil {
			// TODO: extend error
			return answer, err
		}

		if item.FakeName == "" {
			cityName = w.Name
		} else {
			cityName = item.FakeName
		}
	}

	// Тута парсер - направление ветра, градусы итд.
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

	if int(w.Main.TempMin) == int(w.Main.TempMax) {
		if cfg.OpenWeatherMap.Country {
			answer = fmt.Sprintf(
				"Погода %s %s, %s:\\n%s, ветер %s %.1f м/c, температура %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				item.Where,
				cityName,
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
				"Погода %s %s, (ш:%f, д:%f):\n%s, ветер %s %.1f м/c, температура %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				item.Where,
				cityName,
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
		if cfg.OpenWeatherMap.Country {
			answer = fmt.Sprintf(
				"Погода %s %s, %s:\n%s, ветер %s %.1f м/c, температура от %d до %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				item.Where,
				cityName,
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
				"Погода %s %s, (ш:%f, д:%f):\n%s, ветер %s %.1f м/c, температура от %d до %d°C, ощущается как %d°C, относительная влажность %d%%, давление %d мм.рт.ст",
				item.Where,
				cityName,
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

// UpdateOwmCache обновляет кэш.
func UpdateOwmCache(cfg *config.MyConfig, tsNowUnix int64, entry string, value string) error {
	key := fmt.Sprintf("%s+timestamp", entry)

	// "Кэшируем" на 3 часа.
	ts := tsNowUnix + int64((3 * time.Hour).Seconds())
	timestamp := fmt.Sprintf("%d", ts)

	if err := pcachedb.SaveKeyWithValue(cfg, "cache", key, timestamp); err != nil {
		return fmt.Errorf("OWM-Cache: unable to save timestamp to cache: %w", err)
	}

	key = fmt.Sprintf("%s+value", entry)

	if err := pcachedb.SaveKeyWithValue(cfg, "cache", key, value); err != nil {
		return fmt.Errorf("OWM-Cache: unable to save answer to cache: %w", err)
	}

	return nil
}

// QueryOwmCache пытается вынуть из кэша ответ относительно данного города или координат.
// Если кэш пуст, делает запрос к api owm.
func QueryOwmCache(cfg *config.MyConfig, item OwmItem) (string, error) {
	var (
		answer string
		err    error
		tsNow  = time.Now().Unix()
		key    string
	)

	if item.City != "" {
		key = item.City + "+timestamp"
	} else {
		key = fmt.Sprintf("%f+%f+timestamp", item.Lon, item.Lat)
	}

	// А вот тут мы реализовываем механику кэширования.
	tsCacheString := pcachedb.GetValue(cfg, "cache", key)

	// Для этого города в кэше пока ничего нету, вероятно это первый запрос.
	if tsCacheString == "" {
		log.Debugf("OWM no cache entry for city %s", key)

		answer, err := OwmAPIClient(cfg, item)

		// Здесь мы ничего сделать не можем.
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш.
		if err := UpdateOwmCache(cfg, tsNow, key, answer); err != nil {
			return "", err
		}

		return answer, err
	}

	tsCache, err := strconv.ParseInt(tsCacheString, 10, 64)

	// Если конверсия не удалась, значит значение в кэше хреновое, надо делать запрос в api и обновлять кэш.
	// Первое, что приходит в голову - коллизия хэшей sha256.
	if err != nil {
		log.Warnf("OWM-Cache: unable to parse string as int64 (sha256 collision?): %s", err)
		answer, err := OwmAPIClient(cfg, item)

		// Бинго! у нас ещё и ошибка при обращении к api openwethermap, здесь мы ничего сделать не можем.
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш.
		if err := UpdateOwmCache(cfg, tsNow, key, answer); err != nil {
			return "", err
		}

		return answer, err
	}

	// Кэш просрочен
	if tsCache < tsNow {
		log.Debugf("OWM cache miss for city: %s", key)

		answer, err := OwmAPIClient(cfg, item)

		// Здесь мы ничего сделать не можем.
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш.
		if err := UpdateOwmCache(cfg, tsNow, key, answer); err != nil {
			log.Errorf("%s", err)
		}

		return answer, err
	}

	key += "+value"
	answer = pcachedb.GetValue(cfg, "cache", key)

	// Какие-то проблемы с кэшом - ts есть, но самого кэша нету. Более подробная инфа по идее уже попала в лог.
	// Или в кэше что-то не то - длина ответа *точно* не может быть меньше 140 символов.
	// Попробуем достать данные из api и сохранить в тот самый многострадальный кэш.
	if len(answer) < 140 {
		log.Debugf("OWM cache miss for city: %s", key)

		if len(answer) > 0 {
			log.Errorf("OWM-Cache: record for city %s in cache is too short, less than 140 chars", key)
		}

		answer, err = OwmAPIClient(cfg, item)

		// Здесь мы ничего сделать не можем.
		if err != nil {
			return "", err
		}

		// Если всё хорошо, надо обновить кэш.
		if err := UpdateOwmCache(cfg, tsNow, key, answer); err != nil {
			return "", err
		}

		return answer, err
	}

	log.Debugf("OWM cache hit for city: %s", key)

	return answer, err
}

// Записывает в кэш, кто каким городом интересовался в последний раз, когда спрашивал о погоде.
// STUB
func UpdateOwmUserCache(cfg *config.MyConfig, chatid string, entry string, value string) error {
	var err error

	// Chatid как есть лучше на диск не пробовать класть, мало ли что туда юзер впишет. Возьмём от него sha256.
	data := []byte(chatid)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))

	if err = pcachedb.SaveKeyWithValue(cfg, hash, entry, value); err != nil {
		return fmt.Errorf("OWM-Cache: unable to save timestamp to cache: %w", err)
	}

	return err
}

// Возвращает из кэша город, которым интересовался пользователь из даденного чятика.
func QueryOwmUserCache(cfg *config.MyConfig, chatid string, entry string) string {
	// Chatid как есть лучше на диск не пробовать класть, мало ли что туда юзер впишет. Возьмём от него sha256.
	var (
		data  = []byte(chatid)
		hash  = fmt.Sprintf("%x", sha256.Sum256(data))
		value = pcachedb.GetValue(cfg, hash, entry)
	)

	return value
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
