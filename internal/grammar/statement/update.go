//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	builder "github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
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

func (st *Update) Size(b *builder.Builder) int {
	size := len(grammar.Symbols[grammar.SYM_UPDATE]) + len(st.table.Name) + len(grammar.Symbols[grammar.SYM_SET])
	ncols := len(st.columns)
	for _, c := range st.columns {
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
	if st.where != nil {
		size += st.where.Size(b)
	}
	return size
}

func (st *Update) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_UPDATE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(st.table.Name)
	b.Write(grammar.Symbols[grammar.SYM_SET])

	ncols := len(st.columns)
	for x, c := range st.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		b.WriteString(c.Name)
		b.Write(grammar.Symbols[grammar.SYM_EQUAL])
		b.AddInterpolationMarker(*curArg)
		args[*curArg] = st.values[x]
		*curArg++
		if x != (ncols - 1) {
			b.Write(grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}

	if st.where != nil {
		st.where.Scan(b, args, curArg)
	}
}

func (st *Update) AddWhere(e *expression.Expression) *Update {
	if st.where == nil {
		st.where = clause.NewWhere(e)
		return st
	}
	st.where.AddExpression(e)
	return st
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
