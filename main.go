package main

import (
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/reply"
	"github.com/zoulls/provencal-le-gaulois/pkg/twitter"
)

var (
	botID string
)

func main() {
	logger.Log.Infof("Env: %s", os.Getenv("BOT_ENV"))

	conf := config.GetConfig()
	discord, err := discordgo.New("Bot " + conf.Auth.Secret)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, conf.Status)
		if err != nil {
			logger.Log.Errorf("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		logger.Log.Infof("%s has started on %d servers", conf.Name, len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	twitter.StreamTweets(discord)

	<-make(chan struct{})
	logger.Log.Errorf("%s stop to %s", conf.Name, conf.Status)
}

func errCheck(msg string, err error) {
	if err != nil {
		logger.Log.Errorf("%s: %+v", msg, err)
		panic(err)
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// init var
	user := m.Author

	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	res, err := reply.GetReply(s, m)
	if err != nil {
		logger.Log.Errorf("Message send error: %+v", err)
	}
	if res != nil {
		_, err = s.ChannelMessageSendComplex(m.ChannelID, res)
		errCheck("Error during bot reply", err)
	}
}
