package cmd

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/zoulls/provencal-le-gaulois/pkg/utils"
)

func uptime(s *discordgo.Session, i *discordgo.InteractionCreate, opt Option) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "uptime: " + utils.HumanizeDuration(time.Since(opt.LaunchTime)),
		},
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
