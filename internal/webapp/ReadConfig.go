package webapp

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/cockroachdb/pebble"
	"github.com/hjson/hjson-go"
)

// ReadConfig читает и валидирует конфиг, а также выставляет некоторые default-ы, если значений для параметров в конфиге
// нет.
func ReadConfig() error {
	locations, err := DefaultConfigFileLocations()

	if err != nil {
		return err
	}

	for _, location := range locations {
		fileInfo, err := os.Stat(location)

		// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём
		if err != nil {
			continue
		}

		// Конфиг-файл длинноват для конфига, попробуем следующего кандидата
		if fileInfo.Size() > ConfigFileSize {
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
			sampleConfig *config.MyConfig
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
			sampleConfig.Server = Host
		}

		if sampleConfig.Port == 0 {
			sampleConfig.Port = RedisPort
		}

		if sampleConfig.Timeout == 0 {
			sampleConfig.Timeout = NetworkTimeout
		}

		if sampleConfig.Loglevel == "" {
			sampleConfig.Loglevel = Loglevel
		}

		// sampleConfig.Log = "" if not set

		if sampleConfig.Channel == "" {
			log.Errorf("Channel field in config file %s must be set", location)

			continue
		}

		if sampleConfig.DataDir == "" {
			sampleConfig.DataDir = DataDir
		}

		if sampleConfig.Csign == "" {
			log.Errorf("Csign field in config file %s must be set", location)

			continue
		}

		if sampleConfig.ForwardsMax == 0 {
			sampleConfig.ForwardsMax = ForwardMax
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
		uaFile := sampleConfig.DataDir + "/useragents.txt"
		uaList, err := ReadLines(uaFile)

		if err != nil {
			return fmt.Errorf("unable to read %s: %w", uaFile, err)
		}

		sampleConfig.UserAgents = uaList

		sampleConfig.PcacheDB = make(map[string]*pebble.DB)

		log.Infof("Using %s as config file", location)

		Config = sampleConfig

		return nil
	}

	return errors.New("config was not loaded! Refusing to start")
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
