//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <aggregate function>    ::=
//          COUNT <left paren> <asterisk> <right paren> [ <filter clause> ]
//      |     <general set function> [ <filter clause> ]
//      |     <binary set function> [ <filter clause> ]
//      |     <ordered set function> [ <filter clause> ]
//
// <general set function>    ::=   <set function type> <left paren> [ <set quantifier> ] <value expression> <right paren>
//
// <set function type>    ::=   <computational operation>
//
// <computational operation>    ::=
//            AVG | MAX | MIN | SUM
//      |     EVERY | ANY | SOME
//      |     COUNT
//      |     STDDEV_POP | STDDEV_SAMP | VAR_SAMP | VAR_POP
//      |     COLLECT | FUSION | INTERSECTION
//
// <set quantifier>    ::=   DISTINCT | ALL

type ComputationalOperation int

const (
	ComputationalOperationAvg ComputationalOperation = iota
	ComputationalOperationMax
	ComputationalOperationMin
	ComputationalOperationSum
	ComputationalOperationEvery
	ComputationalOperationAny
	ComputationalOperationSome
	ComputationalOperationCount
	ComputationalOperationStdDevPop
	ComputationalOperationStdDevSamp
	ComputationalOperationVarSamp
	ComputationalOperationVarPop
	ComputationalOperationCollect
	ComputationalOperationFusion
	ComputationalOperationIntersection
)

var ComputationalOperationSymbol = map[ComputationalOperation]string{
	ComputationalOperationAvg:          "AVG",
	ComputationalOperationMax:          "MAX",
	ComputationalOperationMin:          "MIN",
	ComputationalOperationSum:          "SUM",
	ComputationalOperationEvery:        "EVERY",
	ComputationalOperationAny:          "ANY",
	ComputationalOperationSome:         "SOME",
	ComputationalOperationCount:        "COUNT",
	ComputationalOperationStdDevPop:    "STDDEV_POP",
	ComputationalOperationStdDevSamp:   "STDDEV_SAMP",
	ComputationalOperationVarSamp:      "VAR_SAMP",
	ComputationalOperationVarPop:       "VAR_POP",
	ComputationalOperationCollect:      "COLLECT",
	ComputationalOperationFusion:       "FUSION",
	ComputationalOperationIntersection: "INTERSECTION",
}

type SetQuantifier int

const (
	SetQuantifierAll SetQuantifier = iota
	SetQuantifierDistinct
)

type AggregateFunction struct {
	CountStar  *struct{}
	GeneralSet *GeneralSetFunction
	BinarySet  *BinarySetFunction
	OrderedSet *OrderedSetFunction
}

func (f *AggregateFunction) ArgCount(count *int) {
	if f.GeneralSet != nil {
		f.GeneralSet.ArgCount(count)
	} else if f.BinarySet != nil {
		f.BinarySet.ArgCount(count)
	} else if f.OrderedSet != nil {
		f.OrderedSet.ArgCount(count)
	}
}

type GeneralSetFunction struct {
	Operation  ComputationalOperation
	Quantifier SetQuantifier
	Value      ValueExpression
}

func (f *GeneralSetFunction) ArgCount(count *int) {
	f.Value.ArgCount(count)
}

// <binary set function>    ::=   <binary set function type> <left paren> <dependent variable expression> <comma> <independent variable expression> <right paren>
//
// <binary set function type>    ::=
//          COVAR_POP | COVAR_SAMP | CORR | REGR_SLOPE
//      |     REGR_INTERCEPT | REGR_COUNT | REGR_R2 | REGR_AVGX | REGR_AVGY
//      |     REGR_SXX | REGR_SYY | REGR_SXY
//
// <dependent variable expression>    ::=   <numeric value expression>
//
// <independent variable expression>    ::=   <numeric value expression>

type BinarySetFunction struct{}

func (f *BinarySetFunction) ArgCount(count *int) {
}

// <ordered set function>    ::=   <hypothetical set function> | <inverse distribution function>
//
// <hypothetical set function>    ::=   <rank function type> <left paren> <hypothetical set function value expression list> <right paren> <within group specification>
//
// <within group specification>    ::=   WITHIN GROUP <left paren> ORDER BY <sort specification list> <right paren>
//
// <hypothetical set function value expression list>    ::=   <value expression> [ { <comma> <value expression> }... ]
//
// <inverse distribution function>    ::=   <inverse distribution function type> <left paren> <inverse distribution function argument> <right paren> <within group specification>
//
// <inverse distribution function argument>    ::=   <numeric value expression>
//
// <inverse distribution function type>    ::=   PERCENTILE_CONT | PERCENTILE_DISC

type OrderedSetFunction struct{}

func (f *OrderedSetFunction) ArgCount(count *int) {
}
