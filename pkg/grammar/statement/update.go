//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	pkgscanner "github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

// UPDATE <table> SET <column_value_list>[ WHERE <predicates>]
type Update struct {
	table   *identifier.Table
	columns []*identifier.Column
	values  []interface{}
	where   *clause.Where
}

func (s *Update) ArgCount() int {
	argc := len(s.values)
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	return argc
}

func (s *Update) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_UPDATE]) + len(s.table.Name) + len(grammar.Symbols[grammar.SYM_SET])
	ncols := len(s.columns)
	for _, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		size += len(c.Name)
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

func (s *Update) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_UPDATE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(s.table.Name)
	b.Write(grammar.Symbols[grammar.SYM_SET])

	ncols := len(s.columns)
	for x, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		b.WriteString(c.Name)
		b.Write(grammar.Symbols[grammar.SYM_EQUAL])
		pkgscanner.ScanInterpolationMarker(scanner.Dialect(), b, *curArg)
		args[*curArg] = s.values[x]
		*curArg++
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}

	if s.where != nil {
		s.where.Scan(scanner, b, args, curArg)
	}
}

func (s *Update) AddWhere(e *expression.Expression) *Update {
	if s.where == nil {
		s.where = clause.NewWhere(e)
		return s
	}
	s.where.AddExpression(e)
	return s
}

// NewUpdate returns a new UpdateStatement struct that scans into an UPDATE SQL
// statement
func NewUpdate(
	table *identifier.Table,
	columns []*identifier.Column,
	values []interface{},
	where *clause.Where,
) *Update {
	return &Update{
		table:   table,
		columns: columns,
		values:  values,
		where:   where,
	}
}
