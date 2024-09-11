SUBDIR := bootstrap docs

.PHONY: all clean test run build upgrade install help $(SUBDIR)

all: $(SUBDIR) 		# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean: $(SUBDIR)	# clean-up environment
	@find . -name '*.sw[po]' -delete

test:				# run test

run:				# run in the local environment

build:				# build the binary/library

upgrade:			# upgrade all the necessary packages
	pre-commit autoupdate

install:			# install in the local system

help:				# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

all clean test run build install: $(SUBDIR)
$(SUBDIR):
	$(MAKE) -C $@ $(MAKECMDGOALS)
