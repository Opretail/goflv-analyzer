.PHONY: cgo

SERVER_NAME     := goflv-analyzer
TARGET_DIR      := .
BUILD_NAME      := flvanalyzer
BUILD_VERSION   := $(shell date "+%Y%m%d.%H%M%S")
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )


all: release

release:
	go build -ldflags \
	"-s -w \
	-X '${SERVER_NAME}/config.Version=${BUILD_VERSION}' \
	-X '${SERVER_NAME}/config.BuildTime=${BUILD_TIME}' \
	-X '${SERVER_NAME}/config.CommitID=${COMMIT_SHA1}' \
	" -o ${TARGET_DIR}/${BUILD_NAME}