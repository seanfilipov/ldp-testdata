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

func ParseLogLevel(logLevel string) logrus.Level {
	switch logLevel {
	case "Trace":
		return logrus.TraceLevel
	case "Debug":
		return logrus.DebugLevel
	case "Info":
		return logrus.InfoLevel
	case "Warning":
		return logrus.WarnLevel
	case "Error":
		return logrus.ErrorLevel
	case "Fatal":
		return logrus.FatalLevel
	case "Panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
