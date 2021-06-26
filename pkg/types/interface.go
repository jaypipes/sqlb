//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

type Scannable interface {
	// scan takes two slices and a pointer to an int. The first slice is a
	// slice of bytes that the implementation should copy its string
	// representation to and the other slice is a slice of interface{} values
	// that the element should add its arguments to. The pointer to an int is
	// the index of the current argument to be processed. The method returns a
	// single int, the number of bytes written to the buffer.
	Scan(Scanner, []byte, []interface{}, *int) int
}

type Element interface {
	// Returns the number of bytes that the scannable element would consume as
	// a SQL string
	Size(Scanner) int
	// Returns the number of interface{} arguments that the element will add to
	// the slice of interface{} arguments passed to Scan()
	ArgCount() int
	// scan takes two slices and a pointer to an int. The first slice is a slice of bytes that the
	// implementation should copy its string representation to and the other slice is a slice of interface{} values that the element should add its
	// arguments to. The pointer to an int is the index of the current argument to be processed. The method returns a single int, the number of bytes written to the buffer.
	Scan(Scanner, []byte, []interface{}, *int) int
}

// A projection is something that produces a scalar value. A column, column
// definition, function, etc. When appearing in the SELECT clause's projection
// list, the projection will output itself using the "AS alias" extended
// notation. When outputting in GROUP BY, ORDER BY or ON clauses, the
// projection will not include the alias extension
type Projection interface {
	From() Selection
	// projections must also implement element
	Size(Scanner) int
	ArgCount() int
	Scan(Scanner, []byte, []interface{}, *int) int
	// disables the outputting of the "AS alias" extended output. Returns a
	// function that resets the outputting of the "AS alias" extended output
	DisableAliasScan() func()
}

// A selection is something that produces rows. A table, table definition,
// view, subselect, etc.
type Selection interface {
	Projections() []Projection
	// selections must also implement element
	Size(Scanner) int
	ArgCount() int
	Scan(Scanner, []byte, []interface{}, *int) int
}

// A Query is a placeholder for something that can be asked for the SQL string
// representation of the underlying query clauses
type Query interface {
	IsValid() bool
	Error() error
	String() string
	StringArgs() (string, []interface{})
}
