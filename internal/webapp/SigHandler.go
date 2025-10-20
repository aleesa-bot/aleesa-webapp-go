package webapp

import (
	"aleesa-webapp-go/internal/config"
	"aleesa-webapp-go/internal/log"
	"os"
	"syscall"
)

// SigHandler хэндлер сигналов закрывает все бд и сваливает из приложения.
func SigHandler(cfg *config.MyConfig) {
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

			for dbName, db := range cfg.PcacheDB {
				if err := db.Flush(); err != nil {
					log.Errorf("Unable to flush %s db to disk: %s", dbName, err)
				}

				if err := db.Close(); err != nil {
					log.Errorf("Unable to close %s db to disk: %s", dbName, err)
				}
			}
		}

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
