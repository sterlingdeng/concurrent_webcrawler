package internal

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger(lvl string) *log.Logger {
	level, err := log.ParseLevel(lvl)
	if err != nil {
		log.Infof("Could not parse log level: %s, Defaulting to Debug level. Err: %s", lvl, err.Error())
		level = log.DebugLevel
	}
	return &log.Logger{
		Out:          os.Stdout,
		Hooks:        log.LevelHooks{},
		Formatter:    &log.TextFormatter{},
		Level:        level,
	}
}
