package rss

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

// containsGUID checks if a GUID is present in the list of GUIDs.
func containsGUID(guid string, guidList []string) bool {
	for _, g := range guidList {
		if g == guid {
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

			// Get only the first 5 items
			if k < 5 {
				lastGUIDs = updateListGUIDs(lastGUIDs, v.GUID, 5)
			}
		}
	}

	return listRSSMsg, lastGUIDs, err
}
