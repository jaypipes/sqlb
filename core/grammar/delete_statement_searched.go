//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <delete statement: searched>    ::=   DELETE FROM <target table> [ WHERE <search condition> ]
//
// <target table>    ::=
//          <table name>
//      |     ONLY <left paren> <table name> <right paren>

// DeleteStatementSearched represents a DELETE FROM SQL statement
type DeleteStatementSearched struct {
	TableName string
	Where     *WhereClause
}

func (s *DeleteStatementSearched) ArgCount(count *int) {
	if s.Where != nil {
		s.Where.ArgCount(count)
	}
}
