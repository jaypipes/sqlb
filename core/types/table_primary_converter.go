//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

import "github.com/jaypipes/sqlb/grammar"

// TablePrimaryConverter knows how to convert itself into a
// `*grammar.TablePrimary`
type TablePrimaryConverter interface {
	Named
	// TablePrimary returns the object as a `*grammar.TablePrimary`
	TablePrimary() *grammar.TablePrimary
}
