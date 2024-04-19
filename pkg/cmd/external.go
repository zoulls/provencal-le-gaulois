package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/event"
	"github.com/zoulls/provencal-le-gaulois/pkg/rss"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

func d4Event(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
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

	// Convert duration to string duration for cron
	durationStr := fmt.Sprintf("@every %dm", duration)

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

	// Job function for Cron
	job := func() {
		log.Debug("check D4 events")

		eventTimers, err = event.RefreshEventTimers(eventTimers)
		if err != nil {
			log.With("err", err).Error("refresh event timers")
		}

		msgEmbed := event.TimerMsg(eventTimers)
		_, err = s.ChannelMessageEditEmbed(channelID, eventMessageId, &msgEmbed)
		if err != nil {
			log.With("err", err).Error("edit embed message")
		}

		log.Debug("check D4 events done")
	}
	// First exec
	job()

	_, err = opt.Cron.AddFunc(durationStr, job)
	if err != nil {
		log.With("err", err).Error("D4 events cron creation")
	}
	opt.Cron.Start()

	log.Infof("init cron schedule to check event diablo IV every %d minutes", duration)
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

	// Convert duration to string duration for cron
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

	// Job function for Cron
	job := func() {
		log.Debugf("check tweetsList %s (%d)", listName, listId)

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

		log.Debugf("check tweetsList %s (%d) done", listName, listId)
	}
	// First exec
	job()

	_, err = opt.Cron.AddFunc(durationStr, job)
	if err != nil {
		log.With("err", err).Error("Twitter cron creation")
	}
	opt.Cron.Start()

	log.Infof("init cron schedule to check tweets every %d minutes for the list %s (%d)", duration, listName, listId)
}

func rssParser(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		taskName string
		duration int
		rssURL   string
		lastGUID string
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "name":
			taskName = optParam.StringValue()
		case "time":
			duration = int(optParam.IntValue())
		case "url":
			rssURL = optParam.StringValue()
		case "last-guid":
			lastGUID = optParam.StringValue()
		}
	}

	// Convert duration to string duration for cron
	durationStr := fmt.Sprintf("@every %dm", duration)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf("Check %s every %d minutes", taskName, duration),
		},
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}

	// Job function for Cron
	job := func() {
		log.Debugf("exec %s", taskName)

		listRSSMsg, err := rss.ParseRSS(rssURL, lastGUID)
		if err != nil {
			log.With("err", err).With("taskName", taskName).Error("RSS parse")
		}

		cpt := len(listRSSMsg)
		if cpt > 0 {
			for idx := cpt - 1; idx >= 0; idx-- {
				_, err := s.ChannelMessageSend(channelID, listRSSMsg[idx])
				if err != nil {
					log.With("err", err).With("taskName", taskName).Error("send error message")
				}
			}
			lastGUID = listRSSMsg[0]
		}

		log.Debugf("exec %s done", taskName)
	}
	// First exec
	job()

	_, err = opt.Cron.AddFunc(durationStr, job)
	if err != nil {
		log.With("err", err).With("taskName", taskName).Error("cron creation")
	}
	opt.Cron.Start()

	log.Infof("init cron schedule to exec %s every %d minutes", taskName, duration)
}
