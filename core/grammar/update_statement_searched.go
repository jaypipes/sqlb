//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <update statement: searched>    ::=   UPDATE <target table> SET <set clause list> [ WHERE <search condition> ]

// UpdateStatementSearched represents an UPDATE SQL statement
type UpdateStatementSearched struct {
	TableName string
	Columns   []string
	Values    []interface{}
	Where     *WhereClause
}

func (s *UpdateStatementSearched) ArgCount(count *int) {
	*count += len(s.Values)
	if s.Where != nil {
		s.Where.ArgCount(count)
	}
}
