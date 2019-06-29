PKG := github.com/dh1tw/infractl
COMMITID := $(shell git describe --always --long --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags)

PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

build:
	go build -v -ldflags="-X github.com/dh1tw/infractl/cmd.commitHash=${COMMIT} \
		-X github.com/dh1tw/remoteSwitch/cmd.version=${VERSION}"

# build and strip off the dwraf table. This will reduce the file size
dist:
	go build -v -ldflags="-w -X github.com/dh1tw/infractl/cmd.commitHash=${COMMIT} \
		-X github.com/dh1tw/remoteSwitch/cmd.version=${VERSION}"

.PHONY: build dist