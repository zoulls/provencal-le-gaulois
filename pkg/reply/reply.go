package reply

import (
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/status"
)

func GetReply(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageSend, error) {
	var err error

	config := config.GetConfig()
	reply := &discordgo.MessageSend{}

	if strings.HasPrefix(m.Content, config.PrefixCmd+"embed ") {
		reply, err = createReplyFromJson(m.Content[7:len(m.Content)])
		if err == nil {
			err = s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		}
	} else {
		switch m.Content {
		case config.PrefixCmd + "ping":
			reply.Content = "Pong! :ping_pong:"
		case config.PrefixCmd + "pong":
			reply.Content = "Ping! :ping_pong: petit malin! :laughing:"
		case config.PrefixCmd + "help":
			reply.Embed = help()
		case config.PrefixCmd + "embedGen":
			reply.Embed = embedGenerator()
		case config.PrefixCmd + "updateStatus":
			err = status.Update(s, true)
			if err != nil {
				logger.Log.Printf("Error attempting to set my status, %v\n", err)
				reply.Content = "Error during the update status"
			} else {
				reply.Content = "Status updated successfully !"
			}
		case config.PrefixCmd + "statusLastSync":
			reply.Content = status.GetLastSync()
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
