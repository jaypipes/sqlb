//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Orderable is an Element that knows whether it is to be part of an ORDER BY
// clause
type Orderable interface {
	Element
	// IsAsc returns true if the element is to be sorted in ascending order
	IsAsc() bool
}
