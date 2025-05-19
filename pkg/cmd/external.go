package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/event"
	"github.com/zoulls/provencal-le-gaulois/pkg/rss"
	"github.com/zoulls/provencal-le-gaulois/pkg/task"
	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

func d4Event(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options
	channelID := i.ChannelID
	// Convert option slice into a map
	var (
		duration       string
		eventMessageId string
		activeHelltide bool
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "time":
			duration = optParam.StringValue()
		case "event-message-url":
			eURL := strings.Split(optParam.StringValue(), "/")
			eventMessageId = eURL[len(eURL)-1]
		case "helltide-active":
			activeHelltide = optParam.BoolValue()
		}
	}

	// Convert duration to string duration for cron
	durationStr := fmt.Sprintf("@every %s", duration)

	// Check duration format
	_, err := time.ParseDuration(duration)
	if err != nil {
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.ErrorMsg(err),
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
		return
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.StringPtr(fmt.Sprintf("Check event diablo IV every %s", duration)),
	})
	if err != nil {
		log.With("err", err).Error("send message")
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
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: utils.ErrorMsg(err),
			})
			if err != nil {
				log.With("err", err).Error("send error message")
			}
		}
		eventMessageId = eventEmbedMsg.ID
	} else {
		// retrieve discord event timer msg
		eventMsg, err := s.ChannelMessage(i.ChannelID, eventMessageId)
		if err != nil {
			log.With("err", err).Error("get event message on discord")
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: utils.ErrorMsg(err),
			})
			if err != nil {
				log.With("err", err).Error("send error message")
			}
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

	// Create cron job
	err = task.CreateTask(task.Option{
		Cron:     opt.Cron,
		Spec:     durationStr,
		TaskName: "d4-event",
		Task:     job,
	})
	if err != nil {
		log.With("err", err).Error("D4 events cron creation")
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.ErrorMsg(cronError),
		})
	}

	log.Infof("init cron schedule to check event diablo IV every %s", duration)
}

func rssParser(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	// Access options in the order provided by the user.
	optionsParam := i.ApplicationCommandData().Options

	// init channel id for go routine
	channelID := i.ChannelID

	// Convert option slice into a map
	var (
		taskName string
		duration string
		rssURL   string
		nbMax    int
	)

	for _, optParam := range optionsParam {
		switch optParam.Name {
		case "name":
			taskName = optParam.StringValue()
		case "time":
			duration = optParam.StringValue()
		case "url":
			rssURL = optParam.StringValue()
		case "nb-last-news":
			nbMax = int(optParam.IntValue())
		}
	}

	// Convert duration to string duration for cron
	durationStr := fmt.Sprintf("@every %s", duration)

	// Check duration format
	_, err := time.ParseDuration(duration)
	if err != nil {
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.ErrorMsg(err),
		})
		if err != nil {
			log.With("err", err).Error("send error message")
		}
		return
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.StringPtr(fmt.Sprintf("Exec %s every %s", taskName, duration)),
	})
	if err != nil {
		log.With("err", err).Error("send message")
	}

	// Init first last GUIDs
	lastGUIDs := make([]string, 0, 5)

	// Job function for Cron
	job := func() {
		log.Debugf("exec %s", taskName)

		listRSSMsg, listGUIDs, err := rss.ParseRSS(rssURL, nbMax, lastGUIDs)
		if err != nil {
			log.With("err", err).With("taskName", taskName).Error("RSS parse")
		}

		// Update variables
		lastGUIDs = listGUIDs
		nbMax = rss.DefaultNbMax

		cpt := len(listRSSMsg)
		if cpt > 0 {
			for idx := cpt - 1; idx >= 0; idx-- {
				_, err := s.ChannelMessageSend(channelID, listRSSMsg[idx])
				if err != nil {
					log.With("err", err).With("taskName", taskName).Error("send error message")
				}
			}
		}

		log.Debugf("exec %s done", taskName)
	}

	err = task.CreateTask(task.Option{
		Cron:     opt.Cron,
		Spec:     durationStr,
		TaskName: taskName,
		Task:     job,
	})
	if err != nil {
		log.With("err", err).With("taskName", taskName).Error("cron creation")
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: utils.ErrorMsg(cronError),
		})
	}

	log.Infof("init cron schedule to exec %s every %s", taskName, duration)
}
