// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package sqlb

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/errors"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

type UpdateQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *statement.Update
	scanner types.Scanner
}

func (q *UpdateQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *UpdateQuery) Error() error {
	return q.e
}

func (q *UpdateQuery) Scan(s types.Scanner, b *strings.Builder, qargs []interface{}, idx *int) {
	q.stmt.Scan(s, b, qargs, idx)
}

func (q *UpdateQuery) ArgCount() int {
	return q.stmt.ArgCount()
}

func (q *UpdateQuery) Size(s types.Scanner) int {
	return q.stmt.Size(s)
}

func (q *UpdateQuery) Where(e *expression.Expression) *UpdateQuery {
	q.stmt.AddWhere(e)
	return q
}

// Given a table and a map of column name to value for that column to update,
// returns an UpdateQuery that will produce an UPDATE SQL statement
func Update(t *identifier.Table, values map[string]interface{}) *UpdateQuery {
	if t == nil {
		return &UpdateQuery{e: errors.InvalidUpdateNoTarget}
	}
	if len(values) == 0 {
		return &UpdateQuery{e: errors.InvalidUpdateNoValues}
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*identifier.Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return &UpdateQuery{e: errors.InvalidUpdateUnknownColumn}
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	return &UpdateQuery{
		scanner: scanner.New(t.Schema().Dialect),
		stmt:    statement.NewUpdate(t, cols, vals, nil),
	}
}

//func (t *Table) Update(values map[string]interface{}) *UpdateQuery {
//	return Update(t, values)
//}
