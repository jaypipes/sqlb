//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Selection is something that produces rows. A table, table definition,
// view, subselect, etc.
type Selection interface {
	Element
	Projections() []Projection
}
