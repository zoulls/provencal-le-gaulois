package cmd

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

func uptime(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.StringPtr("uptime: " + utils.HumanizeDuration(time.Since(opt.LaunchTime))),
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}

func debug(s *discordgo.Session, i *discordgo.InteractionCreate) {
	active := i.ApplicationCommandData().Options[0].BoolValue()
	if active {
		log.SetLevel(log.DebugLevel)
		log.Info("enable debug log")
	} else {
		log.SetLevel(log.InfoLevel)
		log.Info("disable debug log")
	}

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: utils.StringPtr(fmt.Sprintf("debug: %t", active)),
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

	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embedsMsg,
	})
	if err != nil {
		log.With("err", err).Error("send error message")
	}
}
