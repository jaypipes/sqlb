//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

// Scannable is a thing that knows how to describe itself to a Scanner for
// construction in a SQL query string and query argument list.
type Scannable interface {
	// Scan takes two slices and a pointer to an int. The first slice is a
	// slice of bytes that the implementation should copy its string
	// representation to and the other slice is a slice of interface{} values
	// that the element should add its arguments to. The pointer to an int is
	// the index of the current argument to be processed. The method returns a
	// single int, the number of bytes written to the buffer.
	Scan(Scanner, []byte, []interface{}, *int) int
}

// Element adds a Size and ArgCount method to the Scannable interface
type Element interface {
	Scannable
	// Size returns the number of bytes that the scannable element would
	// consume as a SQL string
	Size(Scanner) int
	// ArgCount returns the number of interface{} arguments that the element
	// will add to the slice of interface{} arguments passed to Scan()
	ArgCount() int
}

// Sortable is an Element that knows whether it is to be part of an ORDER BY
// clause
type Sortable interface {
	Element
	// IsAsc returns true if the element is to be sorted in ascending order
	IsAsc() bool
}

// Projection is something that produces a scalar value. A column, column
// definition, function, etc. When appearing in the SELECT clause's projection
// list, the projection will output itself using the "AS alias" extended
// notation. When outputting in GROUP BY, ORDER BY or ON clauses, the
// projection will not include the alias extension
type Projection interface {
	Element
	From() Selection
	// As returns the projection aliased as another name
	As(alias string) Projection
	// disables the outputting of the "AS alias" extended output. Returns a
	// function that resets the outputting of the "AS alias" extended output
	DisableAliasScan() func()
	// Asc returns an Element that describes a sort on the projection in
	// descending order
	Asc() Sortable
	// Desc returns an Element that describes a sort on the projection in
	// descending order
	Desc() Sortable
}

// Selection is something that produces rows. A table, table definition,
// view, subselect, etc.
type Selection interface {
	Element
	Projections() []Projection
}
