//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"github.com/jaypipes/sqlb/grammar"
)

// Projectable things can be included in a SELECT statement and allow
// inspection of the tables or derived tables (subqueries) that they reference.
type Projectable interface {
	// RefersTo returns a slice of tables or derived tables that are referenced
	// by the Projectable
	RefersTo() []interface{}
	// DerivedColumn returns the `*grammar.DerivedColumn` element representing
	// the Projectable
	DerivedColumn() *grammar.DerivedColumn
}
