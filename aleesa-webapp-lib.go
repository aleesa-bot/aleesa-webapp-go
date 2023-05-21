package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
	"unsafe"

	"github.com/hjson/hjson-go"
	log "github.com/sirupsen/logrus"
)

// Горутинка, которая парсит json-чики прилетевшие из REDIS-ки
func msgParser(ctx context.Context, msg string) {
	var sendTo string
	var answer string
	var j rMsg
	var err error
	regexpWeather, err := regexp.Compile("(w|weather|погода|погодка|погадка)[[:space:]]+.+")

	if err != nil {
		log.Errorf("Unable to compile regexp for weather command: %s", err)
		return
	}

	log.Debugf("Incomming raw json: %s", msg)

	if err := json.Unmarshal([]byte(msg), &j); err != nil {
		log.Warnf("Unable to to parse message from redis channel: %s", err)
		return
	}

	j, err = validateRmsg(j, msg)

	if err != nil {
		log.Warn(err)
		return
	}

	// Если у нас циклическая пересылка сообщения, попробуем её тут разорвать, отбросив сообщение
	if j.Misc.Fwdcnt > config.ForwardsMax {
		log.Warnf("Discarding msg with fwd_cnt exceeding max value: %s", msg)
		return
	} else {
		j.Misc.Fwdcnt++
	}

	sendTo = j.Plugin

	// Классифицируем входящие сообщения. Первым делом, попробуем определить команды
	if j.Message[0:len(j.Misc.Csign)] == j.Misc.Csign {
		var cmd = j.Message[len(j.Misc.Csign):]

		switch {
		case cmd == "cat" || cmd == "кис":
			// Пытаемся вычитать ответ thecatapi 3 раза
			for i := 0; i < 3; i++ {
				answer, err = theCatAPIClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query api.thecatapi.com: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету кошечек, все разбежались."
			} else {
				// for telegram return answer in markdown format
				if j.Plugin == "telegram" {
					cats := []string{"龴ↀ◡ↀ龴", "=^..^=", "≧◔◡◔≦", "^ↀᴥↀ^"}
					j.Message = fmt.Sprintf("[%s](%s)", cats[rand.Intn(len(cats))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "fox" || cmd == "лис":
			// Пытаемся вычитать ответ randonfox.ca 3 раза
			for i := 0; i < 3; i++ {
				answer, err = randomFoxClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query randomfox.ca: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету лисичек, все разбежались."
			} else {
				if j.Plugin == "telegram" {
					j.Message = fmt.Sprintf("[-^^,--,~](%s)", answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "xkcd":
			// Пытаемся вычитать ответ xkcd.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = xkcdruClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query xkcd.ru: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Комикс-стрип нарисовать не так-то просто :("
			} else {
				if j.Plugin == "telegram" {
					j.Message = fmt.Sprintf("[xkcdru](%s)", answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "buni":
			// Пытаемся вычитать ответ www.bunicomic.com 3 раза
			for i := 0; i < 3; i++ {
				answer, err = buniComicClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query www.bunicomic.com: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету Buni :("
			} else {
				if j.Plugin == "telegram" {
					j.Message = fmt.Sprintf("[buni](%s)", answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "anek" || cmd == "анек" || cmd == "анекдот":
			// Пытаемся вычитать ответ www.anekdot.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = anekdotruClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query www.anekdot.ru: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Все рассказчики анекдотов отдыхают"
			} else {
				if j.Plugin == "telegram" {
					j.Message = fmt.Sprintf("```\n%s\n```", answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "monkeyuser":
			// Пытаемся вычитать ответ www.bunicomic.com 3 раза
			for i := 0; i < 3; i++ {
				answer, err = monkeyUserClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query monkeyuser.com: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету Monkey User-ов, они все спрятались."
			} else {
				if j.Plugin == "telegram" {
					j.Message = fmt.Sprintf("[monkeyuser](%s)", answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "butt" || cmd == "booty" || cmd == "ass" || cmd == "попа" || cmd == "попка":
			// Пытаемся вычитать ответ api.obutts.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = obuttsClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query api.obutts.ru: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету попок, все разбежались."
			} else {
				if j.Plugin == "telegram" {
					arts := []string{"(__(__)", "(_!_)", "(__.__)"}
					j.Message = fmt.Sprintf("[%s](%s)", arts[rand.Intn(len(arts))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "titts" || cmd == "boobs" || cmd == "tities" || cmd == "boobies" || cmd == "сиси" || cmd == "сисечки":
			// Пытаемся вычитать ответ api.obutts.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = oboobsClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query api.oboobs.ru: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету cисичек, все разбежались."
			} else {
				if j.Plugin == "telegram" {
					arts := []string{"(. )( .)", "(  . Y .  )", "(o)(o)", "( @ )( @ )", "(.)(.)"}
					j.Message = fmt.Sprintf("[%s](%s)", arts[rand.Intn(len(arts))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "drink" || cmd == "праздник":
			// Пытаемся вычитать ответ prazdniki-segodnya.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = psrClient()

				if err != nil {
					log.Errorf("Try %d/3 unable to query prazdniki-segodnya.ru: %v", i+1, err)
					time.Sleep(1 * time.Second)
					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Не знаю праздников - вджобываю весь день на шахтах, как проклятая."
			} else {
				j.Message = answer
			}
		case regexpWeather.Match([]byte(cmd)):
			if re, err := regexp.Compile("[[:space:]]+"); err != nil {
				log.Errorf("Unable to compile regex for weather command prsing: %s", err)
			} else {
				city := re.Split(cmd, 2)[1]

				if answer, err := owmClient(city); err != nil {
					log.Errorf("Unable to handle city %s in openweartermap api: %s", city, err)
					j.Message = fmt.Sprintf("Я не знаю, какая погода в %s", city)
				} else {
					j.Message = answer
				}
			}

		default:
			log.Errorf("Unknown command %s, unable to handle, skipping", j.Message)
			return
		}
	} else {
		log.Errorf("Message is not a command: %s, unable to handle, skipping", j.Message)
		return
	}

	// Настало время формировать json и засылать его в дальше
	var message sMsg
	message.From = j.From
	message.Userid = j.Userid
	message.Chatid = j.Chatid
	message.Threadid = j.Threadid
	message.Message = j.Message
	message.Plugin = j.Plugin
	message.Mode = j.Mode
	message.Misc.Fwdcnt = j.Misc.Fwdcnt
	message.Misc.Csign = j.Misc.Csign
	message.Misc.Username = j.Misc.Username
	message.Misc.Answer = j.Misc.Answer
	message.Misc.Botnick = j.Misc.Botnick
	message.Misc.Msgformat = j.Misc.Msgformat
	message.Misc.GoodMorning = j.Misc.GoodMorning

	data, err := json.Marshal(message)

	if err != nil {
		log.Warnf("Unable to to serialize message for redis: %s", err)
		return
	}

	// Заталкиваем наш json в редиску
	if err := redisClient.Publish(ctx, sendTo, data).Err(); err != nil {
		log.Warnf("Unable to send data to redis channel %s: %s", sendTo, err)
	} else {
		log.Debugf("Send msg to redis channel %s: %s", sendTo, string(data))
	}
}

// Читает и валидирует конфиг, а также выставляет некоторые default-ы, если значений для параметров в конфиге нет
func readConfig() {
	configLoaded := false
	executablePath, err := os.Executable()

	if err != nil {
		log.Errorf("Unable to get current executable path: %s", err)
	}

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	locations := []string{
		"~/.aleesa-webapp-go.json",
		"~/aleesa-webapp-go.json",
		"/etc/aleesa-webapp-go.json",
		configJSONPath,
	}

	for _, location := range locations {
		fileInfo, err := os.Stat(location)

		// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём
		if err != nil {
			continue
		}

		// Конфиг-файл длинноват для конфига, попробуем следующего кандидата
		if fileInfo.Size() > 65535 {
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
		var sampleConfig myConfig
		var tmp map[string]interface{}
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
			sampleConfig.Server = "localhost"
		}

		if sampleConfig.Port == 0 {
			sampleConfig.Port = 6379
		}

		if sampleConfig.Timeout == 0 {
			sampleConfig.Timeout = 10
		}

		if sampleConfig.Loglevel == "" {
			sampleConfig.Loglevel = "info"
		}

		// sampleConfig.Log = "" if not set

		if sampleConfig.Channel == "" {
			log.Errorf("Channel field in config file %s must be set", location)
			continue
		}

		if sampleConfig.DataDir == "" {
			sampleConfig.DataDir = "data"
		}

		if sampleConfig.Csign == "" {
			log.Errorf("Csign field in config file %s must be set", location)
			continue
		}

		if sampleConfig.ForwardsMax == 0 {
			sampleConfig.ForwardsMax = forwardMax
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

		config = sampleConfig
		configLoaded = true
		log.Infof("Using %s as config file", location)
		break
	}

	if !configLoaded {
		log.Error("Config was not loaded! Refusing to start.")
		os.Exit(1)
	}
}

// Хэндлер сигналов закрывает все бд и сваливает из приложения
func sigHandler() {
	var err error

	for {
		var s = <-sigChan
		switch s {
		case syscall.SIGINT:
			log.Infoln("Got SIGINT, quitting")
		case syscall.SIGTERM:
			log.Infoln("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			log.Infoln("Got SIGQUIT, quitting")

		// Заходим на новую итерацию, если у нас "неинтересный" сигнал
		default:
			continue
		}

		// Чтобы не срать в логи ошибками от редиски, проставим shutdown state приложения в true
		shutdown = true

		// Отпишемся от всех каналов и закроем коннект к редиске
		if err = subscriber.Unsubscribe(ctx); err != nil {
			log.Errorf("Unable to unsubscribe from redis channels cleanly: %s", err)
		}

		if err = subscriber.Close(); err != nil {
			log.Errorf("Unable to close redis connection cleanly: %s", err)
		}

		if len(pcacheDB) > 0 {
			log.Debug("Closing persistent cache db")

			for _, db := range pcacheDB {
				_ = db.Close()
			}
		}

		os.Exit(0)
	}
}

// Читает даденный файл построчно в массив строк
func readLines(path string) ([]string, error) {
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

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// Валидирует входящее сообщение
func validateRmsg(j rMsg, msg string) (rMsg, error) {
	if exist := j.From; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no from field: %s", msg)
		return j, errors.New(emsg)
	}

	if exist := j.Chatid; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no chatid field: %s", msg)
		return j, errors.New(emsg)
	}

	if exist := j.Userid; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no userid field: %s", msg)
		return j, errors.New(emsg)
	}

	// j.Threadid может быть пустым, значит либо нам его не дали, либо дали пустым. Это нормально.

	if exist := j.Message; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no message field: %s", msg)
		return j, errors.New(emsg)
	}

	if exist := j.Plugin; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no plugin field: %s", msg)
		return j, errors.New(emsg)
	}

	if exist := j.Mode; exist == "" {
		emsg := fmt.Sprintf("Incorrect msg from redis, no mode field: %s", msg)
		return j, errors.New(emsg)
	}

	// j.Misc.Answer может и не быть, тогда ответа на такое сообщение не будет
	// j.Misc.Botnick тоже можно не передавать, тогда будет записана пустая строка
	// j.Misc.Csign если нам его не передали, возьмём значение из конфига
	if exist := j.Misc.Csign; exist == "" {
		j.Misc.Csign = config.Csign
	}

	// j.Misc.Fwdcnt если нам его не передали, то будет 0
	if exist := j.Misc.Fwdcnt; exist == 0 {
		j.Misc.Fwdcnt = 1
	}

	// j.Misc.GoodMorning может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.Msgformat может быть быть 1 или 0, по-умолчанию 0
	// j.Misc.Username можно не передавать, тогда будет пустая строка
	return j, nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
