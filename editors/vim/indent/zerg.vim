" Vim indent file
" Language:    Zerg
" Maintainer:  cmj0121
" Last Change: 2026-03-25
" Filenames:   *.zg

if exists("b:did_indent")
  finish
endif
let b:did_indent = 1

setlocal indentexpr=GetZergIndent()
setlocal indentkeys=0{,0},!^F,o,O,e,0),0]
setlocal autoindent

if exists("*GetZergIndent")
  finish
endif

function! GetZergIndent()
  let lnum = prevnonblank(v:lnum - 1)
  if lnum == 0
    return 0
  endif

  let prev = getline(lnum)
  let curr = getline(v:lnum)
  let ind = indent(lnum)

  " increase indent after { or lines ending with operators
  if prev =~ '{\s*$'
    let ind += shiftwidth()
  endif

  " decrease indent for closing }
  if curr =~ '^\s*}'
    let ind -= shiftwidth()
  endif

  return ind
endfunction
