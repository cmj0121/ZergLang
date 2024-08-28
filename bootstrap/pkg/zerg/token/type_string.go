// Code generated by "stringer -type=Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[EOF-1]
	_ = x[EOL-2]
	_ = x[LBracket-3]
	_ = x[RBracket-4]
	_ = x[LParen-5]
	_ = x[RParen-6]
	_ = x[LBrace-7]
	_ = x[RBrace-8]
	_ = x[Add-9]
	_ = x[Sub-10]
	_ = x[Mul-11]
	_ = x[Div-12]
	_ = x[Mod-13]
	_ = x[Arrow-14]
	_ = x[Name-15]
	_ = x[String-16]
	_ = x[Int-17]
	_ = x[Fn-18]
	_ = x[Str-19]
	_ = x[Print-20]
}

const _Type_name = "UnknownEOFEOLLBracketRBracketLParenRParenLBraceRBraceAddSubMulDivModArrowNameStringIntFnStrPrint"

var _Type_index = [...]uint8{0, 7, 10, 13, 21, 29, 35, 41, 47, 53, 56, 59, 62, 65, 68, 73, 77, 83, 86, 88, 91, 96}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
