//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// WHERE <search condition>

// WhereClause represents the SQL WHERE clause
type WhereClause struct {
	SearchCondition BooleanValueExpression
}
