# Zerg Vim Syntax Plugin

Vim syntax highlighting and filetype detection for the Zerg programming language (`*.zg` files).

## Install

```sh
make install-vim
```

This copies the syntax and ftdetect files to `~/.vim/`. To use a different prefix:

```sh
make install-vim VIM_PREFIX=~/.config/nvim
```

## Uninstall

```sh
make uninstall-vim
```

## Files

- `syntax/zerg.vim` -- syntax highlighting rules
- `ftdetect/zerg.vim` -- auto-detection for `*.zg` files
