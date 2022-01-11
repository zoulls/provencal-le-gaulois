package command

import (
	"time"

	"github.com/zoulls/provencal-le-gaulois/config"
)

func GetUptime() time.Duration {
	conf := config.GetConfig()
	return time.Since(conf.GetStartDate())
}
