package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// Производит некоторую инициализацию перед запуском main().
func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableQuote:           true,
		DisableLevelTruncation: false,
		DisableColors:          true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
	})

	config.Loglevel = "debug"
	config.DataDir = "../../data"

	var (
		err        error
		useragents = fmt.Sprintf("%s/useragents.txt", config.DataDir)
	)

	userAgents, err = readLines(useragents)

	if err != nil {
		log.Errorf("Unable to load list of useragents from %s: %s", useragents, err)
		os.Exit(1)
	}

	// No panic, no trace.
	switch config.Loglevel {
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	log.Errorln("Monkeyuser test")

	s, err := monkeyUserClient()

	if err != nil {
		log.Errorf("function monkeyUserClient() returns error: %s", err)
		os.Exit(1)
	}

	log.Errorf("Test result is: %s", s)
}
