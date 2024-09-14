//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

// Aliased things can have an alias
type Aliased interface {
	Named
	// Alias returns the thing's alias, or an empty string if not aliased
	Alias() string
	// AliasOrName returns the thing's alias or its name if not aliased
	AliasOrName() string
}
