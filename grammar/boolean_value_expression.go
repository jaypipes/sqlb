//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <boolean value expression>    ::=
//          <boolean term>
//      |     <boolean value expression> OR <boolean term>
//
// <boolean term>    ::=
//          <boolean factor>
//      |     <boolean term> AND <boolean factor>
//
// <boolean factor>    ::=   [ NOT ] <boolean test>
//
// <boolean test>    ::=   <boolean primary> [ IS [ NOT ] <truth value> ]
//
// <truth value>    ::=   TRUE | FALSE | UNKNOWN
//
// <boolean primary>    ::=   <predicate> | <boolean predicand>
//
// <boolean predicand>    ::=
//          <parenthesized boolean value expression>
//      |     <nonparenthesized value expression primary>
//
// <parenthesized boolean value expression>    ::=   <left paren> <boolean value expression> <right paren>

// BooleanValueExpression represents a boolean comparison expression in the SQL
// statement, e.g. "a = b"
type BooleanValueExpression struct {
	Terms []interface{}
}
