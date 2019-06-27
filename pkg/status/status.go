package status

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02 15:04 MST"
	shDate     = "2019-06-28 09:00 UTC"
)

func Update(s *discordgo.Session) error {
	deadline, err := time.Parse(timeFormat, shDate)
	if err != nil {
		return err
	}
	diff := time.Until(deadline)

	status := "Shadowbringers !"
	if diff.Seconds() > float64(0) {
		out := time.Time{}.Add(diff)
		status = fmt.Sprintf("attendre %s pour Shadowbringers", out.Format("15h 04m 05s"))
	}

	return s.UpdateStatus(0, status)
}
