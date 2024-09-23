//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn

import (
	"fmt"

	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/internal/inspect"
)

// AggregateFunction describes a SQL aggregate function (COUNT, AVG, SUM, etc)
// across zero or more referenced tables/columns/value expressions.
type AggregateFunction struct {
	BaseFunction
	*grammar.AggregateFunction
}

// DerivedColumn returns the `*grammar.DerivedColumn` element representing
// the Projection
func (f *AggregateFunction) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		Value: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					SetFunction: &grammar.SetFunctionSpecification{
						Aggregate: f.AggregateFunction,
					},
				},
			},
		},
	}
	if f.alias != "" {
		dc.As = &f.alias
	}
	return dc
}

// As aliases the SQL function as the supplied column name
func (f *AggregateFunction) As(alias string) types.Projection {
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
	if f.AggregateFunction.GeneralSet == nil {
		return f
	}
	f.AggregateFunction.GeneralSet.Quantifier = grammar.SetQuantifierDistinct
	return f
}

// Aggregate returns an AggregateFunction with the supplied subject and
// computational operation.
func Aggregate(
	subjectAny interface{},
	op grammar.ComputationalOperation,
) *AggregateFunction {
	var ref types.Relation
	switch subjectAny := subjectAny.(type) {
	case types.Projection:
		ref = subjectAny.References()
	case *grammar.ValueExpression:
		return &AggregateFunction{
			AggregateFunction: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: op,
					Value:     *subjectAny,
				},
			},
		}
	}
	v := inspect.ValueExpressionFromAny(subjectAny)
	if v == nil {
		msg := fmt.Sprintf(
			"expected coerceable ValueExpression but got %+v(%T)",
			subjectAny, subjectAny,
		)
		panic(msg)
	}
	return &AggregateFunction{
		BaseFunction: BaseFunction{
			ref: ref,
		},
		AggregateFunction: &grammar.AggregateFunction{
			GeneralSet: &grammar.GeneralSetFunction{
				Operation: op,
				Value:     *v,
			},
		},
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
	return Aggregate(args[0], grammar.ComputationalOperationCount)
}

// CountStar returns an AggregateFunction that produces a COUNT(*) against the
// supplied selectable thing.
func CountStar(sel types.Relation) *AggregateFunction {
	return &AggregateFunction{
		BaseFunction: BaseFunction{
			ref: sel,
		},
		AggregateFunction: &grammar.AggregateFunction{
			CountStar: &struct{}{},
		},
	}
}

// Avg returns a AggregateFunction that can be passed to a Select function to
// create a AVG(<value expression>) SQL function. The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Avg(subjectAny interface{}) *AggregateFunction {
	return Aggregate(subjectAny, grammar.ComputationalOperationAvg)
}

// Min returns a AggregateFunction that can be passed to a Select function to
// create a MIN(<value expression>) SQL function. The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Min(subjectAny interface{}) *AggregateFunction {
	return Aggregate(subjectAny, grammar.ComputationalOperationMin)
}

// Max returns a AggregateFunction that can be passed to a Select function to
// create a MAX(<value expression>) SQL function. The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Max(subjectAny interface{}) *AggregateFunction {
	return Aggregate(subjectAny, grammar.ComputationalOperationMax)
}

// Sum returns a AggregateFunction that can be passed to a Select function to
// create a SUM(<value expression>) SQL function. The supplied argument should
// be a ValueExpression or something that can be converted into a
// ValueExpression.
func Sum(subjectAny interface{}) *AggregateFunction {
	return Aggregate(subjectAny, grammar.ComputationalOperationSum)
}
