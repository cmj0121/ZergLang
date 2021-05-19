/* Copyright (C) 2020-2021 cmj <cmj@cmj.tw>. All right reserved. */
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>
#include <sys/mman.h>

#include "zerg.h"

typedef struct _tag_lexer_ {
	int fd;
	char *ptr;
	size_t size;
	size_t cur;
} Lexer;

static int open_lexer(Lexer *lexer, const char *filepath) {
	int ret = -1;

	if (0 > (lexer->fd = open(filepath, O_RDONLY))) {
		_D(WARN, "cannot open file '%s': %s", filepath, strerror(errno));
		goto END;
	}

	// get the file total size
	struct stat st;
	if (0 > fstat(lexer->fd, &st)) {
		_D(WARN, "cannot stat file: %s", strerror(errno));
		goto END;
	}

	lexer->cur = 0;
	lexer->size = st.st_size;
	/* load the source code into memory and process as long char array */
	if (MAP_FAILED == (lexer->ptr = mmap(NULL, lexer->size, PROT_READ, MAP_PRIVATE, lexer->fd, 0))) {
		_D(WARN, "cannot load into memory: %s", strerror(errno));
		goto END;
	}

	_D(INFO, "load %s into memory with size %zu", filepath, lexer->size);
	ret = 0;
END:
	return ret;
}

static void close_lexer(Lexer *lexer) {
	close(lexer->fd);
	if (lexer->ptr) munmap(lexer->ptr, lexer->size);
	return;
}

static int next_token(Lexer *lexer, char *token, size_t token_len) {
	int len = 0;

	if (lexer->cur >= lexer->size) {
		_D(INFO, "end-of-file");
		return -1;
	}

	for (; lexer->cur < lexer->size; ++lexer->cur) {
		switch (lexer->ptr[lexer->cur]) {
			case ' ': case '\t': case '\n': case '\r':
				// get next token
				lexer->cur ++;
				goto END;
			default:
				if (len == token_len) {
					_D(CRIT, "not support token size > %zd", token_len);
					return -1;
				}

				token[len++] = lexer->ptr[lexer->cur];
				break;
		}
	}

END:
	token[len] = '\0';
	return len;
}

// parse and compile source file
int compile(const char *filepath) {
	int ret = -1;
	Lexer lexer = {
		.fd = -1,
		.ptr = NULL,
	};

	if (0 > open_lexer(&lexer, filepath)) {
		_D(CRIT, "cannot open lexer");
		goto END;
	}	

	char token[MAX_TOKEN_LEN] = {0};
	while (0 <= next_token(&lexer, token, MAX_TOKEN_LEN)) {
		_D(WARN, "throw token '%s'", token);
	}

	ret = 0;
END:
	close_lexer(&lexer);
	return ret;
}
