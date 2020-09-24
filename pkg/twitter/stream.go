package twitter

import (
	"errors"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/discord"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/status"
)

func getAPI() *anaconda.TwitterApi {
	conf := config.GetConfig()
	anaconda.SetConsumerKey(conf.Twitter.Config.ConsumerKey)
	anaconda.SetConsumerSecret(conf.Twitter.Config.ConsumerSecret)
	return anaconda.NewTwitterApi(conf.Twitter.Config.AccessToken, conf.Twitter.Config.AccessTokenSecret)
}

func StreamTweets(ds *discordgo.Session, sClient *status.Status, rClient redis.Client) {
	conf := config.GetConfig()
	api := getAPI()
	v := url.Values{}
	v.Set("follow", conf.Twitter.FollowIDstring)
	s := api.PublicStreamFilter(v)
	defer api.Close()

	lastPing := time.Now()

	for t := range s.C {
		// Check for status update
		if conf.StatusUpdate.Enabled {
			lastStatus, err := sClient.Last(false)
			if err != nil {
				logger.Log.Errorf("Error retrieving the last status, %v", err)
			}
			err = ds.UpdateStatus(0, lastStatus)
			if err != nil {
				logger.Log.Errorf("Error during status update, %v", err)
			}
		}

		if utils.MoreThan(conf.Redis.PingTimer, lastPing) {
			lastPing = time.Now()
			ping, err := rClient.Ping()
			if err != nil || ping == nil {
				logger.Log.Errorf("Error during redis ping, %v", err)
				discord.SendPrivateMessage(ds, conf.Discord.AdminID, "Error during redis ping")
			}
		}

		switch tweet := t.(type) {
		case anaconda.Tweet:
			if originalTweet(tweet) {
				err := createMessage(ds, &tweet)
				if err != nil {
					logger.Log.Errorf("Error during send message of tweet, %v", tweet)
				}
			}
		default:
			logger.Log.Debugf("unknown type(%T), %v", tweet, tweet)
		}
	}
}

func createMessage(ds *discordgo.Session, tweet *anaconda.Tweet) error {
	conf := config.GetConfig()

	message := discord.URLFromTweet(tweet)
	reply := &discordgo.MessageSend{
		Content: message,
	}

	discordID, err := getDiscordChanID(conf.Twitter.TwitterFollows, tweet)
	if err != nil {
		return err
	}

	_, err = ds.ChannelMessageSendComplex(discordID, reply)
	return err
}

func originalTweet(tweet anaconda.Tweet) bool {
	return tweet.RetweetedStatus == nil && tweet.InReplyToStatusID == 0 && tweet.InReplyToUserID == 0
}

func getDiscordChanID(tf []*config.TwitterFollow, tweet *anaconda.Tweet) (string, error) {
	for _, follow := range tf {
		if contains(follow.List, tweet.User.IdStr) {
			return follow.DiscordChan, nil
		}
	}
	return "", errors.New("twitter ID unknown")
}

// Contains tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
