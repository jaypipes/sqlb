//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/scanner"
)

// DELETE FROM <table> WHERE <predicates>
type Delete struct {
	table *identifier.Table
	where *clause.Where
}

func (st *Delete) ArgCount() int {
	argc := 0
	if st.where != nil {
		argc += st.where.ArgCount()
	}
	return argc
}

func (st *Delete) Size(s *scanner.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_DELETE]) + len(st.table.Name)
	if st.where != nil {
		size += st.where.Size(s)
	}
	return size
}

func (st *Delete) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(st.table.Name)
	if st.where != nil {
		st.where.Scan(s, b, args, curArg)
	}
}

func (st *Delete) AddWhere(e *expression.Expression) *Delete {
	if st.where == nil {
		st.where = clause.NewWhere(e)
		return st
	}
	st.where.AddExpression(e)
	return st
}

// NewDelete returns a new DeleteStatement struct that scans into a DELETE SQL
// statement
func NewDelete(
	table *identifier.Table,
	where *clause.Where,
) *Delete {
	return &Delete{
		table: table,
		where: where,
	}
}
