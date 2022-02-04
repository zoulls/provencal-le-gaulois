package status

import (
	"fmt"
	"time"

	"github.com/maritimusj/durafmt"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

// See http://golang.org/pkg/time/#Parse
// Time is in UTC
const timeFormat = "2006-01-02 15:04"

var units = map[string]string{
	"years":        "années",
	"weeks":        "semaines",
	"days":         "jours",
	"hours":        "heures",
	"minutes":      "minutes",
	"seconds":      "secondes",
	"milliseconds": "millisecondes",
	"microseconds": "microsecondes",
}

func generateCountdown(force bool) (*string, error) {
	conf := config.GetConfig()
	var status *string

	// Avoid to many update
	if !force && !utils.MoreThan(conf.StatusUpdate.Every, lastSync) {
		return nil, nil
	}

	deadline, err := time.Parse(timeFormat, conf.StatusUpdate.Date)
	if err != nil {
		return nil, err
	}
	timeDuration := time.Until(deadline)

	if timeDuration.Seconds() > float64(0) {
		initUnits()
		msg := fmt.Sprintf(
			"attendre %s avant %s",
			durafmt.Parse(timeDuration).LimitFirstN(conf.StatusUpdate.NbUnits),
			conf.StatusUpdate.Game,
		)
		status = &msg
	}

	return status, nil
}

func initUnits() {
	for u, a := range units {
		durafmt.SetAlias(u, a)
	}
}
