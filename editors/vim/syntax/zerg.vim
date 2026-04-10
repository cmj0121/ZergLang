" Vim syntax file
" Language:    Zerg
" Maintainer:  cmj0121
" Last Change: 2026-03-25
" Filenames:   *.zg

if exists("b:current_syntax")
  finish
endif

" ── Keywords ───────────────────────────────────────────────────────────
syn keyword zergDecl        fn struct enum spec impl type import const pub mut defer as
syn keyword zergConditional if elif else match select
syn keyword zergRepeat      for break continue in
syn keyword zergControl     return nop rush raise
syn keyword zergException   try except finally
syn keyword zergOperator    and or not xor
syn keyword zergConstant    true false nil
syn keyword zergSpecial     this
syn keyword zergStatement   print

" Built-in types and specs
syn keyword zergType        int float bool str byte rune list map set tuple chan sync ptr Result Option
syn keyword zergType        Exception Printable Iterable Comparable Hashable Eq

" ── Numbers ────────────────────────────────────────────────────────────
syn match   zergNumber      "\<\d[0-9_]*\>"
syn match   zergNumber      "\<0x[0-9a-fA-F_]\+\>"
syn match   zergNumber      "\<0b[01_]\+\>"
syn match   zergNumber      "\<0o[0-7_]\+\>"
syn match   zergFloat       "\<\d[0-9_]*\.\d[0-9_]*\>"

" ── Strings ────────────────────────────────────────────────────────────
syn region  zergString      start=/"/ skip=/\\./ end=/"/ contains=zergEscape,zergInterpolation
syn region  zergRawString   start=/r"/ end=/"/
syn region  zergMultiString start=/"""/ end=/"""/ contains=zergEscape,zergInterpolation
syn match   zergEscape      contained "\\[ntr0\\\"'{}]"
syn region  zergInterpolation contained start=/{/ end=/}/ contains=TOP

" Rune literals
syn region  zergRune        start=/'/ end=/'/ contains=zergEscape

" ── Operators ──────────────────────────────────────────────────────────
syn match   zergOp          ":="
syn match   zergOp          "=="
syn match   zergOp          "!="
syn match   zergOp          "<="
syn match   zergOp          ">="
syn match   zergOp          "->"
syn match   zergOp          "<-"
syn match   zergOp          "\.\."
syn match   zergOp          "\.\.="
syn match   zergOp          "??"
syn match   zergOp          "?\."
syn match   zergOp          "<<"
syn match   zergOp          ">>"
syn match   zergOp          "+="
syn match   zergOp          "-="
syn match   zergOp          "\*="
syn match   zergOp          "/="
syn match   zergOp          "%="
syn match   zergOp          "&="
syn match   zergOp          "|="
syn match   zergOp          "\^="
syn match   zergOp          "<<="
syn match   zergOp          ">>="

" ── Definitions ────────────────────────────────────────────────────────
syn match   zergFuncDef     "fn\s\+\zs\w\+" contained
syn match   zergFuncStart   "fn\s\+\w\+" contains=zergDecl,zergFuncDef

syn match   zergTypeDef     "\<struct\s\+\zs\w\+"
syn match   zergTypeDef     "\<enum\s\+\zs\w\+"
syn match   zergTypeDef     "\<spec\s\+\zs\w\+"
syn match   zergTypeDef     "\<impl\s\+\zs\w\+"
syn match   zergTypeDef     "\<type\s\+\zs\w\+"

" ── Inline assembly ────────────────────────────────────────────────────
" ARM64 GNU as syntax inside asm { ... } block.
" # is ARM immediate prefix (not a Zerg comment), // is asm line comment.
" {expr} interpolates Zerg byte values (mut or immut).
syn region  zergAsmBlock    matchgroup=zergAsmKeyword start="\<asm\s*{" end="}" contains=zergAsmReg,zergAsmDirective,zergAsmNumber,zergAsmInterp,zergAsmComment,zergAsmLabel
syn keyword zergAsmReg      contained x0 x1 x2 x3 x4 x5 x6 x7 x8 x9 x10 x11 x12 x13 x14 x15 x16 x17 x18 x19 x20 x21 x22 x23 x24 x25 x26 x27 x28 x29 x30 xzr sp lr fp
syn keyword zergAsmReg      contained w0 w1 w2 w3 w4 w5 w6 w7 w8 w9 w10 w11 w12 w13 w14 w15 w16 w17 w18 w19 w20 w21 w22 w23 w24 w25 w26 w27 w28 w29 w30 wzr
syn match   zergAsmDirective contained "\<\(mov\|add\|sub\|ldr\|str\|svc\|bl\|b\|ret\|cmp\|cbz\|cbnz\|adr\|adrp\|and\|orr\|eor\|lsl\|lsr\|stp\|ldp\|nop\|mul\|neg\|mvn\)\>"
syn match   zergAsmNumber   contained "#\d\+"
syn match   zergAsmNumber   contained "#0x[0-9a-fA-F]\+"
syn region  zergAsmInterp   contained start=/\${/ end=/}/ contains=TOP
syn match   zergAsmComment  contained "//.*$"
syn match   zergAsmLabel    contained "^\s*\w\+:"

" ── Comments ───────────────────────────────────────────────────────────
" # starts a Zerg comment. ## starts a doc comment.
" Exclude #<digit> and #0x so ARM immediates (#1, #0x80) in asm blocks
" are not treated as comments.
syn match   zergDocComment  "##.*$"
syn match   zergComment     "#\(#\|[0-9]\|0x\)\@!.*$"

" ── Highlight groups ───────────────────────────────────────────────────
hi def link zergDecl        Keyword
hi def link zergConditional Conditional
hi def link zergRepeat      Repeat
hi def link zergControl     Statement
hi def link zergException   Exception
hi def link zergOperator    Operator
hi def link zergConstant    Constant
hi def link zergSpecial     Special
hi def link zergStatement   Statement
hi def link zergType        Type
hi def link zergNumber      Number
hi def link zergFloat       Float
hi def link zergString      String
hi def link zergRawString   String
hi def link zergMultiString String
hi def link zergRune        Character
hi def link zergEscape      SpecialChar
hi def link zergInterpolation Special
hi def link zergDocComment  SpecialComment
hi def link zergComment     Comment
hi def link zergOp          Operator
hi def link zergFuncDef     Function
hi def link zergTypeDef     TypeDef
hi def link zergAsmReg      Identifier
hi def link zergAsmDirective Statement
hi def link zergAsmNumber   Number
hi def link zergAsmInterp   Special
hi def link zergAsmComment  Comment
hi def link zergAsmLabel    Label
hi def link zergAsmKeyword  Keyword

let b:current_syntax = "zerg"
