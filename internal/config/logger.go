package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// Set log environment
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "debug"
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.DebugLevel
	}

	logrus.SetLevel(level)
}
