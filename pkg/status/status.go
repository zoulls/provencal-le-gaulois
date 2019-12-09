package status

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"

	"github.com/zoulls/provencal-le-gaulois/config"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02 15:04 MST"
)

var lastSync = time.Now()

func Update(s *discordgo.Session, force bool) error {
	conf := config.GetConfig()
	status := conf.Status

	// Avoid to many update
	if !force && !moreThan(conf.StatusUpdate.Every) {
		return nil
	}

	deadline, err := time.Parse(timeFormat, conf.StatusUpdate.Date)
	if err != nil {
		return err
	}
	timeDuration := time.Until(deadline)

	if timeDuration.Seconds() > float64(0) {
		status = fmt.Sprintf("attendre %s avant MHW Iceborne !", durafmt.ParseShort(timeDuration).String())
	}

	lastSync = time.Now()

	return s.UpdateStatus(0, status)
}

func GetLastSync() string {
	return lastSync.String()
}

// Return if the last sync is older than the min in minutes
func moreThan(min float64) bool {
	diff := time.Since(lastSync)
	return diff.Minutes() > min
}
