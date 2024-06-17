//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type FunctionID int

const (
	FUNC_MAX FunctionID = iota
	FUNC_MIN
	FUNC_SUM
	FUNC_AVG
	FUNC_COUNT_STAR
	FUNC_COUNT_DISTINCT
	FUNC_CAST
	FUNC_CHAR_LENGTH
	FUNC_BIT_LENGTH
	FUNC_ASCII
	FUNC_REVERSE
	FUNC_CONCAT
	FUNC_CONCAT_WS
	FUNC_NOW
	FUNC_CURRENT_TIMESTAMP
	FUNC_CURRENT_TIME
	FUNC_CURRENT_DATE
	FUNC_EXTRACT
)

var (
	// A static table containing information used in constructing the
	// expression's SQL string during scan() calls
	funcScanTable = map[FunctionID]ScanInfo{
		FUNC_MAX: {
			SYM_MAX, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_MIN: {
			SYM_MIN, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_SUM: {
			SYM_SUM, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_AVG: {
			SYM_AVG, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_COUNT_STAR: {
			SYM_COUNT_STAR,
		},
		FUNC_COUNT_DISTINCT: {
			SYM_COUNT_DISTINCT, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CAST: {
			SYM_CAST, SYM_ELEMENT, SYM_AS, SYM_PLACEHOLDER, SYM_RPAREN,
		},
		FUNC_CHAR_LENGTH: {
			SYM_CHAR_LENGTH, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_BIT_LENGTH: {
			SYM_BIT_LENGTH, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_ASCII: {
			SYM_ASCII, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_REVERSE: {
			SYM_REVERSE, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CONCAT: {
			SYM_CONCAT, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_CONCAT_WS: {
			SYM_CONCAT_WS, SYM_ELEMENT, SYM_COMMA_WS, SYM_ELEMENT, SYM_RPAREN,
		},
		FUNC_NOW: {
			SYM_NOW,
		},
		FUNC_CURRENT_TIMESTAMP: {
			SYM_CURRENT_TIMESTAMP,
		},
		FUNC_CURRENT_TIME: {
			SYM_CURRENT_TIME,
		},
		FUNC_CURRENT_DATE: {
			SYM_CURRENT_DATE,
		},
		// This is the MySQL variant of EXTRACT, which follows the form
		// EXTRACT(field FROM source). PostgreSQL has a different format for
		// EXTRACT() which follows the following format:
		// EXTRACT(field FROM [interval|timestamp] source)
		FUNC_EXTRACT: {
			SYM_EXTRACT, SYM_PLACEHOLDER, SYM_SPACE, SYM_FROM, SYM_ELEMENT, SYM_RPAREN,
		},
	}
)

func FunctionScanTable(id FunctionID) ScanInfo {
	return funcScanTable[id]
}
