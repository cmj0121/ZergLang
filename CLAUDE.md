# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Zerg is an open-source programming language: procedural-first, strongly typed,
null-safe, garbage-collected, with built-in concurrency.
The project follows "Dream-Driven Development" — features are driven by the
author's needs and vision.

## Build & Test

```bash
make              # default: installs pre-commit hooks, sets commit template
make test          # run tests
make build         # build binaries
make clean         # remove swap files
make upgrade       # update pre-commit hooks
```

The Makefile uses `SUBDIR` to recurse into subdirectories — add new source directories there as the project grows.

## Commit Convention

Uses conventional commits with the template in `.git-commit-template`:

```text
<type>(scope): <subject>
```

Types: `feat`, `docs`, `test`, `perf`, `build`, `style`, `refactor`

## Pre-commit Hooks

Configured in `.pre-commit-config.yaml`:

- YAML validation, trailing whitespace, EOF newline
- markdownlint (auto-fix)
- prettier (formatting)
- gitleaks (secret detection)

## Language Design Principles

- **Small and crisp** — minimal syntax
- **Procedural-first** — top-down control flow
- **Strongly typed** with null safety
- **Concurrent** — built-in concurrency primitives
- **Garbage-collected**
