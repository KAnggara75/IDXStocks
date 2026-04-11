package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// Set log environment
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
