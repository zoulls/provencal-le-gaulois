package status

import (
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
	"time"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
)

var lastSync = time.Now()
var currentStatus string

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

func (s *Status) GetDefault() (string, error) {
	//Take redis status
	statusInRedis, err := s.rClient.GetDefaultStatus()
	if err != nil {
		return "", err
	}

	status := utils.StringValue(statusInRedis)
	if len(status) > 0 {
		return status, err
	}
	//Take status in config it not defined in redis
	return s.config.Status, err
}

func (s *Status) Last(force bool) (string, error) {
	var err error

	conf := config.GetConfig()

	// Init with default status
	status, err := s.GetDefault()
	if err != nil {
		return status, err
	}

	if conf.StatusUpdate.Enabled {
		statusCtd, err := generateCountdown(force)
		if err != nil {
			return status, err
		}
		status = currentStatus
		if statusCtd != nil {
			status = utils.StringValue(statusCtd)
		}
	}

	lastSync = time.Now()
	currentStatus = status
	return status, err
}

func  (s *Status) GetCurrentStatus() string {
	return currentStatus
}

func  (s *Status) SetCurrentStatus(status string) {
	currentStatus = status
}

func GetLastSync() string {
	return lastSync.String()
}
