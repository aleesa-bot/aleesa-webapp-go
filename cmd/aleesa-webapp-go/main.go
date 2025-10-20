package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aleesa-webapp-go/internal/log"
	"aleesa-webapp-go/internal/webapp"

	"github.com/go-redis/redis/v8"
)

// Основная функция программы, не добавить и не убавить.
func main() {
	var (
		logfile *os.File
	)

	err := webapp.ReadConfig()

	if err != nil {
		log.Errorf("Unable to parse config: %s", err)
		os.Exit(1)
	}

	// Откроем лог и скормим его логгеру.
	if webapp.Config.Log != "" {
		logfile, err = os.OpenFile(webapp.Config.Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

		if err != nil {
			log.Errorf("Unable to open log file %s: %s", webapp.Config.Log, err)
			os.Exit(1)
		}
	}

	log.Init(webapp.Config.Loglevel, logfile)

	// Иницализируем клиента Редиски.
	webapp.Config.RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", webapp.Config.Server, webapp.Config.Port),
	}).WithContext(webapp.Ctx).WithTimeout(time.Duration(webapp.Config.Timeout) * time.Second)

	// Обозначим, что хотим после соединения подписаться на события из канала config.Channel.
	webapp.Subscriber = webapp.Config.RedisClient.Subscribe(webapp.Ctx, webapp.Config.Channel)

	// Самое время поставить трапы на сигналы.
	signal.Notify(webapp.SigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Запустим обработчик сигналов.
	go webapp.SigHandler(webapp.Config)

	// Начнём выгребать события из редиски (длина конвеера/буфера канала по-умолчанию - 100 сообщений).
	ch := webapp.Subscriber.Channel()

	log.Info("Service started.")

	for msg := range ch {
		if !webapp.Shutdown {
			webapp.MsgParser(webapp.Config, webapp.Ctx, msg.Payload)
		}
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
