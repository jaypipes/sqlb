//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

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
	Asc() Orderable
	// Desc returns an Element that describes a sort on the projection in
	// descending order
	Desc() Orderable
}
