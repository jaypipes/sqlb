//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/core/grammar"

// DerivedColumnConverter knows how to convert itself into a
// `*grammar.DerivedColumn`
type DerivedColumnConverter interface {
	// DerivedColumn returns the object as a `*grammar.DerivedColumn`
	DerivedColumn() *grammar.DerivedColumn
}
