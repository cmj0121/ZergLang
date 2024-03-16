.PHONY: all clean test run build upgrade install help

all: 			# default action

clean:			# clean-up environment

test:			# run test

run:			# run in the local environment

build:			# build the binary/library

upgrade:		# upgrade all the necessary packages

install:		# install in the local system
	install -m644 zerg.vim ~/.config/nvim/syntax/zerg.vim

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
