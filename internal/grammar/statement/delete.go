//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
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

func (st *Delete) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.Write(grammar.Symbols[grammar.SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	b.WriteString(st.table.Name)
	if st.where != nil {
		b.WriteString(st.where.String(opts, qargs, curarg))
	}
	return b.String()
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
