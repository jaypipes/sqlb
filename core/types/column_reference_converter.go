//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/core/grammar"

// ColumnReferenceConverter knows how to convert itself into a
// `*grammar.ColumnReference`
type ColumnReferenceConverter interface {
	// ColumnReference returns the object as a `*grammar.ColumnReference`
	ColumnReference() *grammar.ColumnReference
}
