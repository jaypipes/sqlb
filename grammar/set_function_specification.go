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
	AggregateFunction *AggregateFunction
	//GroupingOperation *GroupingOperation
}
