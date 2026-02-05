SUBDIR :=
VIM_PREFIX := $(HOME)/.vim

.PHONY: all clean test run build upgrade help install-vim uninstall-vim $(SUBDIR)

all: $(SUBDIR) 		# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean: $(SUBDIR)	# clean-up environment
	@find . -name '*.sw[po]' -delete

test:				# run test

run:				# run in the local environment

build:				# build the binary/library

install-vim:		# install vim syntax files
	@mkdir -p $(VIM_PREFIX)/syntax $(VIM_PREFIX)/ftdetect
	@cp docs/vim/syntax/zerg.vim $(VIM_PREFIX)/syntax/
	@cp docs/vim/ftdetect/zerg.vim $(VIM_PREFIX)/ftdetect/
	@echo "Installed zerg.vim to $(VIM_PREFIX)"

uninstall-vim:		# uninstall vim syntax files
	@rm -f $(VIM_PREFIX)/syntax/zerg.vim
	@rm -f $(VIM_PREFIX)/ftdetect/zerg.vim
	@echo "Removed zerg.vim from $(VIM_PREFIX)"

upgrade:			# upgrade all the necessary packages
	pre-commit autoupdate

help:				# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(SUBDIR):
	$(MAKE) -C $@ $(MAKECMDGOALS)
