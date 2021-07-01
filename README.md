# Provencal le gaulois
Provencal le gaulois is a bot discord written in go.

## Main features
* Create an embed message via a json object ( [json generator](https://leovoel.github.io/embed-visualizer/) )
* Publish on a discord channel the last tweets of users followed
* Use `~help` in a discord channel to see all the command available

## Installing
* Generate go binary `go build main.go`
* Copy config-prod.json to config.json
* Copy .env-dev to .env
* Write Discord bot secret key, Twitter secret key and Discord channels ID in the .env file
* Place config.json and .env files at the same directory of the binary

## Launch
daemon mode
```
$ provencal-le-gaulois &
```

## Logs
During the first launch, a log file *info.log* is generated at the same directory of the binary.
The log file is wipe at every launch.

## Todo
* Create a command to configure directly in discord the different channels for the tweets push (use a DB or overwrite the config file ?)
* Add a logger with error levels
* Discord Embedded message with Twitter Websites Overview
* Perceval's quotes


## Dev Getting Started
* Use the commands `go mod tidy` to pull the vendor
* Copy config-dev.json into config.json in folder /config
* Copy .env-dev into .env at the root project
* Write Discord bot secret key, Twitter secret key and Discord channels ID in .env file