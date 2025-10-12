package config

import (
	"aleesa-webapp-go/internal/defaults"
	"aleesa-webapp-go/internal/log"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/cockroachdb/pebble"
	"github.com/go-redis/redis/v8"
	"github.com/hjson/hjson-go"
)

// MyConfig структурка, описывающая конфиг.
type MyConfig struct {
	Server         string `json:"server,omitempty"`
	Port           int    `json:"port,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
	Loglevel       string `json:"loglevel,omitempty"`
	Log            string `json:"log,omitempty"`
	Channel        string `json:"channel,omitempty"`
	DataDir        string `json:"datadir,omitempty"`
	Csign          string `json:"csign,omitempty"`
	ForwardsMax    int64  `json:"forwards_max,omitempty"`
	OpenWeatherMap struct {
		Enabled bool   `json:"enabled,omitempty"`
		Country bool   `json:"country,omitempty"`
		Appid   string `json:"appid,omitempty"`
	} `json:"openweathermap,omitempty"`
	Flickr struct {
		Enabled          bool   `json:"enabled,omitempty"`
		Key              string `json:"key,omitempty"`
		Secret           string `json:"secret,omitempty"`
		OAuthToken       string `json:"oauth_token,omitempty"`
		OAuthTokenSecret string `json:"oauth_token_secret,omitempty"`
	} `json:"flickr,omitempty"`
	UserAgents  []string `json:"user_agents,omitempty"`
	PcacheDB    map[string]*pebble.DB
	RedisClient *redis.Client
}

// ReadConfig читает и валидирует конфиг, а также выставляет некоторые default-ы, если значений для параметров в конфиге
// нет.
func ReadConfig() (*MyConfig, error) {
	locations, err := defaults.DefaultConfigFileLocations()

	if err != nil {
		return &MyConfig{}, err
	}

	for _, location := range locations {
		fileInfo, err := os.Stat(location)

		// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём
		if err != nil {
			continue
		}

		// Конфиг-файл длинноват для конфига, попробуем следующего кандидата
		if fileInfo.Size() > defaults.ConfigFileSize {
			log.Warnf("Config file %s is too long for config, skipping", location)

			continue
		}

		buf, err := os.ReadFile(location)

		// Не удалось прочитать, попробуем следующего кандидата
		if err != nil {
			log.Warnf("Skip reading config file %s: %s", location, err)

			continue
		}

		// Исходя из документации, hjson какбы умеет парсить "кривой" json, но парсит его в map-ку.
		// Интереснее на выходе получить структурку: то есть мы вначале конфиг преобразуем в map-ку, затем эту map-ку
		// сериализуем в json, а потом json преврщааем в стркутурку. Не очень эффективно, но он и не часто требуется.
		var (
			sampleConfig *MyConfig
			tmp          map[string]interface{}
		)

		err = hjson.Unmarshal(buf, &tmp)

		// Не удалось распарсить - попробуем следующего кандидата
		if err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		tmpjson, err := json.Marshal(tmp)

		// Не удалось преобразовать map-ку в json
		if err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		if err := json.Unmarshal(tmpjson, &sampleConfig); err != nil {
			log.Warnf("Skip parsing config file %s: %s", location, err)

			continue
		}

		// Валидируем значения из конфига
		if sampleConfig.Server == "" {
			sampleConfig.Server = defaults.Host
		}

		if sampleConfig.Port == 0 {
			sampleConfig.Port = defaults.RedisPort
		}

		if sampleConfig.Timeout == 0 {
			sampleConfig.Timeout = defaults.NetworkTimeout
		}

		if sampleConfig.Loglevel == "" {
			sampleConfig.Loglevel = defaults.Loglevel
		}

		// sampleConfig.Log = "" if not set

		if sampleConfig.Channel == "" {
			log.Errorf("Channel field in config file %s must be set", location)

			continue
		}

		if sampleConfig.DataDir == "" {
			sampleConfig.DataDir = defaults.DataDir
		}

		if sampleConfig.Csign == "" {
			log.Errorf("Csign field in config file %s must be set", location)

			continue
		}

		if sampleConfig.ForwardsMax == 0 {
			sampleConfig.ForwardsMax = defaults.ForwardMax
		}

		if unsafe.Sizeof(sampleConfig.OpenWeatherMap) == 0 {
			log.Errorf("Hash OpenWeatherMap in config file %s must be defined", location)

			continue
		}

		if sampleConfig.OpenWeatherMap.Appid == "" {
			log.Errorf("Openweathermap->appid in config file %s must be set to its valid value", location)

			continue
		}

		// sampleConfig.OpenWeatherMap.Country is false by default

		// Теперь надо бы попробовать загрузить список User Agent-ов, с которыми ходить в разные апишки.
		uaFile := fmt.Sprintf("%s/useragents.txt", sampleConfig.DataDir)
		uaList, err := ReadLines(uaFile)

		if err != nil {
			return &MyConfig{}, fmt.Errorf("unable to read %s: %w", uaFile, err)
		}

		sampleConfig.UserAgents = uaList

		sampleConfig.PcacheDB = make(map[string]*pebble.DB)

		log.Infof("Using %s as config file", location)

		return sampleConfig, nil
	}

	return &MyConfig{}, fmt.Errorf("config was not loaded! Refusing to start")
}

// ReadLines читает даденный файл построчно в массив строк.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Unable to close file %s:%s", path, err)
		}
	}(file)

	var (
		lines   []string
		scanner = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
