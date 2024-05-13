package cmd

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type Option struct {
	Cron       *cron.Cron
	LaunchTime time.Time
	BuildInfo  BuildInfo
}

type BuildInfo struct {
	Version   string
	BuildTime string
	GitBranch string
	GitCommit string
}

// default error message
var cronError = errors.New("cron creation")

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
					Description: "Amount of message parsed",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "before-url",
					Description: "List message before message URL",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "after-url",
					Description: "List message after message URL",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author-id",
					Description: "Filter only messages with this author ID",
					Required:    false,
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
					Description: "Amount of message parsed",
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
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author-id",
					Description: "Filter only messages with this author ID",
					Required:    false,
				},
			},
		},
		{
			Name:        "d4event",
			Description: "Check diablo IV event messages, check every X minutes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "Time between each execution, duration format (https://pkg.go.dev/time#ParseDuration)",
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
			Name:        "rss",
			Description: "Report rss event in the channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Task name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "Time between each execution, duration format (https://pkg.go.dev/time#ParseDuration)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "RSS url",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "last-guid",
					Description: "last GUID message",
					Required:    false,
				},
			},
		},
		{
			Name:        "auto-clean",
			Description: "Auto clean messages every each time (batch of 100 messages max)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Task name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "Time between each execution, duration format (https://pkg.go.dev/time#ParseDuration)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "expiration",
					Description: "Expiration time of message, duration format (https://pkg.go.dev/time#ParseDuration)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "author-id",
					Description: "Filter only messages with this author ID",
					Required:    false,
				},
			},
		},
	}
}

func GetCommandHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate, Option) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option){
		"placeholder": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "placeholder"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			placeholder(s, i)
			log.Debugf("end cmd %s", cmdName)
		},
		"uptime": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "uptime"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			uptime(s, i, opt)
			log.Debugf("end cmd %s", cmdName)
		},
		"version": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "version"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			version(s, i, opt)
			log.Debugf("end cmd %s", cmdName)
		},
		"list": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "list"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			list(s, i)
			log.Debugf("end cmd %s", cmdName)
		},
		"debug": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "debug"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			debug(s, i)
			log.Debugf("end cmd %s", cmdName)
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "delete"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			delete(s, i)
			log.Debugf("end cmd %s", cmdName)
		},
		"d4event": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "d4event"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			d4Event(s, i, opt)
			log.Debugf("end cmd %s", cmdName)
		},
		"rss": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "rss"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			rssParser(s, i, opt)
			log.Debugf("end cmd %s", cmdName)
		},
		"auto-clean": func(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
			cmdName := "auto-clean"
			log.Debugf("received cmd %s", cmdName)
			loadingMessage(s, i)
			autoClean(s, i, opt)
			log.Debugf("end cmd %s", cmdName)
		},
	}
}
