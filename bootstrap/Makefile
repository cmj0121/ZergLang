SRC := $(shell find . -name "*.go" -type f)
BIN := zergb

.PHONY: all clean test run build upgrade help

all: 			# default action

clean:			# clean-up environment
	rm -f $(BIN)

test:			# run test
	@go mod tidy
	@gofmt -s -w $(SRC)
	@golangci-lint run --fix
	go test ./...

run:			# run in the local environment

build: $(BIN)	# build the binary/library
	@go mod tidy

upgrade:		# upgrade all the necessary packages

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(BIN): $(SRC)
	go generate ./...
	go build -o $(BIN) -ldflags "-s -w" cmd/$(BIN)/*.go
