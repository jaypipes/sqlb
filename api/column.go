//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Column describes a column in a Table
type Column struct {
	Alias string
	Name  string
	Table *Table
}
