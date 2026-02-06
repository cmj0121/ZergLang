BIN_DIR := ../../bin
BINARY := $(BIN_DIR)/zerg-bootstrap

.PHONY: all clean test build

all: build

build:				# build the binary
	@mkdir -p $(BIN_DIR)
	go build -o $(BINARY) .

clean:				# clean build artifacts
	@rm -f $(BINARY)

test:				# run tests
	go test ./...
