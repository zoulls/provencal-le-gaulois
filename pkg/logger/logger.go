package logger

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/zoulls/provencal-le-gaulois/config"
)

var Log = logrus.New()

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)

	conf := config.GetConfig()
	lvl, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		panic(err)
	}

	logrus.SetLevel(lvl)
}
