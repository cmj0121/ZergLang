NVIM_SITE := $(HOME)/.local/share/nvim/site

.PHONY: install uninstall

install:			# install vim/neovim syntax files
	@install -d $(NVIM_SITE)/syntax $(NVIM_SITE)/ftdetect $(NVIM_SITE)/indent
	@install -m 644 vim/syntax/zerg.vim $(NVIM_SITE)/syntax/
	@install -m 644 vim/ftdetect/zerg.vim $(NVIM_SITE)/ftdetect/
	@install -m 644 vim/indent/zerg.vim $(NVIM_SITE)/indent/
	@echo "Installed Zerg syntax files to $(NVIM_SITE)"

uninstall:			# remove vim/neovim syntax files
	@rm -f $(NVIM_SITE)/syntax/zerg.vim
	@rm -f $(NVIM_SITE)/ftdetect/zerg.vim
	@rm -f $(NVIM_SITE)/indent/zerg.vim
	@echo "Removed Zerg syntax files from $(NVIM_SITE)"
