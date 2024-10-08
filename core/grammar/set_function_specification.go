//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <set function specification>    ::=   <aggregate function> | <grouping operation>
//
// <grouping operation>    ::=   GROUPING <left paren> <column reference> [ { <comma> <column reference> }... ] <right paren>

type SetFunctionSpecification struct {
	Aggregate *AggregateFunction
	//GroupingOperation *GroupingOperation
}

func (s *SetFunctionSpecification) ArgCount(count *int) {
	if s.Aggregate != nil {
		s.Aggregate.ArgCount(count)
	}
}
