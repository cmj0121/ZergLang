PYTHON := python3.11
VENV   := .venv

.PHONY: all clean test run build install upgrade prologue help

all: 			# default action

clean:			# clean-up environment
	@find . -name '*.pyc' -delete
	@find . -name '__pycache__' -delete
	@rm -rf dist/

test:			# run test
	poetry run pytest -v tests/

run: install	# run in the local environment
	poetry run python zergb/cli.py ../test.zg -o demo.obj

build:			# build the binary/library
	poetry build

install:		# install into the local environment
	poetry install

upgrade:		# upgrade all the necessary packages
	pre-commit autoupdate

prologue: $(VENV)

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(VENV):
	@$(PYTHON) -m venv $(VENV)
