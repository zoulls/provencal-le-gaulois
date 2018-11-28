package main

import (
	"bitbucket.org/zoulls/provencal-le-gaulois/config"
	"bitbucket.org/zoulls/provencal-le-gaulois/pkg/logger"
	"bitbucket.org/zoulls/provencal-le-gaulois/pkg/reply"
	"bitbucket.org/zoulls/provencal-le-gaulois/pkg/twitter"
	"github.com/bwmarrin/discordgo"
)

var (
	botID string
)

func main() {
	config := config.GetConfig()
	discord, err := discordgo.New("Bot " + config.Auth.Secret)
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, config.Status)
		if err != nil {
			logger.Log.Printf("Error attempting to set my status\n")
		}
		servers := discord.State.Guilds
		logger.Log.Printf("%s has started on %d servers\n", config.Name, len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	twitter.StreamTweets(discord)

	<-make(chan struct{})
	logger.Log.Printf("%s stop to %s\n", config.Name, config.Status)
}

func errCheck(msg string, err error) {
	if err != nil {
		logger.Log.Printf("%s: %+v\n", msg, err)
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

	// DEBUG
	// fmt.Printf("Message: %+v || From: %s\n", m.Message, m.Author)

	res, err := reply.GetReply(s, m)
	if err != nil {
		logger.Log.Printf("Message send error: %+v\n", err)
	}
	if res != nil {
		_, err = s.ChannelMessageSendComplex(m.ChannelID, res)
		errCheck("Error during bot reply", err)
	}
}
