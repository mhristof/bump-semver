#
# vim:ft=make
#

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:

GIT_REF := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git name-rev --tags --name-only $(GIT_REF))

.PHONY: all
all: ./bin/semver.darwin ./bin/semver.linux

./bin/semver.%: $(shell find ./ -name '*.go')
	GOOS=$* go build -o $@ -ldflags "-X github.com/mhristof/semver/cmd.version=$(GIT_TAG)+$(GIT_REF)" main.go

.PHONY: fast-test
fast-test:  ## Run fast tests
	go test ./... -tags fast

.PHONY: test
test:	## Run all tests
	go test ./...

.PHONY: clean
clean:
	rm -rf bin/semver.*

.PHONY: help
help:           ## Show this help.
	@grep '.*:.*##' Makefile | grep -v grep  | sort | sed 's/:.*## /:/g' | column -t -s:
