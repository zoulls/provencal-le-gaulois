package reply

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/status"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
	"strings"
)

func GetReply(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageSend, error) {
	var err error

	conf := config.GetConfig()
	reply := &discordgo.MessageSend{}

	// Redis client
	rClient, err := redis.NewClient()
	if err != nil {
		logger.Log.Errorf("Error during Redis init, %v", err)
	}

	if strings.HasPrefix(m.Content, conf.PrefixCmd+"embed ") {
		reply, err = createReplyFromJson(m.Content[7:len(m.Content)])
		if err == nil {
			err = s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		}
	} else {
		switch m.Content {
		case conf.PrefixCmd + "ping":
			botStatus, err := rClient.Ping()
			if err != nil {
				logger.Log.Errorf("Error get bot status, %v", err)
			}
			reply.Content = utils.StringValue(botStatus)
		case conf.PrefixCmd + "pong":
			reply.Content = "Ping! :ping_pong: petit malin! :laughing:"
		case conf.PrefixCmd + "help":
			reply.Embed = help()
		case conf.PrefixCmd + "embedGen":
			reply.Embed = embedGenerator()
		case conf.PrefixCmd + "updateStatus":
			sClient := status.New(conf, rClient)
			lastStatus, err := sClient.Last(true)
			if err != nil {
				logger.Log.Errorf("Error retrieving the last status, %v", err)
				reply.Content = "Error retrieving the last status, use of default config"
			}

			err = s.UpdateStatus(0, lastStatus)
			if err != nil {
				logger.Log.Errorf("Error during status update, %v", err)
			} else {
				reply.Content = "Status updated successfully !"
			}

		case conf.PrefixCmd + "statusLastSync":
			reply.Content = status.GetLastSync()
		case conf.PrefixCmd + "twitterFollows":
			reply.Content = conf.Twitter.FollowIDstring
		default:
			return nil, nil
		}
	}

	return reply, err
}

func createReplyFromJson(str string) (*discordgo.MessageSend, error) {
	message := &discordgo.MessageSend{}
	err := json.Unmarshal([]byte(str), message)

	return message, err
}
