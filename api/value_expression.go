//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import "github.com/jaypipes/sqlb/grammar"

// ValueExpressionFromAny evaluates the supplied interface argument and returns
// a *ValueExpression if the supplied argument can be converted into a
// ValueExpression, or nil if the conversion cannot be done.
func ValueExpressionFromAny(subject interface{}) *grammar.ValueExpression {
	switch v := subject.(type) {
	case *grammar.CommonValueExpression:
		return &grammar.ValueExpression{CommonValueExpression: v}
	case grammar.CommonValueExpression:
		return &grammar.ValueExpression{CommonValueExpression: &v}
	case *grammar.BooleanValueExpression:
		return &grammar.ValueExpression{BooleanValueExpression: v}
	case grammar.BooleanValueExpression:
		return &grammar.ValueExpression{BooleanValueExpression: &v}
	case *grammar.RowValueExpression:
		return &grammar.ValueExpression{RowValueExpression: v}
	case grammar.RowValueExpression:
		return &grammar.ValueExpression{RowValueExpression: &v}
	case *grammar.NumericValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				NumericValueExpression: v,
			},
		}
	case grammar.NumericValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				NumericValueExpression: &v,
			},
		}
	case *grammar.StringValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				StringValueExpression: v,
			},
		}
	case grammar.StringValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				StringValueExpression: &v,
			},
		}
	case *grammar.DatetimeValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				DatetimeValueExpression: v,
			},
		}
	case grammar.DatetimeValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				DatetimeValueExpression: &v,
			},
		}
	case *grammar.IntervalValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				IntervalValueExpression: v,
			},
		}
	case grammar.IntervalValueExpression:
		return &grammar.ValueExpression{
			CommonValueExpression: &grammar.CommonValueExpression{
				IntervalValueExpression: &v,
			},
		}
	case *grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.ValueExpression{
			RowValueExpression: &grammar.RowValueExpression{
				NonParenthesizedValueExpressionPrimary: v,
			},
		}
	case grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.ValueExpression{
			RowValueExpression: &grammar.RowValueExpression{
				NonParenthesizedValueExpressionPrimary: &v,
			},
		}
	}
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.ValueExpression{
			RowValueExpression: &grammar.RowValueExpression{
				NonParenthesizedValueExpressionPrimary: v,
			},
		}
	}
	return nil
}

// ReferredFromCommonValueExpression returns a slice of string names of any tables
// or derived tables (subqueries in the FROM clause) that are referenced within
// a supplied CommonValueExpression.
func ReferredFromCommonValueExpression(
	cve *grammar.CommonValueExpression,
) []string {
	if cve.NumericValueExpression != nil {
		return ReferredFromNumericValueExpression(cve.NumericValueExpression)
	} else if cve.StringValueExpression != nil {
		return ReferredFromStringValueExpression(cve.StringValueExpression)
	} else if cve.DatetimeValueExpression != nil {
		return ReferredFromDatetimeValueExpression(cve.DatetimeValueExpression)
	} else if cve.IntervalValueExpression != nil {
		return ReferredFromIntervalValueExpression(cve.IntervalValueExpression)
	}
	return []string{}
}
