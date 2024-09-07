//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"github.com/jaypipes/sqlb/grammar"
)

// Insert returns a struct that will produce an INSERT SQL statement for a
// given table and map of column name to value for that column to insert.
func Insert(
	t *Table,
	values map[string]interface{},
) (*grammar.InsertStatement, error) {
	if len(values) == 0 {
		return nil, NoValues
	}
	if t == nil {
		return nil, TableRequired
	}
	return t.Insert(values)
}
