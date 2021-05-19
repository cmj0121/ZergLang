/* Copyright (C) 2020-2021 cmj <cmj@cmj.tw>. All right reserved. */
#ifndef _ZERG_H
#	define _ZERG_H

// the project name
#define PROJ_NAME "zerg"
// the project version meta
#define MAJOR 0
#define MINOR 0
#define MACRO 0

#define MAX_TOKEN_LEN 64

// syntax-sugar for the debug message with log level
extern int verbose;

#define _D(lv, msg, ...) \
	do {																			\
		if (lv <= verbose)															\
		fprintf(stderr, "[%s L#%d] " msg "\n", __FILE__, __LINE__, ##__VA_ARGS__);	\
	} while(0)

typedef enum {
	CRIT = 0,
	WARN,
	INFO,
	DEBUG,
} LOG_LEVEL;

#include "zerg_lib.h"

#endif /* _ZERG_H */
