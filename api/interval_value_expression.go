//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import "github.com/jaypipes/sqlb/grammar"

// ReferredFromIntervalValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied IntervalValueExpression.
func ReferredFromIntervalValueExpression(
	cve *grammar.IntervalValueExpression,
) []string {
	return []string{}
}
