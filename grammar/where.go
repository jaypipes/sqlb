//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// WHERE <search condition>

// Where represents the SQL WHERE clause
type Where struct {
	SearchConditions []*BooleanValueExpression
}
