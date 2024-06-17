//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type ExpressionType int

const (
	EXP_EQUAL ExpressionType = iota
	EXP_NEQUAL
	EXP_AND
	EXP_OR
	EXP_IN
	EXP_BETWEEN
	EXP_IS_NULL
	EXP_IS_NOT_NULL
	EXP_GREATER
	EXP_GREATER_EQUAL
	EXP_LESS
	EXP_LESS_EQUAL
)

var (
	// A static table containing information used in constructing the
	// expression's SQL string during scan() calls
	exprScanTable = map[ExpressionType]ScanInfo{
		EXP_EQUAL: {
			SYM_ELEMENT, SYM_EQUAL, SYM_ELEMENT,
		},
		EXP_NEQUAL: {
			SYM_ELEMENT, SYM_NEQUAL, SYM_ELEMENT,
		},
		EXP_AND: {
			SYM_LPAREN, SYM_ELEMENT, SYM_AND, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_OR: {
			SYM_LPAREN, SYM_ELEMENT, SYM_OR, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_IN: {
			SYM_ELEMENT, SYM_IN, SYM_ELEMENT, SYM_RPAREN,
		},
		EXP_BETWEEN: {
			SYM_ELEMENT, SYM_BETWEEN, SYM_ELEMENT, SYM_AND, SYM_ELEMENT,
		},
		EXP_IS_NULL: {
			SYM_ELEMENT, SYM_IS_NULL,
		},
		EXP_IS_NOT_NULL: {
			SYM_ELEMENT, SYM_IS_NOT_NULL,
		},
		EXP_GREATER: {
			SYM_ELEMENT, SYM_GREATER, SYM_ELEMENT,
		},
		EXP_GREATER_EQUAL: {
			SYM_ELEMENT, SYM_GREATER_EQUAL, SYM_ELEMENT,
		},
		EXP_LESS: {
			SYM_ELEMENT, SYM_LESS, SYM_ELEMENT,
		},
		EXP_LESS_EQUAL: {
			SYM_ELEMENT, SYM_LESS_EQUAL, SYM_ELEMENT,
		},
	}
)

func ExpressionScanTable(et ExpressionType) ScanInfo {
	return exprScanTable[et]
}
