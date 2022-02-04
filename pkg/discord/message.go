package discord

import (
	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
)

func SendPrivateMessage(ds *discordgo.Session, userID string, message string) {
	channel, err := ds.UserChannelCreate(userID)
	if err != nil || channel == nil {
		logger.Log().Errorf("Error during creating private channel, %v", err)
	} else {
		_, err = ds.ChannelMessageSend(channel.ID, message)
		if err != nil {
			logger.Log().Errorf("Error sending MP, %v", err)
		}
	}
}
