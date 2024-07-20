SRC := $(shell find . -name "*.go" -type f)
BIN := zergb

.PHONY: all clean test run build upgrade help

all: 			# default action

clean:			# clean-up environment

test:			# run test
	@go mod tidy
	@gofmt -s -w $(SRC)
	go test -v ./...

run:			# run in the local environment

build:			# build the binary/library
	@go mod tidy
	go build -o $(BIN) -ldflags "-s -w" cmd/$(BIN)/main.go

upgrade:		# upgrade all the necessary packages

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
