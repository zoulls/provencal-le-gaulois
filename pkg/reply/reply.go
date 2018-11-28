package reply

import (
	"encoding/json"
	"strings"

	"bitbucket.org/zoulls/provencal-le-gaulois/config"
	"github.com/bwmarrin/discordgo"
)

func GetReply(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.MessageSend, error) {
	var err error

	config := config.GetConfig()
	reply := &discordgo.MessageSend{}

	switch {
	case m.Content == config.PrefixCmd+"ping":
		reply.Content = "Pong! :ping_pong:"
	case m.Content == config.PrefixCmd+"pong":
		reply.Content = "Ping! :ping_pong: petit malin! :laughing:"
	case m.Content == config.PrefixCmd+"help":
		reply.Embed = help()
	case m.Content == config.PrefixCmd+"embedGen":
		reply.Embed = embedGenerator()
	case strings.HasPrefix(m.Content, config.PrefixCmd+"embed "):
		reply, err = createReplyFromJson(m.Content[7:len(m.Content)])
		if err == nil {
			err = s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		}
	default:
		return nil, nil
	}

	return reply, err
}

func createReplyFromJson(str string) (*discordgo.MessageSend, error) {
	message := &discordgo.MessageSend{}
	err := json.Unmarshal([]byte(str), message)

	return message, err
}
