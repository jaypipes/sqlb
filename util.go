//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

// Given a slice of interface{} variables, returns a slice of element members.
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toElements(vars ...interface{}) []element {
	els := make([]element, len(vars))
	for x, v := range vars {
		switch v.(type) {
		case element:
			els[x] = v.(element)
		default:
			els[x] = &value{val: v}
		}
	}
	return els
}

// Given a variable number of interface{} variables, returns a List containing
// Value structs for the variables
// If any of the interface{} variables are *not* of type element already, we
// construct a Value{} for the variable.
func toValueList(vars ...interface{}) *List {
	els := make([]element, len(vars))
	for x, v := range vars {
		els[x] = &value{val: v}
	}
	return &List{elements: els}
}
