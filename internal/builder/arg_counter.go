//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

// ArgCounter updates a supplied count of query arguments.
type ArgCounter interface {
	// ArgCount updates the supplied count of arguments that the ArgCounter
	// contains
	ArgCount(*int)
}
