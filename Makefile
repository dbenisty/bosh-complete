BUILD_TARGET ?= *.go
APP_NAME := bosh-complete
OUTPUT_NAME ?= $(APP_NAME)
SHELL := /bin/bash
COMMIT_HASH := $(shell git log --pretty='format:%h' -n 1)
DIRTY_LINE := $(shell git diff --shortstat 2> /dev/null | tail -n1)
ifneq ("$(DIRTY_LINE)", "")
  DIRTY := +
endif
VERSION ?= development
LDFLAGS := -X "github.com/thomasmmitchell/bosh-complete/version.Version=$(VERSION)-$(COMMIT_HASH)$(DIRTY)"
BUILD := go build -v -ldflags='$(LDFLAGS)' -o $(OUTPUT_NAME) $(BUILD_TARGET)

.PHONY: build darwin linux all clean
.DEFAULT: build
build:
	@echo $(VERSION)-$(COMMIT_HASH)$(DIRTY)
	GOOS=$(GOOS) GOARCH=amd64 $(BUILD)

darwin:
	GOOS=darwin OUTPUT_NAME=$(APP_NAME)-darwin VERSION="$(VERSION)" $(MAKE)

linux:
	GOOS=linux OUTPUT_NAME=$(APP_NAME)-linux VERSION="$(VERSION)" $(MAKE)

all: darwin linux

clean:
	rm -f $(APP_NAME) $(APP_NAME)-darwin $(APP_NAME)-linux
