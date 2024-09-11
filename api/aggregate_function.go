//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"

	"github.com/jaypipes/sqlb/grammar"
)

// AggregateFunction describes a SQL aggregate function (COUNT, AVG, SUM, etc)
// across zero or more referenced tables/columns/value expressions.
type AggregateFunction struct {
	*grammar.AggregateFunction
	// referred is a the Table or DerivedTable that is referred from
	// the aggregate function
	Referred interface{}
	// alias is the aggregate function as an aliased projection
	// (e.g. COUNT(*) AS counter)
	alias string
}

// As aliases the SQL function as the supplied column name
func (f *AggregateFunction) As(alias string) *AggregateFunction {
	f.alias = alias
	return f
}

// Distinct modifies the AggregateFunction by changing the <set quantifier>
// from ALL to DISTINCT. This does nothing unless the AggregateFunction is a
// General Set Function, which is an aggregate function with any of the
// following computational operations: AVG, MAX, MIN, SUM, EVERY, ANY, SOME,
// COUNT, STDDEV_POP, STDDEV_SAMP, VAR_SAMP, VAR_POP, COLLECT, FUSION,
// INTERSECTION.
func (f *AggregateFunction) Distinct() *AggregateFunction {
	if f.AggregateFunction.GeneralSetFunction == nil {
		return f
	}
	f.AggregateFunction.GeneralSetFunction.Quantifier = grammar.SetQuantifierDistinct
	return f
}

func doGeneralSetFunction(op grammar.ComputationalOperation, valAny interface{}) *AggregateFunction {
	var ref interface{}
	switch valAny := valAny.(type) {
	case *Column:
		ref = valAny.t
	case *grammar.ValueExpression:
		return &AggregateFunction{
			AggregateFunction: &grammar.AggregateFunction{
				GeneralSetFunction: &grammar.GeneralSetFunction{
					Operation:       op,
					ValueExpression: *valAny,
				},
			},
		}
	case grammar.ValueExpression:
		return &AggregateFunction{
			AggregateFunction: &grammar.AggregateFunction{
				GeneralSetFunction: &grammar.GeneralSetFunction{
					Operation:       op,
					ValueExpression: valAny,
				},
			},
		}
	}
	v := ValueExpressionFromAny(valAny)
	if v == nil {
		msg := fmt.Sprintf(
			"expected coerceable ValueExpression but got %+v(%T)",
			valAny, valAny,
		)
		panic(msg)
	}
	return &AggregateFunction{
		AggregateFunction: &grammar.AggregateFunction{
			GeneralSetFunction: &grammar.GeneralSetFunction{
				Operation:       op,
				ValueExpression: *v,
			},
		},
		Referred: ref,
	}
}

// Count returns a AggregateFunction that can be passed to a Select function.
// It accepts zero or one parameter. If no parameters are passed, the
// AggregateFunction returned represents a COUNT(*) SQL function. If a
// parameter is passed, it should be a ValueExpression or something that can be
// converted into a ValueExpression.
func Count(args ...interface{}) *AggregateFunction {
	if len(args) > 1 {
		panic("Count expects either zero or one argument")
	}
	if len(args) == 0 {
		return &AggregateFunction{
			AggregateFunction: &grammar.AggregateFunction{
				CountStar: &struct{}{},
			},
		}
	}
	return doGeneralSetFunction(grammar.ComputationalOperationCount, args[0])
}

// Avg returns a AggregateFunction that can be passed to a Select function to
// create a AVG(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Avg(valAny interface{}) *AggregateFunction {
	return doGeneralSetFunction(grammar.ComputationalOperationAvg, valAny)
}

// Min returns a AggregateFunction that can be passed to a Select function to
// create a MIN(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Min(valAny interface{}) *AggregateFunction {
	return doGeneralSetFunction(grammar.ComputationalOperationMin, valAny)
}

// Max returns a AggregateFunction that can be passed to a Select function to
// create a MAX(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Max(valAny interface{}) *AggregateFunction {
	return doGeneralSetFunction(grammar.ComputationalOperationMax, valAny)
}

// Sum returns a AggregateFunction that can be passed to a Select function to
// create a SUM(<value expression>) SQL function.  The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Sum(valAny interface{}) *AggregateFunction {
	return doGeneralSetFunction(grammar.ComputationalOperationSum, valAny)
}
