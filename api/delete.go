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

// DeleteStatement represents a DELETE FROM SQL statement
type DeleteStatement struct {
	ds *grammar.DeleteStatementSearched
}

// Query returns the underlying SQL grammar element which is passed to the
// builder
func (s *DeleteStatement) Query() *grammar.DeleteStatementSearched {
	return s.ds
}

// Where returns the DeleteStatementSearched adapted with a supplied search
// condition
func (s *DeleteStatement) Where(
	expr interface{},
) *DeleteStatement {
	if s.ds == nil {
		panic("cannot call Where on a nil DeleteStatement")
	}
	bve := BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	if s.ds.WhereClause != nil {
		s.ds.WhereClause.SearchCondition = And(&s.ds.WhereClause.SearchCondition, bve)
	} else {
		s.ds.WhereClause = &grammar.WhereClause{
			SearchCondition: *bve,
		}
	}
	return s
}

// Delete returns a Queryable that produces a DELETE SQL statement for a given
// table
func Delete(
	t *Table,
) (*DeleteStatement, error) {
	if t == nil {
		return nil, TableRequired
	}

	return t.Delete(), nil
}
