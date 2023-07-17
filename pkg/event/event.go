package event

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

const NextEventStr = "Next event"
const LatestEventStr = "Latest event"
const ActiveStr = ":fire: Ending"
const NAEventStr = "NA"

// iota number and EventsName array need to have the same order
const (
	EventWB = iota
	EventHelltide
	EventLegions
)

var EventsName = []string{"World Boss", "Helltide", "Legion"}

var RegDiscordTime *regexp.Regexp

type EventTimer struct {
	Name   string
	Latest time.Time
	Next   time.Time
}

func init() {
	RegDiscordTime, _ = regexp.Compile("<t:([0-9]*):R>")
}

func EvenTimerInit() []*EventTimer {
	eventTimers := make([]*EventTimer, 3)
	for k, name := range EventsName {
		eventTimers[k] = &EventTimer{
			Name:   name,
			Latest: time.Time{},
			Next:   time.Time{},
		}
	}

	return eventTimers
}

func (et *EventTimer) SetNext(next time.Time) {
	et.Next = next
}

func (et *EventTimer) SetNextTimestamp(ts string) error {
	next, err := utils.UnixStringToTime(ts)
	if err != nil {
		return err
	}
	et.Next = next
	return nil
}

func (et *EventTimer) GetNextTimestamp() string {
	return strconv.Itoa(int(et.Next.Unix()))
}

func (et *EventTimer) SetLatestTimestamp(ts string) error {
	latest, err := utils.UnixStringToTime(ts)
	if err != nil {
		return err
	}
	et.Latest = latest
	return nil
}

func (et *EventTimer) GetLatestTimestamp() string {
	return strconv.Itoa(int(et.Latest.Unix()))
}

func TimerMsg(eventTimers []*EventTimer, newEvent []bool) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField
	for k, eTimer := range eventTimers {
		msgNext := NextEventStr + " "
		msgLatest := LatestEventStr + " "

		if eTimer.Latest.IsZero() {
			msgLatest += NAEventStr
		} else {
			msgLatest += fmt.Sprintf("<t:%s:R>", eTimer.GetLatestTimestamp())
		}

		latestField := discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("-- %s --", eTimer.Name),
			Value:  msgLatest,
			Inline: false,
		}

		if eTimer.Next.IsZero() {
			msgNext += "NA"
		} else {
			if k == EventHelltide && newEvent[EventHelltide] {
				msgNext = ActiveStr + " "
			}
			msgNext += fmt.Sprintf("<t:%s:R>", eTimer.GetNextTimestamp())
		}
		nextField := discordgo.MessageEmbedField{
			Name:   "",
			Value:  msgNext,
			Inline: false,
		}

		fields = append(fields, &latestField, &nextField)
	}

	return discordgo.MessageEmbed{
		Title: ":timer: Events timer :timer:",
		Color: 0xff0000,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/lXhXQzM.png",
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Event timers from d4armory.io",
			IconURL: "https://i.imgur.com/JSbhnHQ.png",
		},
	}
}

func ParseTimerMsg(dMsg *discordgo.Message, eventTimers []*EventTimer) []*EventTimer {
	msg := dMsg.Embeds[0]
	for k, v := range msg.Fields {
		switch k {
		case 0, 1:
			parseField(v.Value, EventWB, eventTimers)
		case 2, 3:
			parseField(v.Value, EventHelltide, eventTimers)
		case 4, 5:
			parseField(v.Value, EventLegions, eventTimers)
		}
	}

	return eventTimers
}

func parseField(fieldValue string, eventType int, eventTimers []*EventTimer) {
	var err error
	if strings.Contains(fieldValue, NextEventStr) || strings.Contains(fieldValue, ActiveStr) {
		found := RegDiscordTime.FindStringSubmatch(fieldValue)

		if len(found) == 2 {
			tsNext := found[1]

			eventTimers[eventType].Next, err = utils.UnixStringToTime(tsNext)
			if err != nil {
				log.Errorf("Error during ParseTimerMsg for next event type %d ts, err: %s", eventType, err)
			}
		}
	}
	if strings.Contains(fieldValue, LatestEventStr) {
		found := RegDiscordTime.FindStringSubmatch(fieldValue)

		if len(found) == 2 {
			tsLatest := found[1]

			eventTimers[eventType].Latest, err = utils.UnixStringToTime(tsLatest)
			if err != nil {
				log.Errorf("Error during ParseTimerMsg for latest event type %d ts, err: %s", eventType, err)
			}
		}
	}
}

func RefreshEventTimers(eventTimers []*EventTimer) ([]*EventTimer, error) {
	// Get data from d4armory.io
	data, err := getD4EventData()
	if err != nil {
		return eventTimers, err
	}

	// Convert D4armoryData to EventTimer
	for k, eTimer := range eventTimers {
		switch k {
		case EventWB:
			eTimer.Latest = time.Unix(int64(data.Boss.Timestamp), 0)
			eTimer.Next = time.Unix(int64(data.Boss.Expected), 0)
		case EventHelltide:
			eTimer.Latest = time.Unix(int64(data.Helltide.Timestamp), 0)
			eTimer.Next = eTimer.Latest.Add(time.Hour*2 + time.Minute*15)
		case EventLegions:
			eTimer.Latest = time.Unix(int64(data.Legion.Timestamp), 0)
			eTimer.Next = time.Unix(int64(data.Legion.Expected), 0)
		}
	}

	return eventTimers, err
}
