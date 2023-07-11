package utils

import (
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
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
