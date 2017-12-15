package sqlb

import (
	"errors"
)

var (
	ERR_UPDATE_NO_TARGET      = errors.New("No target table supplied.")
	ERR_UPDATE_NO_VALUES      = errors.New("No values supplied.")
	ERR_UPDATE_UNKNOWN_COLUMN = errors.New("Received an unknown column.")
)

type UpdateQuery struct {
	e       error
	b       []byte
	args    []interface{}
	stmt    *updateStatement
	scanner *sqlScanner
}

func (q *UpdateQuery) IsValid() bool {
	return q.e == nil && q.stmt != nil
}

func (q *UpdateQuery) Error() error {
	return q.e
}

func (q *UpdateQuery) String() string {
	size := q.stmt.size()
	argc := q.stmt.argCount()
	size += q.scanner.interpolationLength(argc)
	if len(q.args) != argc {
		q.args = make([]interface{}, argc)
	}
	if len(q.b) != size {
		q.b = make([]byte, size)
	}
	q.scanner.scan(q.b, q.args, q.stmt)
	return string(q.b)
}

func (q *UpdateQuery) StringArgs() (string, []interface{}) {
	size := q.stmt.size()
	argc := q.stmt.argCount()
	size += q.scanner.interpolationLength(argc)
	if len(q.args) != argc {
		q.args = make([]interface{}, argc)
	}
	if len(q.b) != size {
		q.b = make([]byte, size)
	}
	q.scanner.scan(q.b, q.args, q.stmt)
	return string(q.b), q.args
}

func (q *UpdateQuery) Where(e *Expression) *UpdateQuery {
	q.stmt.addWhere(e)
	return q
}

// Given a table and a map of column name to value for that column to update,
// returns an UpdateQuery that will produce an UPDATE SQL statement
func Update(t *Table, values map[string]interface{}) *UpdateQuery {
	if t == nil {
		return &UpdateQuery{e: ERR_UPDATE_NO_TARGET}
	}
	if len(values) == 0 {
		return &UpdateQuery{e: ERR_UPDATE_NO_VALUES}
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]*Column, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return &UpdateQuery{e: ERR_UPDATE_UNKNOWN_COLUMN}
		}
		cols[x] = c
		vals[x] = v
		x++
	}

	scanner := &sqlScanner{
		dialect: t.meta.dialect,
		format:  defaultFormatOptions,
	}
	stmt := &updateStatement{
		table:   t,
		columns: cols,
		values:  vals,
	}
	return &UpdateQuery{
		stmt:    stmt,
		scanner: scanner,
	}
}

func (t *Table) Update(values map[string]interface{}) *UpdateQuery {
	return Update(t, values)
}
