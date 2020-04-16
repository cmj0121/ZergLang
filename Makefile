PROJ_NAME = zerg

CC      := gcc
CFLAGS  := -Iinclude/
LDFLAGS :=

BIN=$(PROJ_NAME)
SRC=$(wildcard src/*.c) $(wildcard src/*/*.c)
OBJ=$(subst .c,.o,$(SRC))

.PHONY: all clean install

all: $(BIN)	# build the binary

help:	# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

clean:	# clean-up the environment
	rm -f $(BIN) $(OBJ)

install: install-syntax	# install into system

VIM_SYNTAX := ~/.vim/syntax
VIM_SRC := $(wildcard vim/*.vim)
VIM_DST := $(subst vim/,$(VIM_SYNTAX)/,$(VIM_SRC))

install-syntax: $(VIM_DST)

$(VIM_SYNTAX)/%.vim: vim/%.vim $(VIM_SYNTAX)
	install -m 644 $< $@

$(VIM_SYNTAX):
	mkdir -p $@

%.o: %.c
	$(CC) -c -o $@ $< $(CFLAGS) $(LDFLAGS)

$(BIN): $(OBJ)
	$(CC) -o $@ $^
