package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"

	"github.com/zoulls/provencal-le-gaulois/pkg/event"
	"github.com/zoulls/provencal-le-gaulois/pkg/message"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

type Option struct {
	Cron          *cron.Cron
	LaunchTime    time.Time
	BuildInfo     BuildInfo
	TwitterClient *anaconda.TwitterApi
}

type BuildInfo struct {
	Version   string
	BuildTime string
	GitBranch string
	GitCommit string
}

func GetApplicationCommand() []*discordgo.ApplicationCommand {
	integerOptionMinValue := 1.0
	return []*discordgo.ApplicationCommand{
		{
			Name:        "placeholder",
			Description: "send placeholder message",
		},
		{
			Name:        "uptime",
			Description: "return uptime duration",
		},
		{
			Name:        "version",
			Description: "return version and build info",
		},
		{
			Name:        "list",
			Description: "list 10 messages of the chan",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Amount of message listed",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
			},
		},
		{
			Name:        "debug",
			Description: "Active or not debug logs",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "enable",
					Description: "Enable debug logs",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete",
			Description: "Delete messages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Amount of message you want to delete",
					MinValue:    &integerOptionMinValue,
					MaxValue:    100,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "before-url",
					Description: "Delete message before message URL",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "after-url",
					Description: "Delete message after message URL",
					Required:    false,
				},
			},
		},
		{
			Name:        "d4event",
			Description: "Check diablo IV event messages, check every X minutes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "time",
					Description: "Time between each check (in minutes)",
					MinValue:    &integerOptionMinValue,
					MaxValue:    60,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "event-message-url",
					Description: "Message URL to put event timer summary",
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "helltide-active",
					Description: "Active Helltide for timer init",
				},
			},
		},
		{
			Name:        "twitter",
			Description: "Report gaming tweets in the channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "time",
					Description: "Time between each check (in minutes)",
					MinValue:    &integerOptionMinValue,
					MaxValue:    60,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "list-id",
					Description: "Twiiter list ID",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "since-id",
					Description: "Tweets more recent than this ID in the list",
					Required:    true,
				},
			},
		},
	}
}

func GetCommandHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate, Option) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option){
		"placeholder": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "placeholder"
			log.Debugf("Received cmd %s", cmdName)
			placeholder(s, i)
			log.Debugf("End cmd %s", cmdName)
		},
		"uptime": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "uptime"
			log.Debugf("Received cmd %s", cmdName)
			uptime(s, i, opt)
			log.Debugf("End cmd %s", cmdName)
		},
		"version": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "version"
			log.Debugf("Received cmd %s", cmdName)
			version(s, i, opt)
			log.Debugf("End cmd %s", cmdName)
		},
		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "list"
			log.Debugf("Received cmd %s", cmdName)
			list(s, i)
			log.Debugf("End cmd %s", cmdName)
		},
		"debug": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "debug"
			log.Debugf("Received cmd %s", cmdName)
			debug(s, i)
			log.Debugf("End cmd %s", cmdName)
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "delete"
			log.Debugf("Received cmd %s", cmdName)
			delete(s, i)
			log.Debugf("End cmd %s", cmdName)
		},
		"d4event": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "d4event"
			log.Debugf("Received cmd %s", cmdName)
			d4Event(s, i, opt)
			log.Debugf("End cmd %s", cmdName)
		},
		"twitter": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "twitter"
			log.Debugf("Received cmd %s", cmdName)
			twitter(s, i, opt)
			log.Debugf("End cmd %s", cmdName)
		},
	}
}

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

func uptime(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "uptime: " + time.Since(opt.LaunchTime).String(),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func debug(s *discordgo.Session, i *discordgo.InteractionCreate) {
	active := i.ApplicationCommandData().Options[0].BoolValue()
	if active {
		log.SetLevel(log.ParseLevel("debug"))
		log.Info("enable debug log")
	} else {
		log.SetLevel(log.ParseLevel("info"))
		log.Info("disable debug log")
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "debug: " + fmt.Sprint(active),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func version(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	var embedsMsg []*discordgo.MessageEmbed

	embedsMsg = append(embedsMsg, &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Bot build info"),
		Description: "",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Version",
				Value: opt.BuildInfo.Version,
			},
			{
				Name:  "Build time UTC",
				Value: opt.BuildInfo.BuildTime,
			},
			{
				Name:  "Git branch",
				Value: opt.BuildInfo.GitBranch,
			},
			{
				Name:  "Git commit",
				Value: opt.BuildInfo.GitCommit,
			},
		},
	})

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embedsMsg,
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func list(s *discordgo.Session, i *discordgo.InteractionCreate) {
	amount := i.ApplicationCommandData().Options[0].IntValue()
	messageList, err := message.List(s, i.ChannelID, int(amount), "", "", "")
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

	var embedsMsg []*discordgo.MessageEmbed

	for key, msg := range messageList {
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

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embedsMsg,
		},
	})
	if err != nil {
		log.With("err", err).Error("send discord embed message")
	}
}

func delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options
	// Convert option slice into a map
	var (
		amount   int
		beforeID string
		afterID  string
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
		}
	}

	messageList, err := message.List(s, i.ChannelID, amount, beforeID, afterID, "")
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

	var failed bool
	for _, msg := range messageList {
		err = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			failed = true
			log.With("msgID", msg.ID, "err", err).Error("can't delete message")
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("can't delete message ID %s", msg.ID),
				},
			})
			if err != nil {
				log.With("err", err).Error("send error message")
			}
		}
	}

	msg := fmt.Sprintf("Done %d messages deleted !", len(messageList))
	if failed {
		msg += " with error"
	}
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func d4Event(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Author ID of D4 tracker
	durationStr := fmt.Sprintf("@every %dm", i.ApplicationCommandData().Options[0].IntValue())

	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options
	channelID := i.ChannelID
	// Convert option slice into a map
	var (
		duration       int
		eventMessageId string
		activeHelltide bool
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "time":
			duration = int(optParam.IntValue())
		case "event-message-url":
			eURL := strings.Split(optParam.StringValue(), "/")
			eventMessageId = eURL[len(eURL)-1]
		case "helltide-active":
			activeHelltide = optParam.BoolValue()
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Check event diablo IV every %d minutes", duration),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}

	// Init event timer array
	eventTimers := event.InitEventTimer()

	if len(eventMessageId) == 0 {
		// create event timer msg
		msgEmbed := event.TimerMsg(eventTimers)
		// create event timer msg in discord
		eventEmbedMsg, err := s.ChannelMessageSendEmbed(channelID, &msgEmbed)
		if err != nil {
			log.With("err", err).Error("send discord embed message")
		}
		eventMessageId = eventEmbedMsg.ID
	} else {
		// retrieve discord event timer msg
		eventMsg, err := s.ChannelMessage(i.ChannelID, eventMessageId)
		if err != nil {
			log.With("err", err).Error("get event message on discord")
		}
		// parse date from event timer msg
		eventTimers = event.ParseTimerMsg(eventMsg, eventTimers)
	}

	// Update data timer with d4armory.io
	eventTimers, err = event.PopulateEventTimer(eventTimers, activeHelltide)
	if err != nil {
		log.With("err", err).Error("populate event timer")
	}

	// create event timer discord msg
	msgEmbed := event.TimerMsg(eventTimers)
	// Update event timer msg on discord
	_, err = s.ChannelMessageEditEmbed(channelID, eventMessageId, &msgEmbed)
	if err != nil {
		log.With("err", err).Error("edit embed message")
	}

	_, err = opt.Cron.AddFunc(
		durationStr,
		func() {
			log.Debug("Check D4 events")

			eventTimers, err = event.RefreshEventTimers(eventTimers)
			if err != nil {
				log.With("err", err).Error("refresh event timers")
			}

			msgEmbed := event.TimerMsg(eventTimers)
			_, err = s.ChannelMessageEditEmbed(channelID, eventMessageId, &msgEmbed)
			if err != nil {
				log.With("err", err).Error("edit embed message")
			}

			log.Debug("Check D4 events done")
		})
	if err != nil {
		log.With("err", err).Error("D4 events cron creation")
	}
	opt.Cron.Start()

	log.Infof("Init cron schedule to check event diablo IV every %d minutes", duration)
}

func twitter(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		duration  int
		listId    int64
		listIdStr string
		listName  string
		sinceId   string
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "time":
			duration = int(optParam.IntValue())
		case "list-id":
			listIdStr = optParam.StringValue()
		case "since-id":
			sinceId = optParam.StringValue()
		}
	}

	// duration for cron
	durationStr := fmt.Sprintf("@every %dm", duration)

	// listId conversion
	listId, err := strconv.ParseInt(listIdStr, 10, 64)
	if err != nil {
		log.With("err", err).Error("listId conversion")
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

	// Retrieve Twitter list data
	tList, err := opt.TwitterClient.GetList(listId, url.Values{})
	if err != nil {
		log.With("err", err).Error("retrieve Twitter list data")
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
	listName = tList.Name

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Check tweets every %d minutes for the list %s (%d)", duration, listName, listId),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}

	_, err = opt.Cron.AddFunc(
		durationStr,
		func() {
			log.Debugf("Check tweetsList %s (%d)", listName, listId)

			v := url.Values{}
			v.Set("since_id", sinceId)

			tweetsList, err := opt.TwitterClient.GetListTweets(listId, false, v)
			if err != nil {
				log.With("err", err).Error("GetListTweets")
			}

			tweetsNb := len(tweetsList)
			log.Debugf("tweetsList %s (%d) count: %d", listName, listId, tweetsNb)

			if tweetsNb > 0 {
				// retrieve the most recent tweet id for next schedule
				sinceId = strconv.FormatInt(tweetsList[0].Id, 10)

				// the most recent tweet being first, the loop is done in the descending direction
				for idx := tweetsNb - 1; idx >= 0; idx-- {
					// generate twitter url
					tUrl := utils.URLFromTweet(tweetsList[idx])
					// send message
					_, err := s.ChannelMessageSend(channelID, tUrl)
					if err != nil {
						log.With("err", err).Error("send error message")
					}
				}
			}

			log.Debugf("Check tweetsList %s (%d) done", listName, listId)
		})
	if err != nil {
		log.With("err", err).Error("Twitter cron creation")
	}
	opt.Cron.Start()

	log.Infof("Init cron schedule to check tweets every %d minutes for the list %s (%d)", duration, listName, listId)
}
