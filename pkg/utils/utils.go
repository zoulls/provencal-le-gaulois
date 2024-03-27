package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

// UnixStringToTime convert unix string to Time
func UnixStringToTime(timestamp string) (time.Time, error) {
	var tm time.Time
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return tm, err
	}
	tm = time.Unix(i, 0)
	return tm, nil
}

// URLFromTweet generate Twitter URL from tweet data
func URLFromTweet(t anaconda.Tweet) string {
	return "https://twitter.com/" + t.User.ScreenName + "/status/" + t.IdStr
}

// HumanizeDuration humanizes time.Duration output to a meaningful value
func HumanizeDuration(d time.Duration) string {
	if d < day {
		return d.Round(time.Second).String()
	}

	var b strings.Builder
	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d.Round(time.Second).String())

	return b.String()
}
