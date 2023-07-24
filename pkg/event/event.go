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

const ActiveHelltideTimer = time.Hour
const NormalHelltideTimer = time.Hour*2 + time.Minute*15

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
	Active bool
}

func init() {
	RegDiscordTime, _ = regexp.Compile("<t:([0-9]*):R>")
}

func InitEventTimer() []*EventTimer {
	eventTimers := make([]*EventTimer, 3)
	for k, name := range EventsName {
		eventTimers[k] = &EventTimer{
			Name:   name,
			Latest: time.Time{},
			Next:   time.Time{},
			Active: false,
		}
	}

	return eventTimers
}

func PopulateEventTimer(eventTimers []*EventTimer, activeHelltide bool) ([]*EventTimer, error) {
	// Get data from d4armory.io only if a next timer expire
	data, err := getD4EventData()
	if err != nil {
		return eventTimers, err
	}

	for k, eTimer := range eventTimers {
		switch k {
		case EventWB:
			eTimer.Latest = time.Unix(int64(data.Boss.Timestamp), 0)
			eTimer.Next = time.Unix(int64(data.Boss.Expected), 0)
		case EventHelltide:
			eTimer.Active = activeHelltide
			eTimer.Latest = time.Unix(int64(data.Helltide.Timestamp), 0)

			nTimer := eTimer.Latest.Add(NormalHelltideTimer)
			if eTimer.Active {
				nTimer = eTimer.Latest.Add(ActiveHelltideTimer)
			}
			eTimer.Next = nTimer
		case EventLegions:
			eTimer.Latest = time.Unix(int64(data.Legion.Timestamp), 0)
			eTimer.Next = time.Unix(int64(data.Legion.Expected), 0)
		}
	}

	return eventTimers, err
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

func TimerMsg(eventTimers []*EventTimer) discordgo.MessageEmbed {
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
			if k == EventHelltide && eventTimers[EventHelltide].Active {
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
	// Init local value
	now := time.Now()
	var getData bool
	var data *d4armoryData
	var err error

	// Convert D4armoryData to EventTimer
	for k, eTimer := range eventTimers {
		// retrieve diff time between now and next timer
		diff := now.Sub(eTimer.Next)
		// check if a next timer is before 4 minutes or expired
		if diff.Minutes() <= 4 {
			log.Debugf("refresh timer event %s", eTimer.Name)

			if !getData {
				// Get data from d4armory.io only if a next timer expire
				data, err = getD4EventData()
				if err != nil {
					return eventTimers, err
				}
				getData = true
			}

			switch k {
			case EventWB:
				eTimer.Latest = eTimer.Next
				eTimer.Next = time.Unix(int64(data.Boss.Expected), 0)
			case EventHelltide:
				// only store next event if not active (next time will be active timer)
				if !eTimer.Active {
					eTimer.Latest = eTimer.Next
				}

				// toggle active boolean
				eTimer.Active = !eTimer.Active

				// init to normal timer
				nTimer := eTimer.Latest.Add(NormalHelltideTimer)
				if eTimer.Active {
					// if active change for active timer
					nTimer = eTimer.Latest.Add(ActiveHelltideTimer)
				}
				eTimer.Next = nTimer
			case EventLegions:
				eTimer.Latest = eTimer.Next
				eTimer.Next = time.Unix(int64(data.Legion.Expected), 0)
			}
		}
	}

	return eventTimers, err
}
