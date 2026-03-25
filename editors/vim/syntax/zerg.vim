" Vim syntax file
" Language:    Zerg
" Maintainer:  cmj0121
" Last Change: 2026-03-25
" Filenames:   *.zg

if exists("b:current_syntax")
  finish
endif

" Keywords — declarations
syn keyword zergDecl        fn struct enum spec impl type import const pub mut defer
" Keywords — control flow
syn keyword zergConditional if elif else match
syn keyword zergRepeat      for break continue in
syn keyword zergControl     return nop
" Keywords — operators (word-form)
syn keyword zergOperator    and or not xor
" Keywords — special values
syn keyword zergConstant    true false nil
" Keywords — special identifiers
syn keyword zergSpecial     this
" Keywords — statements
syn keyword zergStatement   print

" Built-in types
syn keyword zergType        int float bool str byte rune

" Numbers — decimal, hex, binary, octal with digit separators
syn match   zergNumber      "\<\d[0-9_]*\>"
syn match   zergNumber      "\<0x[0-9a-fA-F_]\+\>"
syn match   zergNumber      "\<0b[01_]\+\>"
syn match   zergNumber      "\<0o[0-7_]\+\>"
syn match   zergFloat       "\<\d[0-9_]*\.\d[0-9_]*\>"

" Strings — regular, raw, multi-line
syn region  zergString      start=/"/ skip=/\\./ end=/"/ contains=zergEscape,zergInterpolation
syn region  zergRawString   start=/r"/ end=/"/
syn region  zergMultiString start=/"""/ end=/"""/ contains=zergEscape,zergInterpolation
syn match   zergEscape      contained "\\[ntr0\\\"'{}]"
syn region  zergInterpolation contained start=/{/ end=/}/ contains=TOP

" Rune literals
syn region  zergRune        start=/'/ end=/'/ contains=zergEscape

" Comments
syn match   zergDocComment  "##.*$"
syn match   zergComment     "#[^#].*$"
syn match   zergComment     "^#$"

" Operators
syn match   zergOp          ":="
syn match   zergOp          "=="
syn match   zergOp          "!="
syn match   zergOp          "<="
syn match   zergOp          ">="
syn match   zergOp          "->"
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

" Function definitions (highlight name after fn)
syn match   zergFuncDef     "fn\s\+\zs\w\+" contained
syn match   zergFuncStart   "fn\s\+\w\+" contains=zergDecl,zergFuncDef

" Struct/enum/spec/impl/type names
syn match   zergTypeDef     "\<struct\s\+\zs\w\+"
syn match   zergTypeDef     "\<enum\s\+\zs\w\+"
syn match   zergTypeDef     "\<spec\s\+\zs\w\+"
syn match   zergTypeDef     "\<impl\s\+\zs\w\+"
syn match   zergTypeDef     "\<type\s\+\zs\w\+"

" Highlight groups
hi def link zergDecl        Keyword
hi def link zergConditional Conditional
hi def link zergRepeat      Repeat
hi def link zergControl     Statement
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

let b:current_syntax = "zerg"
