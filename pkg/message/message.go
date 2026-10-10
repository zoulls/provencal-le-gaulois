package message

import (
	"errors"

	"charm.land/log/v2"
	"github.com/bwmarrin/discordgo"
)

func List(s *discordgo.Session, channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	messageList, err := s.ChannelMessages(channelID, limit, beforeID, afterID, aroundID)
	if err != nil {
		log.With("err", err).Error("can't list message")
		return messageList, errors.New("can't list message")
	}
	return messageList, nil
}
