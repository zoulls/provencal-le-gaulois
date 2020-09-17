package twitter

import (
	"strings"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

func SyncList(rClient redis.Client, tConfig config.Twitter) (config.Twitter, error) {
	var followIDstring string

	for _, follow := range tConfig.TwitterFollows {
		list, err := rClient.GetTwitterFollows(follow)
		if err != nil {
			return tConfig, err
		}
		listStr := utils.StringValue(list)
		follow.ListStr = listStr
		follow.List = strings.Split(listStr, ",")
		if len(followIDstring) > 0 {
			followIDstring = followIDstring + ","
		}
		followIDstring = followIDstring + listStr
	}
	tConfig.FollowIDstring = followIDstring

	return tConfig, nil
}
