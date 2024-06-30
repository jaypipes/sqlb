//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type JoinType int

const (
	JoinInner JoinType = iota
	JoinOuter
	JoinCross
)

// <joined table>    ::=
//          <cross join>
//      |     <qualified join>
//      |     <natural join>
//      |     <union join>
//
// <cross join>    ::=   <table reference> CROSS JOIN <table primary>
//
// <qualified join>    ::=   <table reference> [ <join type> ] JOIN <table reference> <join specification>
//
// <natural join>    ::=   <table reference> NATURAL [ <join type> ] JOIN <table primary>
//
// <union join>    ::=   <table reference> UNION JOIN <table primary>
//
// <join specification>    ::=   <join condition> | <named columns join>
//
// <join condition>    ::=   ON <search condition>
//
// <named columns join>    ::=   USING <left paren> <join column list> <right paren>
//
// <join type>    ::=   INNER | <outer join type> [ OUTER ]
//
// <outer join type>    ::=   LEFT | RIGHT | FULL
//
// <join column list>    ::=   <column name list>

// JoinedTable represents a table derived from a Cartesian product, inner or outer join, or union join.
type JoinedTable struct {
	Type  JoinType
	Left  *TableReference
	right *TableReference
	On    *BooleanValueExpression
}
