package status

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02 15:04 MST"
	shDate     = "2019-06-28 09:00 UTC"
)

var lastSync = time.Now()

func Update(s *discordgo.Session, force bool) error {
	conf := config.GetConfig()

	// Avoid to many update
	if !force && !moreThan(conf.StatusUpdate.Every) {
		return nil
	}

	deadline, err := time.Parse(timeFormat, shDate)
	if err != nil {
		return err
	}
	diff := time.Until(deadline)

	status := "Shadowbringers !"
	if diff.Seconds() > float64(0) {
		out := time.Time{}.Add(diff)
		status = fmt.Sprintf("attendre %s pour Shadowbringers", out.Format("15h 04m"))
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
