" Vim syntax file
" Language:    Zerg
" Maintainer:  cmj <cmj@cmj.tw>
" Last Change: 2026

if exists("b:current_syntax")
  finish
endif

" --- Keywords ---
syn keyword zergKeyword       pub mut const static this
syn keyword zergKeyword       impl nextgroup=zergClassName skipwhite
syn keyword zergFnDecl        fn nextgroup=zergFnName skipwhite
syn keyword zergConditional   if else match
syn keyword zergRepeat        for in
syn keyword zergStatement     return break continue del raise yield go import assert nop
syn keyword zergException     try expect finally with as
syn keyword zergStructure     class nextgroup=zergClassName skipwhite
syn keyword zergStructure     spec nextgroup=zergSpecName skipwhite
syn keyword zergStructure     enum nextgroup=zergEnumName skipwhite
syn keyword zergStructure     type
syn keyword zergOperator      and or xor not is

" --- Constants ---
syn keyword zergBoolean       true false
syn keyword zergKeyword       Ok Err
syn keyword zergNil           nil

" --- Built-in types ---
syn keyword zergType          int float bool string list map set chan object iter range Self

" --- Built-in functions ---
syn keyword zergBuiltin       print len input str int float

" --- Wildcard pattern ---
syn match   zergWildcard      /\<_\>/

" --- Comments ---
syn match   zergComment       /#.*$/ contains=zergTodo
syn keyword zergTodo          contained TODO FIXME XXX NOTE HACK

" --- Numbers ---
syn match   zergNumber        /\<0\>/
syn match   zergNumber        /\<[1-9][0-9_]*\>/
syn match   zergNumber        /\<0[xX][0-9a-fA-F][0-9a-fA-F_]*\>/
syn match   zergNumber        /\<0[oO][0-7][0-7_]*\>/
syn match   zergNumber        /\<0[bB][01][01_]*\>/
syn match   zergFloat         /\<[0-9][0-9_]*\.[0-9][0-9_]*\([eE][+-]\?[0-9][0-9_]*\)\?\>/
syn match   zergFloat         /\<[0-9][0-9_]*[eE][+-]\?[0-9][0-9_]*\>/

" --- Strings ---
syn region  zergString        start=/"/ skip=/\\./ end=/"/ contains=zergEscape,zergInterpolation
syn match   zergEscape        contained /\\[ntr\\"0{}]/
syn match   zergEscape        contained /\\x[0-9a-fA-F]\{2}/
syn match   zergEscape        contained /\\u{[0-9a-fA-F]\+}/
" Note: keepend means nested {} (e.g. map literals) will end interpolation early.
" This is a known Vim limitation; simple expressions highlight correctly.
syn region  zergInterpolation contained start=/{/ end=/}/ contains=TOP keepend
syn region  zergRawString     start=/\<r"/ end=/"/

" --- Operators ---
syn match   zergOperatorSym   /??/
syn match   zergOperatorSym   /?\./
syn match   zergOperatorSym   /?\[/
syn match   zergOperatorSym   /=>/
syn match   zergOperatorSym   /->/
syn match   zergOperatorSym   /:=/
syn match   zergOperatorSym   /==/
syn match   zergOperatorSym   /!=/
syn match   zergOperatorSym   /<=/
syn match   zergOperatorSym   />=/
syn match   zergOperatorSym   /<</
syn match   zergOperatorSym   />>/
syn match   zergOperatorSym   /\*\*/
syn match   zergOperatorSym   /\/\//
syn match   zergOperatorSym   /\.\.=/
syn match   zergOperatorSym   /\.\./
syn match   zergChanOp        /<-/
syn match   zergRefOp         /&\ze\w/
syn match   zergOperatorSym   /+=/
syn match   zergOperatorSym   /-=/
syn match   zergOperatorSym   /\*=/
syn match   zergOperatorSym   /\/=/
syn match   zergOperatorSym   /\/\/=/
syn match   zergOperatorSym   /%=/
syn match   zergOperatorSym   /\*\*=/
syn match   zergOperatorSym   /&=/
syn match   zergOperatorSym   /|=/
syn match   zergOperatorSym   /\^=/
syn match   zergOperatorSym   /<<=/
syn match   zergOperatorSym   />>=/
syn match   zergOperatorSym   /++/
syn match   zergOperatorSym   /--/

" --- Declaration names (chained via nextgroup) ---
syn match   zergFnName        /\w\+/ contained
syn match   zergClassName     /\w\+/ contained
syn match   zergSpecName      /\w\+/ contained
syn match   zergEnumName      /\w\+/ contained

" --- Highlight links ---
hi def link zergKeyword       Keyword
hi def link zergConditional   Conditional
hi def link zergRepeat        Repeat
hi def link zergStatement     Statement
hi def link zergException     Exception
hi def link zergStructure     Structure
hi def link zergOperator      Operator
hi def link zergOperatorSym   Operator
hi def link zergBoolean       Boolean
hi def link zergNil           Constant
hi def link zergType          Type
hi def link zergWildcard      Special
hi def link zergNumber        Number
hi def link zergFloat         Float
hi def link zergString        String
hi def link zergRawString     String
hi def link zergEscape        SpecialChar
hi def link zergInterpolation Special
hi def link zergComment       Comment
hi def link zergTodo          Todo
hi def link zergFnDecl        Keyword
hi def link zergFnName        Function
hi def link zergClassName     Function
hi def link zergSpecName      Function
hi def link zergEnumName      Function
hi def link zergChanOp        Special
hi def link zergRefOp         Special
hi def link zergBuiltin       Function

let b:current_syntax = "zerg"
