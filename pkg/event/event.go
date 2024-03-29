package event

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

const (
	NextEventStr   = "Next event"
	LatestEventStr = "Latest event"
	ActiveStr      = ":fire: Ending"
	NAEventStr     = "NA"
	SoonStr        = ":rotating_light:"
)

// iota number and EventsName array need to have the same order
const (
	EventWB = iota
	EventHelltide
	EventLegions
)

var (
	EventsName          = []string{"World Boss", "Helltide", "Legion"}
	RegDiscordTime      *regexp.Regexp
	ActiveHelltideTimer time.Duration
	NextHelltideTimer   time.Duration
)

type EventTimer struct {
	Name    string
	Latest  time.Time
	Next    time.Time
	Active  bool
	Soon    bool
	Updated time.Time
}

func init() {
	RegDiscordTime, _ = regexp.Compile("<t:([0-9]*):R>")
}

func InitEventTimer() []*EventTimer {
	// Init Helltide timer
	hActiveDuration, err := time.ParseDuration(os.Getenv("D4_HELLTIDE_ACTIVE_DURATION"))
	if err != nil {
		log.With("err", err).Error("no D4 Helltide active duration set")
	}
	ActiveHelltideTimer = hActiveDuration

	hNextDuration, err := time.ParseDuration(os.Getenv("D4_HELLTIDE_NEXT_DURATION"))
	if err != nil {
		log.With("err", err).Error("no D4 Helltide active duration set")
	}
	NextHelltideTimer = hNextDuration

	// Init EventTimer slice
	eventTimers := make([]*EventTimer, 3)
	for k, name := range EventsName {
		eventTimers[k] = &EventTimer{
			Name: name,
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

			nTimer := eTimer.Latest.Add(NextHelltideTimer)
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
		prefix := ""
		suffix := ""

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
				// replace start msgNext
				msgNext = ActiveStr + " "
			}
			if eTimer.Soon {
				// redefine prefix and suffix
				prefix = SoonStr + " "
				suffix = " " + SoonStr
			}
			msgNext += fmt.Sprintf("%s<t:%s:R>%s", prefix, eTimer.GetNextTimestamp(), suffix)
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
				log.With("eventType", eventType, "err", err).Error("error during ParseTimerMsg for next event")
			}
		}
	}
	if strings.Contains(fieldValue, LatestEventStr) {
		found := RegDiscordTime.FindStringSubmatch(fieldValue)

		if len(found) == 2 {
			tsLatest := found[1]

			eventTimers[eventType].Latest, err = utils.UnixStringToTime(tsLatest)
			if err != nil {
				log.With("eventType", eventType, "err", err).Error("error during ParseTimerMsg for latest event")
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
		refresh := false
		// retrieve diff time between now and next timer
		diff := eTimer.Next.Sub(now)

		// refresh timer for WB and Legion before 10 minutes expire
		if k != EventHelltide && diff.Minutes() <= 10 && diff.Minutes() > 0 {
			log.Debugf("refresh next timer event %s", eTimer.Name)

			if !getData {
				// get data from d4armory.io only if a next timer expire
				data, err = getD4EventData()
				if err != nil {
					return eventTimers, err
				}
				getData = true
			}

			switch k {
			case EventWB:
				eTimer.Next = time.Unix(int64(data.Boss.Expected), 0)
			case EventLegions:
				eTimer.Next = time.Unix(int64(data.Legion.Expected), 0)
			}

			refresh = true
		}

		if refresh {
			// refresh diff with new next timer
			diff = eTimer.Next.Sub(now)
		}

		if diff.Minutes() <= 5 && diff.Minutes() > 0 {
			// soon flag
			eTimer.Soon = true
		}

		// check if a next timer is expired
		if diff.Minutes() <= 0 {
			log.Debugf("refresh timer event %s", eTimer.Name)

			if !getData {
				// get data from d4armory.io only if a next timer expire
				data, err = getD4EventData()
				if err != nil {
					return eventTimers, err
				}
				getData = true
			}

			// reset soon flag
			eTimer.Soon = false

			switch k {
			case EventWB:
				eTimer.Latest = eTimer.Next
				eTimer.Next = time.Unix(int64(data.Boss.Expected), 0)
			case EventHelltide:
				// re-sync latest timer
				eTimer.Latest = time.Unix(int64(data.Helltide.Timestamp), 0)

				// toggle active boolean
				eTimer.Active = !eTimer.Active

				// init to normal timer
				nTimer := eTimer.Latest.Add(NextHelltideTimer)
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
