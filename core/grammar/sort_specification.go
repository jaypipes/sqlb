//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type OrderSpecification int

const (
	OrderSpecificationAsc OrderSpecification = iota
	OrderSpecificationDesc
)

type NullOrderSpecification int

const (
	NullOrderSpecificationNone NullOrderSpecification = iota
	NullOrderSpecificationNullsFirst
	NullOrderSpecificationNullsLast
)

var NullOrderSpecificationSymbol = map[NullOrderSpecification]string{
	NullOrderSpecificationNullsFirst: "NULLS FIRST",
	NullOrderSpecificationNullsLast:  "NULLS LAST",
}

// <sort specification list>    ::=   <sort specification> [ { <comma> <sort specification> }... ]
//
// <sort specification>    ::=   <sort key> [ <ordering specification> ] [ <null ordering> ]
//
// <sort key>    ::=   <value expression>
//
// <ordering specification>    ::=   ASC | DESC
//
// <null ordering>    ::=   NULLS FIRST | NULLS LAST

type SortSpecification struct {
	Key       ValueExpression
	Order     OrderSpecification
	NullOrder NullOrderSpecification
}
