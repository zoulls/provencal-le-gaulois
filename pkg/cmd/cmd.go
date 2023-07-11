package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
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
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author-id",
					Description: "Author ID message to check",
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
			log.Printf("Received cmd %s", cmdName)
			placeholder(s, i)
			log.Printf("End cmd %s", cmdName)
		},
		"uptime": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "uptime"
			log.Printf("Received cmd %s", cmdName)
			uptime(s, i, opt)
			log.Printf("End cmd %s", cmdName)
		},
		"version": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "version"
			log.Printf("Received cmd %s", cmdName)
			version(s, i, opt)
			log.Printf("End cmd %s", cmdName)
		},
		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "list"
			log.Printf("Received cmd %s", cmdName)
			list(s, i)
			log.Printf("End cmd %s", cmdName)
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "delete"
			log.Printf("Received cmd %s", cmdName)
			delete(s, i)
			log.Printf("End cmd %s", cmdName)
		},
		"d4event": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "d4event"
			log.Printf("Received cmd %s", cmdName)
			d4Event(s, i, opt)
			log.Printf("End cmd %s", cmdName)
		},
		"twitter": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "twitter"
			log.Printf("Received cmd %s", cmdName)
			twitter(s, i, opt)
			log.Printf("End cmd %s", cmdName)
		},
	}
}

func placeholder(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := s.ChannelMessageSend(i.ChannelID, "placeholder")
	if err != nil {
		log.Print(err.Error())
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		if err != nil {
			log.Print(err.Error())
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
		log.Print(err.Error())
	}
}

func uptime(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "uptime: " + time.Since(opt.LaunchTime).String(),
		},
	})
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embedsMsg,
		},
	})
}

func list(s *discordgo.Session, i *discordgo.InteractionCreate) {
	amount := i.ApplicationCommandData().Options[0].IntValue()
	messageList, err := message.List(s, i.ChannelID, int(amount), "", "", "")
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: embedsMsg,
		},
	})
}

func delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options
	// Convert option slice into a map
	var (
		amount   int
		beforeID string
		afterID  string
	)
	for _, opt := range options {
		switch opt.Name {
		case "amount":
			amount = int(opt.IntValue())
		case "before-url":
			bURL := strings.Split(opt.StringValue(), "/")
			beforeID = bURL[len(bURL)-1]
		case "after-url":
			aURL := strings.Split(opt.StringValue(), "/")
			afterID = aURL[len(aURL)-1]
		}
	}

	messageList, err := message.List(s, i.ChannelID, amount, beforeID, afterID, "")
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
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
		errMsg := "Error, can't send messages"
		log.Print(errMsg)

		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		if err != nil {
			log.Print(err.Error())
		}

		return
	}

	var failed bool
	for _, msg := range messageList {
		err = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			failed = true
			errMsg := fmt.Sprintf("Error, can't delete messages with ID, %s", msg.ID)
			log.Print(errMsg)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: errMsg,
				},
			})
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
		log.Print(err.Error())
	}
}

