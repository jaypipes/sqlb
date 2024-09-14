//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package types

// Named things have a Name() method that returns a string name
type Named interface {
	// Name returns the name of the Named thing. The name may be qualified
	// (i.e. have one or more periods in it)
	Name() string
}
