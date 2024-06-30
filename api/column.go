//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Column describes a column in a Table
type Column struct {
	// Table is a pointer to the Table housing this Column
	Table *Table
	// Name is the name of the Column in the Table
	Name string
}
