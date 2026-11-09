.PHONY: build build-image binary test unit-test run run-dev dev-run

SHELL := /bin/bash
APP_NAME = "provencal-le-gaulois"
BINARY_SRC ?= main.go
DOCKER_RUN_FLAGS ?= --rm -d

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

## Run unit tests
test:
	go test -v ./...

## Run unit tests (alias for test)
unit-test: test

## Build docker image
build-image:
	docker build . --tag $(APP_NAME):$(git_tag)

## Build Dev docker image
build-dev-image:
	docker build . --file Dockerfile-dev --tag $(APP_NAME):$(git_tag)-dev

## Run Dev docker image
run-dev:
	docker rm -f $(APP_NAME) 2>/dev/null || true
	docker run $(DOCKER_RUN_FLAGS) --name $(APP_NAME) $(APP_NAME):$(git_tag)-dev

## Run Dev docker image (alias for run-dev)
dev-run: run-dev

## Run docker image
run:
	docker rm -f $(APP_NAME) 2>/dev/null || true
	docker run $(DOCKER_RUN_FLAGS) --name $(APP_NAME) $(APP_NAME):$(git_tag)
