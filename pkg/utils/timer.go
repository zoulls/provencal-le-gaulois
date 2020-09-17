package utils

import "time"

// Return if the last sync is older than the min in minutes
func MoreThan(min float64, lastSync time.Time) bool {
	diff := time.Since(lastSync)
	return diff.Minutes() > min
}
