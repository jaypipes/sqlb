//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/core/grammar"

// NumericValueExpressionFromAny evaluates the supplied interface argument and
// returns a *NumericValueExpression if the supplied argument can be converted
// into a NumericValueExpression, or nil if the conversion cannot be done.
func NumericValueExpressionFromAny(subject interface{}) *grammar.NumericValueExpression {
	switch v := subject.(type) {
	case *grammar.NumericValueExpression:
		return v
	case grammar.NumericValueExpression:
		return &v
	case *grammar.Term:
		return &grammar.NumericValueExpression{
			Unary: v,
		}
	case grammar.Term:
		return &grammar.NumericValueExpression{
			Unary: &v,
		}
	case *grammar.Factor:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: v,
			},
		}
	case grammar.Factor:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &v,
			},
		}
	case *grammar.NumericValueFunction:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: v,
					},
				},
			},
		}
	case grammar.NumericValueFunction:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Function: &v,
					},
				},
			},
		}
	case *grammar.ValueExpressionPrimary:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Primary: v,
					},
				},
			},
		}
	case grammar.ValueExpressionPrimary:
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Primary: &v,
					},
				},
			},
		}
	}
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.NumericValueExpression{
			Unary: &grammar.Term{
				Unary: &grammar.Factor{
					Primary: grammar.NumericPrimary{
						Primary: &grammar.ValueExpressionPrimary{
							Primary: v,
						},
					},
				},
			},
		}
	}
	return nil
}

// ReferredFromNumericValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied NumericValueExpression.
func ReferredFromNumericValueExpression(
	cve *grammar.NumericValueExpression,
) []string {
	return []string{}
}
