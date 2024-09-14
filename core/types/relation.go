//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

// Relation is a thing that can be included in a SELECT statement's FROM clause
// and produce zero or more Projection things that can be referenced in the
// SELECT statement.
type Relation interface {
	Aliased
	QuerySpecificationConverter
	TablePrimaryConverter
	TableReferenceConverter
	// Projections returns a slice of Projection things referenced by the
	// Selectable. The slice should be sorted by the Projection's name.
	Projections() []Projection
	// C returns the Projection with the given or nil if no such Projection
	// was found
	C(name string) Projection
}
