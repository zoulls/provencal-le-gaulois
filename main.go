package main

import (
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
	"github.com/zoulls/provencal-le-gaulois/pkg/discord"
	"github.com/zoulls/provencal-le-gaulois/pkg/logger"
	"github.com/zoulls/provencal-le-gaulois/pkg/redis"
	"github.com/zoulls/provencal-le-gaulois/pkg/status"
	"github.com/zoulls/provencal-le-gaulois/pkg/twitter"
)

var (
	// BuildTime is replaced at compile time using ldflags
	BuildTime string
	// Version is replaced at compile time using ldflags
	Version string
	// GitBranch is replaced at compile time using ldflags
	GitBranch string
	// GitCommit is replaced at compile time using ldflags
	GitCommit string

	// Bot identification
	botID string
)

func main() {

	// Build logs
	logger.Log().Infof("go version: %v, git branch: %v, git commit: %v, build time: %v", Version, GitBranch, GitCommit, BuildTime)

	// Init Config
	conf := config.GetConfig()

	logger.Log().Infof("Env: %s", os.Getenv("BOT_ENV"))

	logger.Log().Infof("Logger level: %s", logger.Log().Level.String())

	// Redis client
	rClient := redis.NewClient()

	// Sync Twitter follows list
	tConf, err := twitter.SyncList(rClient, *conf.Twitter)
	if err != nil {
		logger.Log().Errorf("Error during Twitter sync list, %v", err)
	} else {
		conf = config.UpdateTwitter(tConf)
	}

	// Discord client
	ds, err := discordgo.New("Bot " + conf.Auth.Secret)
	errCheck("error creating discord session", err)
	user, err := ds.User("@me")
	errCheck("error retrieving account", err)
	botID = user.ID

	// Get default status
	sClient := status.New(conf, rClient)
	defaultStatus, err := sClient.Last(true)
	if err != nil {
		logger.Log().Errorf("Error during status init, %v", err)
	}

	ds.AddHandler(commandHandler)
	ds.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateGameStatus(0, defaultStatus)
		if err != nil {
			logger.Log().Errorf("Error attempting to set my status, %v", err)
		}
		servers := discord.State.Guilds
		logger.Log().Infof("%s has started on %d servers", conf.Name, len(servers))
	})

	err = ds.Open()
	errCheck("Error opening connection to Discord", err)
	defer ds.Close()

	twitter.StreamTweets(ds, sClient, rClient)

	<-make(chan struct{})
	logger.Log().Errorf("%s stop to %s", conf.Name, conf.Status)
}

func errCheck(msg string, err error) {
	if err != nil {
		logger.Log().Errorf("%s: %v", msg, err)
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// init var
	user := m.Author

	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	res, err := discord.GetReply(s, m)
	if err != nil {
		logger.Log().Errorf("Message send error: %v", err)
	}
	if res != nil {
		_, err = s.ChannelMessageSendComplex(m.ChannelID, res)
		errCheck("Error during bot reply", err)
	}
}
