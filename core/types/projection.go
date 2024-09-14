//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import (
	"github.com/jaypipes/sqlb/grammar"
)

// Projection is a thing that can be included in a SELECT statement and allow
// inspection of the Relations (tables or derived tables/subqueries) that they
// reference.
type Projection interface {
	Named
	// References returns the table or derived table that is referenced
	// by the Projection, or nil if the Projection references no table
	References() Relation
	// DerivedColumn returns the `*grammar.DerivedColumn` element representing
	// the Projection
	DerivedColumn() *grammar.DerivedColumn
	// As returns a copy of the Projection aliased with a new name
	As(alias string) Projection
	// Asc returns a SortSpecification indicating the Column should used in an
	// ORDER BY clause in ASCENDING sort order
	Asc() grammar.SortSpecification
	// Desc returns a SortSpecification indicating the Column should used in an
	// ORDER BY clause in DESCENDING sort order
	Desc() grammar.SortSpecification
}
