# Provencal le gaulois
Provencal le gaulois is a bot discord written in go.

## Main features
* Create an embed message via a json object ( [json generator](https://leovoel.github.io/embed-visualizer/) )
* Publish on a discord channel the last tweets of users followed
* Use `~help` in a discord channel to see all the command available
* `~helpAdmin` to check redis status and bot updtime

## Setup
* Check and adjust config.json file
* Copy `.env-example` to `.env-prod` and `.env`
* `.env` is used for local run,`.env-prod` is for docker container
* Write Discord bot secret key, Twitter secret key and Discord channels ID in the different env files

## Run in docker
* Generate binary and docker image `make build`
* Run redis container `docker run --name redis -d -p 6379:6379 redis`
* Run bot image link to redis container `docker run --name provencal-le-gaulois -d --link redis:redis provencal-le-gaulois:<latest_tag>`

## Run Localy
* Generate go binary `make binary`
* Place config.json and .env files at the same directory of the binary
* To run in daemon mode
```
$ provencal-le-gaulois &
```

## Logs
The logs are print in the Stdout

## Todo
* Create a command to configure directly in discord the different channels for the tweets push (use a DB or overwrite the config file ?)
* Perceval's quotes
* Auto detect Free game tweet
* Multi config in redis for multi server usage

## Dev Getting Started
* Use the commands `go mod tidy` to pull the vendor
* Update .env values at the root project with Discord bot secret key, Twitter secret key and Discord channels ID etc...
* In the env varibale `CONFIG_FILENAME` put `config-dev` to use the config file `config-dev.json` of config folder
