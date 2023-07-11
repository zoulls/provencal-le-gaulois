package message

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

func List(s *discordgo.Session, channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error) {
	messageList, err := s.ChannelMessages(channelID, limit, beforeID, afterID, aroundID)
	if err != nil {
		log.Printf("Error, can't list messages with err: %s", err.Error())
		return messageList, errors.New("error, can't list messages")
	}

	if len(messageList) == 0 {
		return messageList, errors.New("no messages")
	}

	return messageList, nil
}
