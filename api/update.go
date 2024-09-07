//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"

	"github.com/jaypipes/sqlb/grammar"
)

// UpdateStatement represents an UPDATE SQL statement
type UpdateStatement struct {
	us *grammar.UpdateStatementSearched
}

// Query returns the underlying SQL grammar element which is passed to the
// builder
func (s *UpdateStatement) Query() *grammar.UpdateStatementSearched {
	return s.us
}

// Where returns the UpdateStatementSearched adapted with a supplied search
// condition
func (s *UpdateStatement) Where(
	expr interface{},
) *UpdateStatement {
	if s.us == nil {
		panic("cannot call Where on a nil UpdateStatement")
	}
	bve := BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	if s.us.WhereClause != nil {
		s.us.WhereClause.SearchCondition = And(&s.us.WhereClause.SearchCondition, bve)
	} else {
		s.us.WhereClause = &grammar.WhereClause{
			SearchCondition: *bve,
		}
	}
	return s
}

// Update returns a struct that will produce an UPDATE SQL statement for a
// given table and map of column name to value for that column to update.
func Update(
	t *Table,
	values map[string]interface{},
) (*UpdateStatement, error) {
	if len(values) == 0 {
		return nil, NoValues
	}
	if t == nil {
		return nil, TableRequired
	}
	return t.Update(values)
}
