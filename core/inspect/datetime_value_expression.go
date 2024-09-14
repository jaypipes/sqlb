//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/grammar"

// ReferredFromDatetimeValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied DatetimeValueExpression.
func ReferredFromDatetimeValueExpression(
	cve *grammar.DatetimeValueExpression,
) []string {
	return []string{}
}
