package sqlb

import "strconv"

type Dialect int

const (
	DIALECT_UNKNOWN = iota
	DIALECT_MYSQL
	DIALECT_POSTGRESQL
)

// Returns the total length of the characters representing interpolation
// markers for query parameters. Different SQL dialects use different character
// sequences for marking query parameters during query preparation. For
// instance, MySQL and SQLite use the ? character. PostgreSQL uses a numbered
// $N schema with N starting at 1, SQL Server uses a :N scheme, etc.
func interpolationLength(dialect Dialect, argc int) int {
	if dialect == DIALECT_POSTGRESQL {
		// $ character for each interpolated parameter plus ones digit of
		// number
		size := 2 * argc
		if argc > 9 {
			// tens digit
			size += argc - 9
		}
		if argc > 99 {
			// hundreds digit
			size += argc - 99
		}
		if argc > 999 {
			// thousands digit
			size += argc - 999
		}
		return size
	}
	return argc // Single question mark used as interpolation marker
}

func scanInterpolationMarker(dialect Dialect, b []byte, position int) int {
	if dialect == DIALECT_POSTGRESQL {
		bw := copy(b, Symbols[SYM_DOLLAR])
		bw += copy(b[bw:], []byte(strconv.Itoa(position+1)))
		return bw
	}
	return copy(b, Symbols[SYM_QUEST_MARK])
}
