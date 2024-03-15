SUBDIR := vim

.PHONY: all clean test run build upgrade install help $(SUBDIR)

all: 			# default action

clean:			# clean-up environment

test:			# run test

run:			# run in the local environment

build:			# build the binary/library

upgrade:		# upgrade all the necessary packages

install:		# install in the local system

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

all clean test run build install: $(SUBDIR)
$(SUBDIR):
	$(MAKE) -C $@ $(MAKECMDGOALS)
