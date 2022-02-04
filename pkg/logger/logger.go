package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/zoulls/provencal-le-gaulois/config"
	"os"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	conf := config.GetConfig()
	lvl, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
}

func Log() *logrus.Logger {
	return log
}
