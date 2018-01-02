//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"errors"
)

var (
	ERR_DELETE_NO_TARGET = errors.New("No target table supplied.")
)

type DeleteQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *deleteStatement
	scanner *sqlScanner
}

func (q *DeleteQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *DeleteQuery) Error() error {
	return q.e
}

func (q *DeleteQuery) String() string {
	sizes := q.scanner.size(q.stmt)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.scan(q.b, q.args, q.stmt)
	return string(q.b)
}

func (q *DeleteQuery) StringArgs() (string, []interface{}) {
	sizes := q.scanner.size(q.stmt)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.scan(q.b, q.args, q.stmt)
	return string(q.b), q.args
}

func (q *DeleteQuery) Where(e *Expression) *DeleteQuery {
	q.stmt.addWhere(e)
	return q
}

// Given a table and a map of column name to value for that column to insert,
// returns an DeleteQuery that will produce an INSERT SQL statement
func Delete(t *Table) *DeleteQuery {
	if t == nil {
		return &DeleteQuery{e: ERR_DELETE_NO_TARGET}
	}

	scanner := &sqlScanner{
		dialect: t.meta.dialect,
		format:  defaultFormatOptions,
	}
	stmt := &deleteStatement{
		table: t,
	}
	return &DeleteQuery{
		stmt:    stmt,
		scanner: scanner,
	}
}

func (t *Table) Delete() *DeleteQuery {
	return Delete(t)
}
