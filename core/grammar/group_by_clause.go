//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <group by clause>    ::=   GROUP BY [ <set quantifier> ] <grouping element list>
//
// <grouping element list>    ::=   <grouping element> [ { <comma> <grouping element> }... ]
//
// <grouping element>    ::=
//          <ordinary grouping set>
//      |     <rollup list>
//      |     <cube list>
//      |     <grouping sets specification>
//      |     <empty grouping set>
//
// <ordinary grouping set>    ::=
//          <grouping column reference>
//      |     <left paren> <grouping column reference list> <right paren>
//
// <grouping column reference>    ::=   <column reference> [ <collate clause> ]
//
// <grouping column reference list>    ::=   <grouping column reference> [ { <comma> <grouping column reference> }... ]
//
// <rollup list>    ::=   ROLLUP <left paren> <ordinary grouping set list> <right paren>
//
// <ordinary grouping set list>    ::=   <ordinary grouping set> [ { <comma> <ordinary grouping set> }... ]
//
// <cube list>    ::=   CUBE <left paren> <ordinary grouping set list> <right paren>
//
// <grouping sets specification>    ::=   GROUPING SETS <left paren> <grouping set list> <right paren>
//
// <grouping set list>    ::=   <grouping set> [ { <comma> <grouping set> }... ]
//
// <grouping set>    ::=
//          <ordinary grouping set>
//      |     <rollup list>
//      |     <cube list>
//      |     <grouping sets specification>
//      |     <empty grouping set>
//
// <empty grouping set>    ::=   <left paren> <right paren>

// GroupByClause represents the SQL GROUP BY clause
type GroupByClause struct {
	GroupingElements []GroupingElement
}

func (c *GroupByClause) ArgCount(count *int) {
	for _, ge := range c.GroupingElements {
		ge.ArgCount(count)
	}
}

type GroupingElement struct {
	OrdinaryGroupingSet *OrdinaryGroupingSet
	//Rollup []OrdinaryGroupingSet
	//Cube []OrdinaryGroupingSet
	//GroupingSetsSpecification *GroupingSetsSpecification
}

func (e *GroupingElement) ArgCount(count *int) {
	if e.OrdinaryGroupingSet != nil {
		e.OrdinaryGroupingSet.ArgCount(count)
	}
}

type OrdinaryGroupingSet struct {
	GroupingColumnReference *GroupingColumnReference
}

func (s *OrdinaryGroupingSet) ArgCount(count *int) {
	if s.GroupingColumnReference != nil {
		s.GroupingColumnReference.ArgCount(count)
	}
}

type GroupingColumnReference struct {
	ColumnReference *ColumnReference
	Collation       *string
}

func (r *GroupingColumnReference) ArgCount(count *int) {
	// Column references don't produce query arguments
}
