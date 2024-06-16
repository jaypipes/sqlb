//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package scanner

import (
	"strconv"
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/types"
)

// InterpolationLength returns the total length of the characters representing
// interpolation markers for query parameters. Different SQL dialects use
// different character sequences for marking query parameters during query
// preparation. For instance, MySQL and SQLite use the ? character. PostgreSQL
// uses a numbered $N schema with N starting at 1, SQL Server uses a :N scheme,
// etc.
func InterpolationLength(dialect types.Dialect, argc int) int {
	if dialect == types.DialectPostgreSQL {
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

// ScanInterpolationMarker adds an interpolation marker of the specified
// dialect and position into the supplied bytestream
func ScanInterpolationMarker(dialect types.Dialect, b *strings.Builder, position int) {
	if dialect == types.DialectPostgreSQL {
		b.Write(grammar.Symbols[grammar.SYM_DOLLAR])
		b.WriteString(strconv.Itoa(position + 1))
	} else {
		b.Write(grammar.Symbols[grammar.SYM_QUEST_MARK])
	}
}
