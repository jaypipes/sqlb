//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/types"
)

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []types.Element {
	els := make([]types.Element, len(vars))
	for x, v := range vars {
		switch v.(type) {
		case types.Element:
			els[x] = v.(types.Element)
		default:
			els[x] = ast.NewValue(nil, v)
		}
	}
	return els
}

// Given a variable number of interface{} variables, returns a List containing
// Value structs for the variables
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toValueList(vars ...interface{}) *List {
	els := make([]types.Element, len(vars))
	for x, v := range vars {
		els[x] = ast.NewValue(nil, v)
	}
	return &List{elements: els}
}
