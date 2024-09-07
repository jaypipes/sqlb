//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <having clause>    ::=   HAVING <search condition>

// HavingClause represents the SQL HAVING clause
type HavingClause struct {
	SearchCondition BooleanValueExpression
}
