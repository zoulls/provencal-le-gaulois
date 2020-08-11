package status

import (
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
	"time"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
)

var lastSync = time.Now()

type Status struct {
	config  *config.Config
	rClient redis.Client
}

func New(config *config.Config, rClient redis.Client) *Status {
	return &Status{
		config:  config,
		rClient: rClient,
	}
}

func (s *Status) GetDefault() (*string, error) {
	status, err := s.rClient.GetDefaultStatus()
	if err != nil {
		return nil, err
	}
	if status != nil {
		return status, err
	}
	return &s.config.Status, err
}

func (s *Status) Last(force bool) (string, error) {
	var err error

	conf := config.GetConfig()

	// Init default status
	status := utils.String(conf.Status)

	if conf.StatusUpdate.Enabled {
		status, err = generateCountdown(force)
		if err != nil {
			return utils.StringValue(status), err
		}
	}


	status, err = s.rClient.GetDefaultStatus()
	if err != nil {
		return utils.StringValue(status), err
	}

	lastSync = time.Now()
	return utils.StringValue(status), err
}

func GetLastSync() string {
	return lastSync.String()
}
