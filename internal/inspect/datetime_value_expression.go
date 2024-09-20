//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/types"
)

// DatetimeValueExpressionFromAny evaluates the supplied interface argument and
// returns a *DatetimeValueExpression if the supplied argument can be converted
// into a DatetimeValueExpression, or nil if the conversion cannot be done.
func DatetimeValueExpressionFromAny(
	subject interface{},
) *grammar.DatetimeValueExpression {
	switch v := subject.(type) {
	case *grammar.DatetimeValueExpression:
		return v
	case grammar.DatetimeValueExpression:
		return &v
	case *grammar.CommonValueExpression:
		if v.Datetime != nil {
			return v.Datetime
		}
	// Columns can produce datetime values...
	case types.ColumnReferenceConverter:
		return &grammar.DatetimeValueExpression{
			Unary: &grammar.DatetimeTerm{
				Factor: grammar.DatetimeFactor{
					Primary: grammar.DatetimePrimary{
						Primary: &grammar.ValueExpressionPrimary{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: v.ColumnReference(),
							},
						},
					},
				},
			},
		}
	}
	return nil
}

// IntervalValueExpressionFromAny evaluates the supplied interface argument and
// returns a *IntervalValueExpression if the supplied argument can be converted
// into a IntervalValueExpression, or nil if the conversion cannot be done.
func IntervalValueExpressionFromAny(
	subject interface{},
) *grammar.IntervalValueExpression {
	switch v := subject.(type) {
	case *grammar.IntervalValueExpression:
		return v
	case grammar.IntervalValueExpression:
		return &v
	case *grammar.CommonValueExpression:
		if v.Interval != nil {
			return v.Interval
		}
	// Columns can produce interval values...
	case types.ColumnReferenceConverter:
		return &grammar.IntervalValueExpression{
			Unary: &grammar.IntervalTerm{
				Unary: &grammar.IntervalFactor{
					Primary: grammar.IntervalPrimary{
						Primary: &grammar.ValueExpressionPrimary{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: v.ColumnReference(),
							},
						},
					},
				},
			},
		}
	}
	return nil
}

// ReferredFromDatetimeValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied DatetimeValueExpression.
func ReferredFromDatetimeValueExpression(
	cve *grammar.DatetimeValueExpression,
) []string {
	return []string{}
}
