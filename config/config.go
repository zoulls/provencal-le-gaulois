package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// BOT
type Config struct {
	ID           string
	Auth         *AuthConfig
	Redis        *RedisConfig
	Name         string
	Status       string
	PrefixCmd    string
	Twitter      *Twitter
	StatusUpdate *StatusUpdate
	Logger       *Logger
	Discord      *Discord
	startDate    time.Time
}

type AuthConfig struct {
	Secret string
}

type RedisConfig struct {
	URL       string
	Pool      int64
	PingTimer float64
}

// TWITTER
type Twitter struct {
	TwitterFollows []*TwitterFollow
	FollowIDstring string
	Config         *TwitterConfig
}

type TwitterFollow struct {
	Name        string
	Key         string
	ListStr     string
	List        []string
	DiscordChan string
}

type TwitterConfig struct {
	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

type StatusUpdate struct {
	Date    string
	NbUnits int
	Enabled bool
	Every   float64
	Game    string
}

type Logger struct {
	Level string
}

type Discord struct {
	AdminID string
}

// Config singleton
var config *Config

// Check initialized exactly once
var once sync.Once

// Read the config file from the current directory and marshal
// into the conf config struct.
func firstInit() *Config {
	var err error

	configName := os.Getenv("CONFIG_FILENAME")
	if configName == "" {
		configName = "config"
	}

	viper.SetConfigName(configName) // name of config file (without extension)
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	viper.AddConfigPath("./config") // optionally look for config in the working directory

	viper.SetEnvPrefix("plg")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file, %v\n", err))
	}

	if _, err := os.Stat(".env"); err == nil {
		// Load .env var
		err = godotenv.Load()
		if err != nil {
			panic(fmt.Errorf("Error loading .env file, %v\n", err))
		}
	}

	conf := &Config{}

	// Init start date
	conf.startDate = time.Now()

	err = viper.Unmarshal(conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode into config struct, %v\n", err))
	}

	conf.Auth = &AuthConfig{
		Secret: os.Getenv("BOT_SECRET"),
	}

	for _, follow := range conf.Twitter.TwitterFollows {
		if len(conf.Twitter.FollowIDstring) > 0 {
			conf.Twitter.FollowIDstring = conf.Twitter.FollowIDstring + ","
		}
		conf.Twitter.FollowIDstring = conf.Twitter.FollowIDstring + follow.ListStr
		follow.List = strings.Split(follow.ListStr, ",")
		follow.DiscordChan = os.Getenv("DISCORD_CHANNEL_FOR_TWEET_" + follow.Name)
	}

	conf.Discord = &Discord{
		AdminID: os.Getenv("DISCORD_ADMIN_ID"),
	}

	conf.Twitter.Config = &TwitterConfig{
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	// Redis conf
	conf.Redis.URL = os.Getenv("REDIS_URL")

	return conf
}

// Return current config
func GetConfig() *Config {
	once.Do(func() {
		config = firstInit()
	})
	return config
}

func UpdateTwitter(twitterCf Twitter) *Config {
	config.Twitter = &twitterCf
	return config
}

func (c *Config) SetID(id string) {
	c.ID = id
}

func (c *Config) GetStartDate() time.Time {
	return c.startDate
}
