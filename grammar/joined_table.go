//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type JoinType int

const (
	JoinTypeInner JoinType = iota
	JoinTypeLeftOuter
	JoinTypeRightOuter
	JoinTypeFullOuter
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
	Cross     *CrossJoin
	Qualified *QualifiedJoin
	Natural   *NaturalJoin
	Union     *UnionJoin
}

type QualifiedJoin struct {
	Type  JoinType
	Left  TableReference
	Right TableReference
	On    BooleanValueExpression
}

type NaturalJoin struct {
	Type  JoinType
	Left  TableReference
	Right TablePrimary
}

type CrossJoin struct {
	Left  TableReference
	Right TablePrimary
}

type UnionJoin struct {
	Left  TableReference
	Right TablePrimary
}
