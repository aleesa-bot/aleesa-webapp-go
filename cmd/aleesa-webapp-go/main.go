package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"aleesa-webapp-go/internal/webapp"

	"github.com/go-redis/redis/v8"
)

var (
	// ctx контекст редиски.
	Ctx = context.Background()

	// Subscriber объектик сабскрайбера редиски.
	Subscriber *redis.PubSub

	// Shutdown ставится в true, если мы получили сигнал на выключение.
	Shutdown = false

	// SigChan канал, в который приходят уведомления для хэндлера сигналов от траппера сигналов.
	SigChan = make(chan os.Signal, 1)
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
	}).WithContext(Ctx).WithTimeout(time.Duration(webapp.Config.Timeout) * time.Second)

	// Обозначим, что хотим после соединения подписаться на события из канала config.Channel.
	Subscriber = webapp.Config.RedisClient.Subscribe(Ctx, webapp.Config.Channel)

	// Самое время поставить трапы на сигналы.
	signal.Notify(SigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// Запустим обработчик сигналов.
	go SigHandler(webapp.Config)

	// Начнём выгребать события из редиски (длина конвеера/буфера канала по-умолчанию - 100 сообщений).
	ch := Subscriber.Channel()

	log.Info("Service started.")

	for msg := range ch {
		if !Shutdown {
			webapp.MsgParser(webapp.Config, Ctx, msg.Payload)
		}
	}
}

// SigHandler хэндлер сигналов закрывает все бд и сваливает из приложения.
func SigHandler(cfg *config.MyConfig) {
	// TODO: утащить отсюда в модули.

	var err error

	for {
		var s = <-SigChan
		switch s {
		case syscall.SIGINT:
			log.Info("Got SIGINT, quitting")
		case syscall.SIGTERM:
			log.Info("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			log.Info("Got SIGQUIT, quitting")

		// Заходим на новую итерацию, если у нас "неинтересный" сигнал
		default:
			continue
		}

		Shutdown = true

		// Отпишемся от всех каналов и закроем коннект к редиске
		if err = Subscriber.Unsubscribe(Ctx); err != nil {
			log.Errorf("Unable to unsubscribe from redis channels cleanly: %s", err)
		}

		if err = Subscriber.Close(); err != nil {
			log.Errorf("Unable to close redis connection cleanly: %s", err)
		}

		if len(cfg.PcacheDB) > 0 {
			log.Debug("Closing persistent cache db")

			for _, db := range cfg.PcacheDB {
				_ = db.Close()
			}
		}

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
