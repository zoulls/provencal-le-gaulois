.PHONY: help build binary

SHELL := /bin/bash
APP_NAME = "provencal-le-gaulois"
BINARY_SRC ?= main.go

git_commit = $(shell git rev-parse HEAD)
git_tag    = $(shell git describe --tags --abbrev=0)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
build_time = $(shell date -u -Iseconds)

go_build_version_flags = "-X main.Version=$(git_tag) -X main.GitBranch=$(git_branch) -X main.BuildTime=$(build_time) -X main.GitCommit=$(git_commit)"

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20
## Show help
# help:
# 	@echo ''
# 	@echo 'Usage:'
# 	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
# 	@echo ''
# 	@echo 'Targets:'
# 	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
# 		helpMessage = match(lastLine, /^## (.*)/); \
# 		if (helpMessage) { \
# 			helpCommand = substr($$1, 0, index($$1, ":")-1); \
# 			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
# 			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
# 		} \
# 	} \
# 	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## Build the app
build: binary build-image

## Build go binary
binary:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -o bin/$(APP_NAME) -ldflags $(go_build_version_flags) $(BINARY_SRC)

## Build docker image
build-image:
	docker build . --tag $(APP_NAME):$(git_tag)