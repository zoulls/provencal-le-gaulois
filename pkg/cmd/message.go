package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/message"
)

func placeholder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := s.ChannelMessageSend(i.ChannelID, "placeholder")
	if err != nil {
		log.With("err", err).Error("send placeholder message")
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Done",
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func list(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// Init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		amount   int
		beforeID string
		afterID  string
		authorID string
	)
	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "amount":
			amount = int(optParam.IntValue())
		case "before-url":
			bURL := strings.Split(optParam.StringValue(), "/")
			beforeID = bURL[len(bURL)-1]
		case "after-url":
			aURL := strings.Split(optParam.StringValue(), "/")
			afterID = aURL[len(aURL)-1]
		case "author-id":
			authorID = optParam.StringValue()
		}
	}

	messageList, err := message.List(s, channelID, amount, beforeID, afterID, "")
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
		return
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}

	if len(messageList) == 0 {
		resp.Data.Content = "No message to list"
	} else {
		embedsMsg := make([]*discordgo.MessageEmbed, 0)

		for key, msg := range messageList {
			addFlag := true
			if len(authorID) > 0 {
				if msg.Author.ID != authorID {
					addFlag = false
				}
			}

			if addFlag {
				msgEmb := discordgo.MessageEmbed{
					Title:       fmt.Sprintf("Message %d", key+1),
					Description: msg.Content,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Author",
							Value: msg.Author.Username,
						},
						{
							Name:  "Author-ID",
							Value: msg.Author.ID,
						},
					},
				}
				embedsMsg = append(embedsMsg, &msgEmb)
			}
		}

		if len(embedsMsg) == 0 {
			resp.Data.Content = "No message to list after filter"
		} else {
			resp.Data.Embeds = embedsMsg
		}
	}

	err = s.InteractionRespond(i.Interaction, resp)
	if err != nil {
		log.With("err", err).Error("send discord embed message")
	}
}

func delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// Init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		amount   int
		beforeID string
		afterID  string
		authorID string
	)
	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "amount":
			amount = int(optParam.IntValue())
		case "before-url":
			bURL := strings.Split(optParam.StringValue(), "/")
			beforeID = bURL[len(bURL)-1]
		case "after-url":
			aURL := strings.Split(optParam.StringValue(), "/")
			afterID = aURL[len(aURL)-1]
		case "author-id":
			authorID = optParam.StringValue()
		}
	}

	messageList, err := message.List(s, channelID, amount, beforeID, afterID, "")
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
		return
	}

	if len(messageList) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "No message to delete",
			},
		})
		return
	}

	// loading message
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(5),
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.With("err", err).Error("send discord loading message")
		return
	}

	msgToDelete := make([]string, 0)
	for _, msg := range messageList {
		addFlag := true
		if len(authorID) > 0 {
			if msg.Author.ID != authorID {
				addFlag = false
			}
		}

		if addFlag {
			msgToDelete = append(msgToDelete, msg.ID)
		}
	}

	doneMsg := fmt.Sprintf("Done %d messages deleted", len(msgToDelete))
	err = s.ChannelMessagesBulkDelete(channelID, msgToDelete)
	if err != nil {
		log.With("err", err).
			With("channelID", channelID).
			Error("delete bulk message")

		doneMsg += fmt.Sprintf(" with error: %s", err.Error())
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &doneMsg,
	})
	if err != nil {
		log.With("err", err).Error("send done message")
	}
}

func autoClean(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// Init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		taskName   string
		duration   string
		expiration string
		authorID   string
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "name":
			taskName = optParam.StringValue()
		case "time":
			duration = optParam.StringValue()
		case "expiration":
			expiration = optParam.StringValue()
		case "author-id":
			authorID = optParam.StringValue()
		}
	}

	// Convert duration to string duration for cron
	durationStr := fmt.Sprintf("@every %s", duration)

	// Convert expiration to duration time object
	exp, err := time.ParseDuration(expiration)
	if err != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Exec %s every %s", taskName, duration),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}

	// Job function for Cron
	job := func() {
		log.Debugf("exec %s", taskName)

		messageList, err := message.List(s, channelID, 100, "", "", "")
		if err != nil {
			log.With("err", err).
				With("taskName", taskName).
				With("channelID", channelID).
				Error("list message")
		}

		if len(messageList) == 0 {
			log.Debug("no message listed")
			return
		}
		log.Debugf("%d messages listed before filter", len(messageList))

		msgList := make([]string, 0)

		timeExp := time.Now().Add(-exp)

		for _, msg := range messageList {
			if msg.Timestamp.Before(timeExp) {
				addFlag := true
				if len(authorID) > 0 {
					if msg.Author.ID != authorID {
						addFlag = false
					}
				}

				if addFlag {
					msgList = append(msgList, msg.ID)
				}
			}
		}

		err = s.ChannelMessagesBulkDelete(channelID, msgList)
		if err != nil {
			log.With("err", err).
				With("taskName", taskName).
				With("channelID", channelID).
				Error("delete bulk message")
		}

		log.Debugf("%d messages deleted", len(msgList))

		log.Debugf("exec %s done", taskName)
	}
	// First exec
	job()

	_, err = opt.Cron.AddFunc(durationStr, job)
	if err != nil {
		log.With("err", err).With("taskName", taskName).Error("cron creation")
	}
	opt.Cron.Start()

	log.Infof("init cron schedule to exec %s every %s", taskName, duration)
}
