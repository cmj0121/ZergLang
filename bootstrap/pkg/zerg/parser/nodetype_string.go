// Code generated by "stringer -type=NodeType"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Root-0]
	_ = x[Fn-1]
	_ = x[Args-2]
	_ = x[Scope-3]
	_ = x[Type-4]
	_ = x[ReturnStmt-5]
	_ = x[PrintStmt-6]
	_ = x[Expression-7]
}

const _NodeType_name = "RootFnArgsScopeTypeReturnStmtPrintStmtExpression"

var _NodeType_index = [...]uint8{0, 4, 6, 10, 15, 19, 29, 38, 48}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
