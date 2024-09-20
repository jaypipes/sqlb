//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/grammar"

// ValueExpressionFromAny evaluates the supplied interface argument and returns
// a *ValueExpression if the supplied argument can be converted into a
// ValueExpression, or nil if the conversion cannot be done.
func ValueExpressionFromAny(subject interface{}) *grammar.ValueExpression {
	switch v := subject.(type) {
	case *grammar.ValueExpression:
		return v
	case grammar.ValueExpression:
		return &v
	case *grammar.CommonValueExpression:
		return &grammar.ValueExpression{Common: v}
	case grammar.CommonValueExpression:
		return &grammar.ValueExpression{Common: &v}
	case *grammar.BooleanValueExpression:
		return &grammar.ValueExpression{Boolean: v}
	case grammar.BooleanValueExpression:
		return &grammar.ValueExpression{Boolean: &v}
	case *grammar.RowValueExpression:
		return &grammar.ValueExpression{Row: v}
	case grammar.RowValueExpression:
		return &grammar.ValueExpression{Row: &v}
	case *grammar.NumericValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Numeric: v,
			},
		}
	case grammar.NumericValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Numeric: &v,
			},
		}
	case *grammar.StringValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				String: v,
			},
		}
	case grammar.StringValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				String: &v,
			},
		}
	case *grammar.DatetimeValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Datetime: v,
			},
		}
	case grammar.DatetimeValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Datetime: &v,
			},
		}
	case *grammar.IntervalValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Interval: v,
			},
		}
	case grammar.IntervalValueExpression:
		return &grammar.ValueExpression{
			Common: &grammar.CommonValueExpression{
				Interval: &v,
			},
		}
	}
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: v,
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
	if cve.Numeric != nil {
		return ReferredFromNumericValueExpression(cve.Numeric)
	} else if cve.String != nil {
		return ReferredFromStringValueExpression(cve.String)
	} else if cve.Datetime != nil {
		return ReferredFromDatetimeValueExpression(cve.Datetime)
	} else if cve.Interval != nil {
		return ReferredFromIntervalValueExpression(cve.Interval)
	}
	return []string{}
}
