//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

// Element describes a single component of a SQL string
type Element interface {
	// String returns the element as a SQL string, populating the supplied
	// arguments with any values from the element.
	String(opts Options, args []interface{}, curarg *int) string
	// ArgCount returns the number of interface{} arguments that the element
	// will add to the slice of interface{} arguments passed to String()
	ArgCount() int
}
