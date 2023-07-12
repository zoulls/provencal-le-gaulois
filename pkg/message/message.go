package message

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func List(s *discordgo.Session, channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	messageList, err := s.ChannelMessages(channelID, limit, beforeID, afterID, aroundID)
	if err != nil {
		log.With("err", err).Error("can't list message")
		return messageList, errors.New("can't list message")
	}

	if len(messageList) == 0 {
		return messageList, errors.New("no message")
	}

	return messageList, nil
}
