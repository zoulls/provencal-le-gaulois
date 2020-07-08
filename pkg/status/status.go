package status

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
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

func (s *Status) Update(d *discordgo.Session, force bool) error {
	var status *string
	var err error

	conf := config.GetConfig()

	if conf.StatusUpdate.Enabled {
		status, err = generateCountdown(force)
		if err != nil {
			return err
		}
	}

	if status == nil {
		status, err = s.rClient.GetDefaultStatus()
		if err != nil {
			return err
		}
	}

	lastSync = time.Now()

	return d.UpdateStatus(0, utils.StringValue(status))
}

func GetLastSync() string {
	return lastSync.String()
}
