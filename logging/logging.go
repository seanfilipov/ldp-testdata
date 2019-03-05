package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func Init() {
	Logger.Out = os.Stdout
	Logger.Level = logrus.InfoLevel
	Logger.Formatter = new(logrus.TextFormatter)
}
