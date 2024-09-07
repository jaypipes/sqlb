//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <predicate>    ::=
//          <comparison predicate>
//      |     <between predicate>
//      |     <in predicate>
//      |     <like predicate>
//      |     <similar predicate>
//      |     <null predicate>
//      |     <quantified comparison predicate>
//      |     <exists predicate>
//      |     <unique predicate>
//      |     <normalized predicate>
//      |     <match predicate>
//      |     <overlaps predicate>
//      |     <distinct predicate>
//      |     <member predicate>
//      |     <submultiset predicate>
//      |     <set predicate>
//      |     <type predicate>

type Predicate struct {
	Comparison *ComparisonPredicate
	Between    *BetweenPredicate
	In         *InPredicate
	//Like *LikePredicate
	//Similar *SimilarPredicate
	Null *NullPredicate
	//QuantifiedComparison *QuantifiedComparisonPredicate
	//Exists *ExistsPredicate
	//Unique *UniquePredicate
	//Normalized *NormalizedPredicate
	//Match *MatchPredicate
	//Overlaps *OverlapsPredicate
	//Distinct *DistinctPredicate
	//Member *MemberPredicate
	//Submultiset *SubmultisetPredicate
	//Set *SetPredicate
	//Type *TypePredicate
}

// <comparison predicate>    ::=   <row value predicand> <comparison predicate part 2>
//
// <comparison predicate part 2>    ::=   <comp op> <row value predicand>
//
// <comp op>    ::=
//          <equals operator>
//      |     <not equals operator>
//      |     <less than operator>
//      |     <greater than operator>
//      |     <less than or equals operator>
//      |     <greater than or equals operator>

type ComparisonOperator int

const (
	ComparisonOperatorEquals ComparisonOperator = iota
	ComparisonOperatorNotEquals
	ComparisonOperatorLessThan
	ComparisonOperatorGreaterThan
	ComparisonOperatorLessThanEquals
	ComparisonOperatorGreaterThanEquals
)

type ComparisonPredicate struct {
	Operator ComparisonOperator
	A        RowValuePredicand
	B        RowValuePredicand
}

// <between predicate>    ::=   <row value predicand> <between predicate part 2>
//
// <between predicate part 2>    ::=   [ NOT ] BETWEEN [ ASYMMETRIC | SYMMETRIC ] <row value predicand> AND <row value predicand>

type BetweenPredicate struct {
	Target RowValuePredicand
	Start  RowValuePredicand
	End    RowValuePredicand
}

// <in predicate>    ::=   <row value predicand> <in predicate part 2>
//
// <in predicate part 2>    ::=   [ NOT ] IN <in predicate value>
//
// <in predicate value>    ::=
//          <table subquery>
//      |     <left paren> <in value list> <right paren>
//
// <in value list>    ::=   <row value expression> [ { <comma> <row value expression> }... ]

type InPredicate struct {
	Target RowValuePredicand
	Values []RowValueExpression
}

// <null predicate>    ::=   <row value predicand> <null predicate part 2>
//
// <null predicate part 2>    ::=   IS [ NOT ] NULL

type NullPredicate struct {
	Target RowValuePredicand
	Not    bool
}
