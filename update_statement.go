//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// UPDATE <table> SET <column_value_list>[ WHERE <predicates>]

type updateStatement struct {
	table   *Table
	columns []*Column
	values  []interface{}
	where   *whereClause
}

func (s *updateStatement) ArgCount() int {
	argc := len(s.values)
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	return argc
}

func (s *updateStatement) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_UPDATE]) + len(s.table.name) + len(grammar.Symbols[grammar.SYM_SET])
	ncols := len(s.columns)
	for _, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		size += len(c.name)
	}
	// NOTE(jaypipes): We do not include the length of interpolation markers,
	// since that differs based on the SQL dialect
	size += len(grammar.Symbols[grammar.SYM_EQUAL]) * ncols
	// Two comma-delimited lists of same number of elements (columns and
	// values)
	size += 2 * (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
	if s.where != nil {
		size += s.where.Size(scanner)
	}
	return size
}

func (s *updateStatement) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_UPDATE])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.name)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_SET])

	ncols := len(s.columns)
	for x, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		bw += copy(b[bw:], c.name)
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_EQUAL])
		bw += scanInterpolationMarker(scanner.Dialect(), b[bw:], *curArg)
		args[*curArg] = s.values[x]
		*curArg++
		if x != (ncols - 1) {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}

	if s.where != nil {
		bw += s.where.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func (s *updateStatement) addWhere(e *Expression) *updateStatement {
	if s.where == nil {
		s.where = &whereClause{filters: make([]*Expression, 0)}
	}
	s.where.filters = append(s.where.filters, e)
	return s
}