func d4Event(s *discordgo.Session, i *discordgo.InteractionCreate, optEvent Option) {
	// Author ID of D4 tracker
	authorID := "1116956812432904323"
	durationStr := fmt.Sprintf("@every %dm", i.ApplicationCommandData().Options[0].IntValue())

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options
	channelID := i.ChannelID
	// Convert option slice into a map
	var (
		duration       int
		eventMessageId string
	)

	for _, opt := range options {
		switch opt.Name {
		case "time":
			duration = int(opt.IntValue())
		case "event-message-url":
			eURL := strings.Split(opt.StringValue(), "/")
			eventMessageId = eURL[len(eURL)-1]
		case "author-id":
			if len(opt.StringValue()) > 0 {
				authorID = opt.StringValue()
			}
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
		log.Print(err.Error())
	}

	// Init event timer array
	eventTimers := event.EvenTimerInit()

	if len(eventMessageId) == 0 {
		msgEmbed := event.TimerMsg(eventTimers, make([]bool, 3))
		eventEmbedMsg, err := s.ChannelMessageSendEmbed(channelID, &msgEmbed)
		if err != nil {
			log.Printf("Error during send message, err: %s", err)
		}
		eventMessageId = eventEmbedMsg.ID
	} else {
		eventMsg, err := s.ChannelMessage(i.ChannelID, eventMessageId)
		if err != nil {
			log.Printf("Error during Get event message on discord, err: %s", err)
		}
		event.ParseTimerMsg(eventMsg, eventTimers)
	}

	_, err = optEvent.Cron.AddFunc(
		durationStr,
		func() {
			log.Print("Check D4 events")
			newEvent := make([]bool, 3)

			messageList, err := message.List(s, i.ChannelID, 10, "", "", "")
			if err != nil {
				log.Printf("Error during list message, err: %s", err)
			}

			for _, msg := range messageList {
				var deleteMsg bool
				if msg.Author.ID == authorID {
					logMsg := "type: "
					switch {
					case strings.Contains(msg.Content, event.EventsName[event.EventWB]):
						logMsg += event.EventsName[event.EventWB]

						found := event.RegDiscordTime.FindStringSubmatch(msg.Content)
						ts := found[1]
						err := eventTimers[event.EventWB].SetNextTimestamp(ts)
						if err != nil {
							log.Printf("Error during SetNextTimestamp, err: %s", err)
						}
						newEvent[event.EventWB] = true

						if eventTimers[event.EventWB].Next.Before(time.Now()) {
							deleteMsg = true
						}
					case strings.Contains(msg.Content, event.EventsName[event.EventHelltide]):
						logMsg += event.EventsName[event.EventHelltide]

						found := event.RegDiscordTime.FindStringSubmatch(msg.Content)
						ts := found[1]
						err := eventTimers[event.EventHelltide].SetNextTimestamp(ts)
						if err != nil {
							log.Printf("Error during SetNextTimestamp, err: %s", err)
						}

						newEvent[event.EventHelltide] = true

						if eventTimers[event.EventHelltide].Next.Before(time.Now()) {
							deleteMsg = true
						}
					case strings.Contains(msg.Content, event.EventsName[event.EventLegions]):
						logMsg += event.EventsName[event.EventLegions]

						found := event.RegDiscordTime.FindStringSubmatch(msg.Content)
						ts := found[1]
						err := eventTimers[event.EventLegions].SetNextTimestamp(ts)
						if err != nil {
							log.Printf("Error during SetNextTimestamp, err: %s", err)
						}
						newEvent[event.EventLegions] = true

						if eventTimers[event.EventLegions].Next.Before(time.Now()) {
							deleteMsg = true
						}
					case strings.EqualFold(msg.Content, "[Original Message Deleted]"):
						deleteMsg = true
					default:
						logMsg += "Other"
					}

					if deleteMsg {
						logMsg += "Delete"
						err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
						if err != nil {
							log.Printf("Error during ChannelMessageDelete, err: %s", err)
						}
					}

					log.Printf("%s, authorName: %s, authorID: %s, message: %s", logMsg, msg.Author.Username, msg.Author.ID, msg.Content)
				}
			}

			event.RefreshEventTimers(eventTimers, newEvent)

			msgEmbed := event.TimerMsg(eventTimers, newEvent)
			_, err = s.ChannelMessageEditEmbed(channelID, eventMessageId, &msgEmbed)
			if err != nil {
				log.Printf("Error during send message, err: %s", err)
			}

			log.Print("Check D4 events done")
		})
	if err != nil {
		log.Printf("Error, during cron creation with err: %s", err.Error())
	}
	optEvent.Cron.Start()

	log.Printf("Init cron schedule to check event diablo IV every %d minutes", duration)
}

func twitter(s *discordgo.Session, i *discordgo.InteractionCreate, optEvent Option) {
	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

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

	for _, opt := range options {
		switch opt.Name {
		case "time":
			duration = int(opt.IntValue())
		case "list-id":
			listIdStr = opt.StringValue()
		case "since-id":
			sinceId = opt.StringValue()
		}
	}

	// duration for cron
	durationStr := fmt.Sprintf("@every %dm", duration)

	// listId conversion
	listId, err := strconv.ParseInt(listIdStr, 10, 64)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		log.Printf("Error during listId conversion, err: %s", err.Error())
		return
	}

	// Retrieve Twitter list data
	tList, err := optEvent.TwitterClient.GetList(listId, url.Values{})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: err.Error(),
			},
		})
		log.Printf("Error during retrieve Twitter list data, err: %s", err.Error())
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
		log.Print(err.Error())
	}

	_, err = optEvent.Cron.AddFunc(
		durationStr,
		func() {
			log.Printf("Check tweetsList %s (%d)", listName, listId)

			v := url.Values{}
			v.Set("since_id", sinceId)

			tweetsList, err := optEvent.TwitterClient.GetListTweets(listId, false, v)
			if err != nil {
				log.Print(err.Error())
			}

			tweetsNb := len(tweetsList)
			log.Printf("tweetsList %s (%d) count: %d", listName, listId, tweetsNb)

			if tweetsNb > 0 {
				// retrieve tweet id for next schedule
				sinceId = strconv.FormatInt(tweetsList[0].Id, 10)

				for _, tweet := range tweetsList {
					// generate twitter url
					tUrl := utils.URLFromTweet(tweet)
					// send message
					_, err := s.ChannelMessageSend(channelID, tUrl)
					if err != nil {
						log.Print(err.Error())
					}
				}
			}

			log.Printf("Check tweetsList %s (%d) done", listName, listId)
		})
	if err != nil {
		log.Printf("Error, during cron creation with err: %s", err.Error())
	}
	optEvent.Cron.Start()

	log.Printf("Init cron schedule to check tweets every %d minutes for the list %s (%d)", duration, listName, listId)
}
