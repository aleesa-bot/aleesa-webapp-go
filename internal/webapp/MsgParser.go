package webapp

import (
	"aleesa-webapp-go/internal/anekdotru"
	"aleesa-webapp-go/internal/bunicomic"
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/flickr"
	"aleesa-webapp-go/internal/log"
	"aleesa-webapp-go/internal/monkeyuser"
	"aleesa-webapp-go/internal/oboobs"
	"aleesa-webapp-go/internal/obutts"
	"aleesa-webapp-go/internal/openweathermap"
	"aleesa-webapp-go/internal/prazdnikisegodnyaru"
	"aleesa-webapp-go/internal/randomfox"
	"aleesa-webapp-go/internal/thecatapi"
	"aleesa-webapp-go/internal/xkcdru"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"regexp"
	"time"
)

// MsgParser горутинка, которая парсит json-чики прилетевшие из REDIS-ки.
func MsgParser(cfg *config.MyConfig, ctx context.Context, msg string) { //nolint: revive
	var (
		sendTo string
		answer string
		j      RMsg
		err    error
	)

	regexpWeather := regexp.MustCompile("(w|weather|п|погода|погодка|погадка)5?[[:space:]]+.+")

	log.Debugf("Incomming raw json: %s", msg)

	if err := json.Unmarshal([]byte(msg), &j); err != nil {
		log.Warnf("Unable to to parse message from redis channel: %s", err)

		return
	}

	j, err = ValidateRmsg(cfg, j, msg)

	if err != nil {
		log.Warnf("%s", err)

		return
	}

	// Если у нас циклическая пересылка сообщения, попробуем её тут разорвать, отбросив сообщение
	if j.Misc.Fwdcnt > cfg.ForwardsMax {
		log.Warnf("Discarding msg with fwd_cnt exceeding max value: %s", msg)

		return
	}

	j.Misc.Fwdcnt++

	sendTo = j.Plugin

	// Классифицируем входящие сообщения. Первым делом, попробуем определить команды
	if j.Message[0:len(j.Misc.Csign)] == j.Misc.Csign {
		var cmd = j.Message[len(j.Misc.Csign):]

		switch {
		case cmd == "cat" || cmd == "кис":
			// Пытаемся вычитать ответ thecatapi 3 раза
			for i := 0; i < 3; i++ {
				answer, err = thecatapi.APIClient(cfg)

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
					j.Message = fmt.Sprintf("[%s](%s)", cats[rand.IntN(len(cats))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "fox" || cmd == "лис":
			// Пытаемся вычитать ответ randonfox.ca 3 раза
			for i := 0; i < 3; i++ {
				answer, err = randomfox.APIClient(cfg)

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
				answer, err = xkcdru.APIClient(cfg)

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
				answer, err = bunicomic.APIClient(cfg)

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
				answer, err = anekdotru.APIClient(cfg)

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
				answer, err = monkeyuser.APIClient(cfg)

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
				answer, err = obutts.APIClient(cfg)

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
					j.Message = fmt.Sprintf("[%s](%s)", arts[rand.IntN(len(arts))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "titts" || cmd == "boobs" || cmd == "tities" || cmd == "boobies" || cmd == "сиси" || cmd == "сисечки":
			// Пытаемся вычитать ответ api.oboobs.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = oboobs.APIClient(cfg)

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
					j.Message = fmt.Sprintf("[%s](%s)", arts[rand.IntN(len(arts))], answer)
				} else {
					j.Message = answer
				}
			}
		case cmd == "drink" || cmd == "праздник":
			// Пытаемся вычитать ответ prazdniki-segodnya.ru 3 раза
			for i := 0; i < 3; i++ {
				answer, err = prazdnikisegodnyaru.PsrClient(cfg)

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

		case cmd == "frog" || cmd == "лягушка":
			for i := 0; i < 3; i++ {
				answer, err = flickr.SearchByTags(cfg, []string{"frog", "toad", "amphibian"})

				if err != nil {
					log.Errorf("Try %d/3 unable to query flickr api about frogs: %s", i+1, err)
					time.Sleep(1 * time.Second)

					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Лягушки упрыгали, ни одной не нашлось."
			} else {
				// Костыль для красивого ответа в telegram.
				if j.Plugin == "telegram" {
					art := []string{"frog", "toad", "лягушка"}
					answer = fmt.Sprintf("[%s](%s)", art[rand.IntN(len(art))], answer)
				}

				j.Message = answer
			}

		case cmd == "owl" || cmd == "сова" || cmd == "сыч":
			for i := 0; i < 3; i++ {
				answer, err = flickr.SearchByTags(cfg, []string{"owlet", "owl", "raptor", "bird of prey", "nocturnal bird", "barn owl", "tawny owl", "brown owl"})

				if err != nil {
					log.Errorf("Try %d/3 unable to query flickr api about owls: %s", i+1, err)
					time.Sleep(1 * time.Second)

					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету сов, все разлетелись."
			} else {
				// Костыль для красивого ответа в telegram.
				if j.Plugin == "telegram" {
					answer = fmt.Sprintf("[{ O v O }](%s)", answer)
				}

				j.Message = answer
			}

		case cmd == "horse" || cmd == "лошадь" || cmd == "лошадка":
			for i := 0; i < 3; i++ {
				answer, err = flickr.SearchByTags(cfg, []string{"horse", "equine", "stallion", "mare", "steed", "thoroughbred"})

				if err != nil {
					log.Errorf("Try %d/3 unable to query flickr api about horses: %s", i+1, err)
					time.Sleep(1 * time.Second)

					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету коняшек, все разбежались."
			} else {
				// Костыль для красивого ответа в telegram.
				if j.Plugin == "telegram" {
					art := []string{"horse", "лошадь", "лошадка"}
					answer = fmt.Sprintf("[%s](%s)", art[rand.IntN(len(art))], answer)
				}

				j.Message = answer
			}

		case cmd == "rabbit" || cmd == "bunny" || cmd == "кролик":
			for i := 0; i < 3; i++ {
				answer, err = flickr.SearchByTags(cfg, []string{"bunny", "rabbit", "buck"})

				if err != nil {
					log.Errorf("Try %d/3 unable to query flickr api about rabbits: %s", i+1, err)
					time.Sleep(1 * time.Second)

					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету кроликов, все разбежались."
			} else {
				// Костыль для красивого ответа в telegram.
				if j.Plugin == "telegram" {
					answer = fmt.Sprintf("[(\\_/)](%s)", answer)
				}

				j.Message = answer
			}

		case cmd == "snail" || cmd == "улитка":
			for i := 0; i < 3; i++ {
				answer, err = flickr.SearchByTags(cfg, []string{"snail", "slug"})

				if err != nil {
					log.Errorf("Try %d/3 unable to query flickr api about snails: %s", i+1, err)
					time.Sleep(1 * time.Second)

					continue
				}

				break
			}

			if answer == "" {
				j.Message = "Нету улиток, все расползлись."
			} else {
				// Костыль для красивого ответа в telegram.
				if j.Plugin == "telegram" {
					art := []string{"'-'_@_", "@╜", "@_'-'"}
					answer = fmt.Sprintf("[%s](%s)", art[rand.IntN(len(art))], answer)
				}

				j.Message = answer
			}

		case cmd == "w" || cmd == "п" || cmd == "погода" || cmd == "weather" || cmd == "погодка" || cmd == "погадка":
			// TODO: научиться различать 5-дневные и моментальные прогнозы. Пока умеем только моментальные.
			city := openweathermap.QueryOwmUserCache(cfg, j.Chatid, j.Userid)

			if city == "" {
				j.Message = "Не припоминаю, какой город вас интересовал в прошлый раз."
			} else {
				if answer, err := openweathermap.OwmClient(cfg, city, 0); err != nil {
					log.Errorf("Unable to handle city %s in openweartermap api: %s", city, err)
					j.Message = "Я не знаю, какая погода в " + city
				} else {
					j.Message = answer
				}
			}

		case regexpWeather.MatchString(cmd):
			re := regexp.MustCompile("[[:space:]]+")
			cmdStr := re.Split(cmd, 2)[0]
			city := re.Split(cmd, 2)[1]

			if regexp.MustCompile(`5$`).MatchString(cmdStr) {
				if answer, err := openweathermap.OwmClient(cfg, city, 5); err != nil {
					log.Errorf("Unable to handle city %s in openweartermap api: %s", city, err)
					j.Message = "Я не знаю, какая погода в " + city
				} else {
					j.Message = answer
				}
			} else {
				if answer, err := openweathermap.OwmClient(cfg, city, 0); err != nil {
					log.Errorf("Unable to handle city %s in openweartermap api: %s", city, err)
					j.Message = "Я не знаю, какая погода в " + city
				} else {
					j.Message = answer
				}
			}

			if err = openweathermap.UpdateOwmUserCache(cfg, j.Chatid, j.Userid, city); err != nil {
				log.Errorf(
					"Unable to save value to cache, for chatid %s, userid %s and city %s: %s",
					j.Chatid,
					j.Userid,
					city,
					err,
				)
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
	var message SMsg

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
	if err := cfg.RedisClient.Publish(ctx, sendTo, data).Err(); err != nil {
		log.Warnf("Unable to send data to redis channel %s: %s", sendTo, err)
	} else {
		log.Debugf("Send msg to redis channel %s: %s", sendTo, string(data))
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
