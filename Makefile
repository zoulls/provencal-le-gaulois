.PHONY: build build-image binary

SHELL := /bin/bash
APP_NAME = "provencal-le-gaulois"
BINARY_SRC ?= main.go

git_commit = $(shell git rev-parse HEAD)
git_tag    = $(shell git describe --tags --abbrev=0)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
build_time = $(shell date -u -Iseconds)

go_build_version_flags = "-X main.Version=$(git_tag) -X main.GitBranch=$(git_branch) -X main.BuildTime=$(build_time) -X main.GitCommit=$(git_commit)"

# Include help command
help.mk:
include help.mk

## Build app and docker image
build: binary build-image

## Build app and docker dev image
dev: binary build-dev-image

## Build go binary
binary:
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -o bin/$(APP_NAME) -ldflags $(go_build_version_flags) $(BINARY_SRC)

## Build docker image
build-image:
	docker build . --tag $(APP_NAME):$(git_tag)

## Build Dev docker image
build-dev-image:
	docker build . --file Dockerfile-dev --tag $(APP_NAME):$(git_tag)-dev
