package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	cron "github.com/robfig/cron/v3"

	"github.com/zoulls/provencal-le-gaulois/pkg/cmd"
	"github.com/zoulls/provencal-le-gaulois/pkg/event"
)

var (
	s               *discordgo.Session
	c               *cron.Cron
	commands        = cmd.GetApplicationCommand()
	commandHandlers = cmd.GetCommandHandlers()
	// BuildTime is replaced at compile time using ldflags
	BuildTime string
	// Version is replaced at compile time using ldflags
	Version = "Dev"
	// GitBranch is replaced at compile time using ldflags
	GitBranch string
	// GitCommit is replaced at compile time using ldflags
	GitCommit string
)

func init() {
	// Load .env var
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	_, err := os.Stat(envFile)
	if err != nil {
		log.With("err", err).Fatal("no .env file")
	}
	err = godotenv.Load(envFile)
	if err != nil {
		log.With("err", err).Fatal("loading .env file")
	}

	// Init log level
	lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.With("err", err).Fatal("bad LOG_LEVEL value")
	}
	log.SetLevel(lvl)

	// Init event host API
	event.InitHost()

	// Init discord session
	s, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.With("err", err).Fatal("invalid bot parameters")
	}

	// Init cron
	c = cron.New()

	// Init command option
	opt := cmd.Option{
		Cron:       c,
		LaunchTime: time.Now(), // Set launch time for uptime
		BuildInfo: cmd.BuildInfo{
			Version:   Version,
			BuildTime: BuildTime,
			GitBranch: GitBranch,
			GitCommit: GitCommit,
		},
	}

	// Declare ApplicationCommandData
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i, opt)
		}
	})
}

func main() {
	// Build logs
	log.With("version", Version, "git branch", GitBranch, "git commit", GitCommit, "build time", BuildTime).Info("build info")
	// Bot env
	log.Infof("env: %s", os.Getenv("BOT_ENV"))

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.With("err", err).Fatal("cannot open the session")
	}

	log.Info("adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, os.Getenv("SERVER_ID"), v)
		if err != nil {
			log.With("command", v.Name, "err", err).Error("cannot create command")
		}
		registeredCommands[i] = cmd
	}

	// Set bot status
	err = s.UpdateGameStatus(0, Version)
	if err != nil {
		log.With("err", err).Error("cannot set status")
	}

	// Close discord connection
	defer s.Close()

	// Stop cron scheduler
	defer c.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("press Ctrl+C to exit")
	<-stop

	log.Info("removing commands...")
	// We need to fetch the commands, since deleting requires the command ID.
	// We are doing this from the returned commands on line 375, because using
	// this will delete all the commands, which might not be desirable, so we
	// are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, os.Getenv("SERVER_ID"), v.ID)
		if err != nil {
			log.With("command", v.Name, "err", err).Error("cannot delete command")
		}
	}

	log.Info("gracefully shutting down.")
}
