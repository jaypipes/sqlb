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

func (j *JoinedTable) ArgCount(count *int) {
	if j.Qualified != nil {
		j := j.Qualified
		j.Left.ArgCount(count)
		j.Right.ArgCount(count)
		j.On.ArgCount(count)
	} else if j.Cross != nil {
		j := j.Cross
		j.Left.ArgCount(count)
		j.Right.ArgCount(count)
	} else if j.Union != nil {
		j := j.Union
		j.Left.ArgCount(count)
		j.Right.ArgCount(count)
	} else if j.Natural != nil {
		j := j.Natural
		j.Left.ArgCount(count)
		j.Right.ArgCount(count)
	}
}

type QualifiedJoin struct {
	Type  JoinType
	Left  TableReference
	Right TableReference
	On    BooleanValueExpression
}

func (j *QualifiedJoin) ArgCount(count *int) {
	j.Left.ArgCount(count)
	j.Right.ArgCount(count)
	j.On.ArgCount(count)
}

type NaturalJoin struct {
	Type  JoinType
	Left  TableReference
	Right TablePrimary
}

func (j *NaturalJoin) ArgCount(count *int) {
	j.Left.ArgCount(count)
	j.Right.ArgCount(count)
}

type CrossJoin struct {
	Left  TableReference
	Right TablePrimary
}

func (j *CrossJoin) ArgCount(count *int) {
	j.Left.ArgCount(count)
	j.Right.ArgCount(count)
}

type UnionJoin struct {
	Left  TableReference
	Right TablePrimary
}

func (j *UnionJoin) ArgCount(count *int) {
	j.Left.ArgCount(count)
	j.Right.ArgCount(count)
}
