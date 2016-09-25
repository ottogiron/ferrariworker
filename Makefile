# Meta info
NAME := ferrariworker
VERSION := 0.6.0
MAINTAINER := "Otto Giron <ottog2486@gmail.com"
SOURCE_URL := https://github.com/ottogiron/ferrariworker.git
DATE := $(shell date -u +%Y%m%d.%H%M%S)
COMMIT_ID := $(shell git rev-parse --short HEAD)
GIT_REPO := $(shell git config --get remote.origin.url)
# Go tools flags
LD_FLAGS := -X 	github.com/ottogiron/ferrariworker/cmd.buildVersion=$(VERSION)
LD_FLAGS += -X github.com/ottogiron/ferrariworker/cmd.buildCommit=$(COMMIT_ID)
LD_FLAGS += -X github.com/ottogiron/ferrariworker/cmd.buildDate=$(DATE)
EXTRA_BUILD_VARS := CGO_ENABLED=0 GOARCH=amd64
SOURCE_DIRS ?= $(shell go list ./... | grep -v /vendor/)
BACKEND_PATH := $(shell go list)/backend/$(BACKEND_NAME)

all: test package-linux package-darwin

build-release: container

lint:
	@go fmt $(SOURCE_DIRS)
	@go vet $(SOURCE_DIRS)

test: lint
	 @go test -v $(SOURCE_DIRS) -cover -bench .

test-backend: lint 
	@go test -v $(BACKEND_PATH) -cover -bench .

cover: 
	gocov test $(SOURCE_DIRS) | gocov-html > coverage.html && open coverage.html

cover-backend:
	gocov test $(BACKEND_PATH) | gocov-html > coverage.html && open coverage.html	

