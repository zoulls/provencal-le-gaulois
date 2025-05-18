package rss

import (
	"context"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/mmcdole/gofeed"
)

// RSS_NB_LAST_GUID is the number of last GUIDs to keep in memory
var nbLastGUID int

// DefaultNbMax is the default number of items to fetch from the RSS feed
var DefaultNbMax int

func InitRSS() {
	var err error
	// Get the number of last GUIDs from environment variable
	strNbLastGUID := os.Getenv("RSS_NB_LAST_GUID")

	if strNbLastGUID == "" {
		// Default value if not set
		nbLastGUID = 5
	}
	// Convert string to int
	nbLastGUID, err = strconv.Atoi(strNbLastGUID)
	if err != nil {
		log.With("err", err).Fatal("RSS_NB_LAST_GUID bad conversion")
	}

	// Get the number of last GUIDs from environment variable
	strDefaultNbMax := os.Getenv("RSS_DEFAULT_NB_MAX")

	if strDefaultNbMax == "" {
		// Default value if not set
		DefaultNbMax = 10
	}
	// Convert string to int
	DefaultNbMax, err = strconv.Atoi(strDefaultNbMax)
	if err != nil {
		log.With("err", err).Fatal("RSS_DEFAULT_NB_MAX bad conversion")
	}
}

// containsGUID checks if a GUID is present in the list of GUIDs.
func containsGUID(guid string, guidList []string) bool {
	for idx := len(guidList) - 1; idx >= 0; idx-- {
		if guidList[idx] == guid {
			return true
		}
	}
	return false
}

// updateLastGUIDs updates the list of last GUIDs, ensuring it does not exceed maxSize.
func updateListGUIDs(listGUIDs []string, newGUID string, maxSize int) []string {
	// Check if the new GUID already exists in the list
	for _, guid := range listGUIDs {
		if guid == newGUID {
			return listGUIDs
		}
	}

	// Add new GUID to the list
	listGUIDs = append(listGUIDs, newGUID)

	// If the size exceeds the limit, remove the oldest
	if len(listGUIDs) > maxSize {
		listGUIDs = listGUIDs[1:]
	}

	return listGUIDs
}

// ParseRSS fetches and parses the RSS feed from the given URL.
func ParseRSS(url string, nbMax int, lastGUIDs []string) ([]string, []string, error) {
	// Init var
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	listRSSMsg := make([]string, 0, nbMax)
	first := false

	// detect first loop
	if len(lastGUIDs) == 0 {
		first = true
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return listRSSMsg, lastGUIDs, err
	}

	// Check if the feed is empty
	if len(feed.Items) > 0 {
		for k, v := range feed.Items {
			if containsGUID(v.GUID, lastGUIDs) {
				break
			}

			// Get only the first nbMax items
			if k < nbMax {
				listRSSMsg = append(listRSSMsg, v.GUID)
			}

			// Get only the first nbLastGUID items
			if k < nbLastGUID {
				lastGUIDs = updateListGUIDs(lastGUIDs, v.GUID, nbLastGUID)
			}
		}
	}

	if first {
		// Reverse the order of the last GUIDs
		slices.Reverse(lastGUIDs)
	}

	return listRSSMsg, lastGUIDs, err
}
