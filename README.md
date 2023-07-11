# Provencal le gaulois
Provencal le gaulois is a bot discord written in go.

## Main features
* Use `/` in discord to see all the commands available
* List and delete messages
* Publish on a discord channel the latest tweets of Twitter list
* Event timer for Diablo IV

## Setup
* Copy `.env` to `.env-prod` file
* Adjust `.env-prod` file with all info and token (Discord bot secret key, Twitter secret key ...)

## Run in docker
* Generate binary and docker image `make build`
* Run bot image link to redis container `docker run --name provencal-le-gaulois -d provencal-le-gaulois:<latest_tag>`

## Dev Getting Started
* Copy `.env` to `.env-dev` file
* Adjust `.env-dev` file with all development tokens
* Use the commands `go mod tidy` to pull the vendor
* Generate binary and docker image for development with `make dev`

## Run Locally
* Generate go binary `make binary`
* Place `.env` files at the same directory of the binary
* To run in daemon mode
```
$ provencal-le-gaulois &
```

## Logs
The logs are print in the Stdout
