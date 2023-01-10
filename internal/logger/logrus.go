package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetLogger() *log.Logger {
	var l = log.New()
	l.SetFormatter(&log.JSONFormatter{})

	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		l.SetLevel(log.DebugLevel)
	} else {
		l.SetLevel(log.InfoLevel)
	}

	return l
}
