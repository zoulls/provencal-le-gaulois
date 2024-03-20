package rss

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

func ParseRSS(url string, lastGUID string) ([]string, error) {
	// Init var
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	listRSSMsg := make([]string, 0)

	// Init RSS parser
	fp := gofeed.NewParser()

	// Parse RSS URL
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return listRSSMsg, err
	}

	// Retrieve GUID since lastGUID
	if len(feed.Items) > 0 {
		for _, v := range feed.Items {
			if v.GUID == lastGUID {
				break
			}
			listRSSMsg = append(listRSSMsg, v.GUID)
		}
	}

	return listRSSMsg, err
}
