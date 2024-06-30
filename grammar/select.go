//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <query specification>    ::=   SELECT [ <set quantifier> ] <select list> <table expression>
//
// <select list>    ::=   <asterisk> | <select sublist> [ { <comma> <select sublist> }... ]
//
// <select sublist>    ::=   <derived column> | <qualified asterisk>
//
// <qualified asterisk>    ::=
//          <asterisked identifier chain> <period> <asterisk>
//      |     <all fields reference>
//
// <asterisked identifier chain>    ::=   <asterisked identifier> [ { <period> <asterisked identifier> }... ]
//
// <asterisked identifier>    ::=   <identifier>
//
// <derived column>    ::=   <value expression> [ <as clause> ]
//
// <as clause>    ::=   [ AS ] <column name>
//
// <all fields reference>    ::=   <value expression primary> <period> <asterisk> [ AS <left paren> <all fields column name list> <right paren> ]
//
// <all fields column name list>    ::=   <column name list>

// SelectStatement represents a SELECT SQL statement
type SelectStatement struct {
	SelectList      []interface{}
	TableExpression *TableExpression
}
