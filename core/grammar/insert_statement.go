//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <insert statement>    ::=   INSERT INTO <insertion target> <insert columns and source>
//
// <insertion target>    ::=   <table name>
//
// <insert columns and source>    ::=
//          <from subquery>
//      |     <from constructor>
//      |     <from default>
//
// <from subquery>    ::=   [ <left paren> <insert column list> <right paren> ] [ <override clause> ] <query expression>
//
// <from constructor>    ::=
//          [ <left paren> <insert column list> <right paren> ] [ <override clause> ] <contextually typed table value constructor>
//
// <override clause>    ::=   OVERRIDING USER VALUE | OVERRIDING SYSTEM VALUE
//
// <from default>    ::=   DEFAULT VALUES
//
// <insert column list>    ::=   <column name list>

// InsertStatement represents an INSERT SQL statement
type InsertStatement struct {
	TableName string
	Columns   []string
	Values    []interface{}
}
