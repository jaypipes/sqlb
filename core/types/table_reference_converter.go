//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/core/grammar"

// TableReferenceConverter knows how to convert itself into a
// `*grammar.TableReference`
type TableReferenceConverter interface {
	// TableReference returns the object as a `*grammar.TableReference`
	TableReference() *grammar.TableReference
}
