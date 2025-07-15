package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.JSONFormatter{})

	Log.SetLevel(logrus.InfoLevel)

	Log.SetOutput(os.Stdout)
}
