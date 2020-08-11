package twitter

import (
	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
	"strings"
)

func SyncList(rClient redis.Client, tConfig config.Twitter) (config.Twitter, error) {
	var FollowIDstring string

	for _, follow := range tConfig.TwitterFollows {
		list, err := rClient.GetTwitterFollows(follow)
		if err != nil {
			return tConfig, err
		}
		listStr := utils.StringValue(list)
		follow.ListStr = listStr
		follow.List = strings.Split(listStr, ",")
		FollowIDstring =  FollowIDstring + listStr
	}
	tConfig.FollowIDstring = FollowIDstring

	return tConfig, nil
}