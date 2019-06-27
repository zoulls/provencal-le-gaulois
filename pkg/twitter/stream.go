package twitter

import (
	"errors"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/reply"
	"github.com/zoulls/provencal-le-gaulois/pkg/status"
)

func getAPI() *anaconda.TwitterApi {
	conf := config.GetConfig()
	anaconda.SetConsumerKey(conf.Twitter.Config.ConsumerKey)
	anaconda.SetConsumerSecret(conf.Twitter.Config.ConsumerSecret)
	return anaconda.NewTwitterApi(conf.Twitter.Config.AccessToken, conf.Twitter.Config.AccessTokenSecret)
}

func StreamTweets(discord *discordgo.Session) {
	conf := config.GetConfig()
	api := getAPI()
	v := url.Values{}
	v.Set("follow", conf.Twitter.FollowIDstring)
	s := api.PublicStreamFilter(v)
	// DEBUG
	// s := api.PublicStreamSample(nil)
	defer api.Close()

	for t := range s.C {
		switch tweet := t.(type) {
		case anaconda.Tweet:
			// DEBUG
			// fmt.Printf("%-15s: %s\n", tweet.User.ScreenName, tweet.Text)
			if originalTweet(tweet) {
				err := createMessage(discord, &tweet)
				if err != nil {
					logger.Log.Printf("Error during send message of tweet : %v \n", tweet)
				}
				if conf.StatusUpdate {
					err = status.Update(discord)
					if err != nil {
						logger.Log.Printf("Error attempting to set my status, %v\n", err)
					}
				}
			}
		default:
			logger.Log.Printf("unknown type(%T) : %v \n", tweet, tweet)
		}
	}
}

func createMessage(discord *discordgo.Session, tweet *anaconda.Tweet) error {
	conf := config.GetConfig()
	message := reply.FromTweet(tweet)
	reply := &discordgo.MessageSend{
		Embed: message,
	}

	discordID, err := getDiscordChanID(conf.Twitter.TwitterFollows, tweet)
	if err != nil {
		return err
	}

	_, err = discord.ChannelMessageSendComplex(discordID, reply)
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
	return "", errors.New("Twitter ID unknow")
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
