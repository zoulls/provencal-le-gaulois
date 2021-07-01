package discord

import (
	"regexp"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"

	"github.com/zoulls/provencal-le-gaulois/config"
)

func help() *discordgo.MessageEmbed {
	config := config.GetConfig()

	fileds := make([]*discordgo.MessageEmbedField, 0)
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "ping",
		Value:  "Ping pong party :ping_pong:",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "help",
		Value:  "Just what you read at this moment",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "embedGen",
		Value:  "URL for Message Embed Generator",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "embed <json>",
		Value:  "Create a embed message from json of Embed Generator",
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title:       ":information_source: Help",
		Description: "List of commands",
		Color:       30935,
		Fields:      fileds,
	}
}

func helpAdmin() *discordgo.MessageEmbed {
	config := config.GetConfig()

	fileds := make([]*discordgo.MessageEmbedField, 0)
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "updateStatus",
		Value:  "Force update bot status",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "statusLastSync",
		Value:  "Last timestamp for bot status sync",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "redisInfo",
		Value:  "Redis config info",
		Inline: false,
	})
	fileds = append(fileds, &discordgo.MessageEmbedField{
		Name:   config.PrefixCmd + "twitterFollows",
		Value:  "List of Twitter follows IDs for store and game news",
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title:       ":information_source: Help Admin",
		Description: "List of Admin commands",
		Color:       30935,
		Fields:      fileds,
	}
}

func embedGenerator() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Message Embed Generator",
		Description: "Visualizer and validator for Discord message embed https://leovoel.github.io/embed-visualizer/",
		Color:       7506902,
		URL:         "https://leovoel.github.io/embed-visualizer/",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/ZimNc57.png",
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/jqrYK44.png",
		},
	}
}

// EmbedFromTweet return an embed message with the tweet data
func EmbedFromTweet(t *anaconda.Tweet) *discordgo.MessageEmbed {
	tweetTime, _ := t.CreatedAtTime()
	timestamp := tweetTime.UTC().Format(time.RFC3339)

	message := &discordgo.MessageEmbed{
		Description: t.FullText,
		Color:       1811438,
		Timestamp:   timestamp,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    t.User.Name,
			IconURL: t.User.ProfileImageURL,
			URL:     "https://twitter.com/" + t.User.ScreenName + "/status/" + t.IdStr,
		},
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: "https://images-ext-1.discordapp.net/external/bXJWV2Y_F3XSra_kEqIYXAAsI3m1meckfLhYuWzxIfI/https/abs.twimg.com/icons/apple-touch-icon-192x192.png",
			Text:    "@" + t.User.ScreenName,
		},
	}

	if len(t.Entities.Urls) > 0 {
		for _, url := range t.Entities.Urls {
			urlReplace := "[" + url.Display_url + "](" + url.Expanded_url + ")"
			message.Description = tweetReplace(message.Description, url.Url, urlReplace)
		}
	}

	if len(t.Entities.Media) > 0 {
		media := t.Entities.Media[0]
		message.Image = &discordgo.MessageEmbedImage{
			URL: media.Media_url,
		}
		message.Description = tweetReplace(message.Description, media.Url, "")
	}

	return message
}

// URLFromTweet return the tweet url
func URLFromTweet(t *anaconda.Tweet) string {
	return "https://twitter.com/" + t.User.ScreenName + "/status/" + t.IdStr
}

func tweetReplace(str string, strSearch string, strRplace string) string {
	r, _ := regexp.Compile(strSearch)
	return r.ReplaceAllString(str, strRplace)
}
