" Zerg Syntax File
" Language:      zerg
" Maintainer:    cmj <cmj@cmj.tw>
" Last Change:   2024 Feb 29
"
if exists("b:current_syntax")
    finish
endif

let s:cpo_save = &cpo
set cpo&vim

let zerg_highlight=1

syn match   ZergOperator    display "\%(+\|-\|*\|/\|%\|<\|>\|&\|\^\||\|:=\)=\?"
syn match   ZergComment     "//.*$" contains=ZergToDo,@Spell
syn region  ZergCommentDoc  start="//!" end="$" contains=ZergTodo,@Spell
syn keyword ZergToDo        contained TODO FIXME XXX NOTE HACK
syn keyword ZergStatement   fn struct enum impl nextgroup=ZergName skipwhite
syn keyword ZergStatement   print return defer async for loop in asm if not nop break continue
syn match   ZergName        contained /\w\+/
syn region  ZergString      start=+"+ end=+"+ skip=+\\+ contains=ZergEscape,ZergStringVar
syn match   ZergStringVar   display contained "\$\w\+"
syn match   ZergEscape      display contained "\\[0-7]\{3}"
syn match   ZergEscape      display contained "\\x\x\{2}"
syn match   ZergEscape      display contained "\\u\x\{4}"
syn match   ZergEscape      display contained "\\U\x\{8}"
syn match   ZergEscape      display contained +\\[abfnrtv\\'"]+
syn match   ZergNumber      "\<\d\+\>"
syn keyword ZergType        i8 u8 i16 u16 i32 u32 i64 u64 i128 u128
syn keyword ZergType        size usize chan
syn keyword ZergType        mut str Self
syn keyword ZergBuiltIn     make len

" Zerg highlight syntax definition
hi def link ZergCommentDoc  SpecialComment
hi def link ZergComment     Comment
hi def link ZergToDo        Todo
hi def link ZergStatement   Statement
hi def link ZergLogic       Statement
hi def link ZergName        Function
hi def link ZergString      String
hi def link ZergStringVar   Special
hi def link ZergEscape      SpecialChar
hi def link ZergNumber      Number
hi def link ZergType        Type
hi def link ZergOperator    Operator
hi def link ZergBuiltIn     Function

let &cpo = s:cpo_save
unlet s:cpo_save
