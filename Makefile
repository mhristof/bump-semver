#
# vim:ft=make
#

SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.ONESHELL:

GIT_TAG := $(shell git describe --tags --always --dirty=+)


.PHONY: all
all: ./bin/semver.darwin ./bin/semver.linux

./bin/semver.%: $(shell find ./ -name '*.go')
	GOOS=$* go build -o $@ -ldflags "-X github.com/mhristof/semver/cmd.version=$(GIT_TAG)" main.go

.PHONY: install
install: ./bin/semver.darwin
	rm -f ~/.local/bin/semver
	cp $< ~/.local/bin/semver

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
