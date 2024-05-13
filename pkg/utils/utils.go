package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

// StringPtr return the pointer ref of a string
func StringPtr(s string) *string {
	return &s
}

// ErrorMsg build error message for discord feedback
func ErrorMsg(err error) *string {
	msg := fmt.Sprintf("Error: %s", err.Error())
	return &msg
}
