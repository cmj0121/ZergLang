/* Copyright (C) 2020-2021 cmj <cmj@cmj.tw>. All right reserved. */
#include <stdio.h>
#include <stdlib.h>
#include <getopt.h>

#include "zerg.h"

int verbose = CRIT;

static void help(char *name) {
	fprintf(stderr, "%s (v%d.%d.%d) usage: %s [OPTIONS] FILE\n", PROJ_NAME, MAJOR, MINOR, MACRO, name);
	fprintf(stderr, "\n");
	fprintf(stderr, "option\n");
	fprintf(stderr, "  -h, --help     show this message\n");
	fprintf(stderr, "  -v, --verbose  verbose message\n");
	exit(-1);
}


int main(int argc, char *argv[]) {
	int opt, opt_idx = 0, ret = 1;
	const char opts[] = "vh";
	struct option long_options[] = {
		{"verbose"	, no_argument	, 0, 'v'},
		{"help"		, no_argument	, 0, 'h'},
	};

	while (-1 != (opt = getopt_long(argc, argv, opts, long_options, &opt_idx))) {
		switch (opt) {
			case 'h':
				help(argv[0]);
				break;
			case 'v':
				verbose ++;
				break;
			default:
				fprintf(stderr, "error: unknown option: '%c'\n", opt);
				help(argv[0]);
				goto END;
		}
	}

	_D(DEBUG, "finish parse the command, start parse source file");
	for (int idx = optind; idx < argc; ++idx) {
		_D(INFO, "process file '%s'", argv[idx]);
		if (0 > compile(argv[idx])) {
			_D(CRIT, "failed to compile '%s'", argv[idx]);
			goto END;	
		}
	}

	ret = 0;
END:
	return ret;
}
