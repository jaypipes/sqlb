//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import (
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/core/grammar"
)

// RowValuePredicandFromAny evaluates the supplied interface argument and
// returns a *RowValuePredicand if the supplied argument can be converted into
// a RowValuePredicand, or nil if the conversion cannot be done.
func RowValuePredicandFromAny(subject interface{}) *grammar.RowValuePredicand {
	switch v := subject.(type) {
	case *grammar.RowValuePredicand:
		return v
	case grammar.RowValuePredicand:
		return &v
	case *grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.RowValuePredicand{
			NonParenthesizedValueExpressionPrimary: v,
		}
	case grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.RowValuePredicand{
			NonParenthesizedValueExpressionPrimary: &v,
		}
	case *grammar.CommonValueExpression:
		return &grammar.RowValuePredicand{
			CommonValueExpression: v,
		}
	case grammar.CommonValueExpression:
		return &grammar.RowValuePredicand{
			CommonValueExpression: &v,
		}
	case types.CommonValueExpressionConverter:
		return &grammar.RowValuePredicand{
			CommonValueExpression: v.CommonValueExpression(),
		}
	case *grammar.BooleanPredicand:
		return &grammar.RowValuePredicand{
			BooleanPredicand: v,
		}
	case grammar.BooleanPredicand:
		return &grammar.RowValuePredicand{
			BooleanPredicand: &v,
		}
	}
	// We could have a simple literal passed to us. See if we can convert it
	// into a simple row value predicand with a non-parenthesized value
	// expression primary
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.RowValuePredicand{
			NonParenthesizedValueExpressionPrimary: v,
		}
	}
	return nil
}

// ReferredFromRowValuePredicand returns a slice of string names of any tables
// or derived tables (subqueries in the FROM clause) that are referenced within
// a supplied RowValuePredicand.
func ReferredFromRowValuePredicand(
	rvp *grammar.RowValuePredicand,
) []string {
	if rvp.NonParenthesizedValueExpressionPrimary != nil {
		return ReferredFromNonParenthesizedValueExpressionPrimary(rvp.NonParenthesizedValueExpressionPrimary)
	} else if rvp.CommonValueExpression != nil {
		return ReferredFromCommonValueExpression(rvp.CommonValueExpression)
	} else if rvp.BooleanPredicand != nil {
		return ReferredFromBooleanPredicand(rvp.BooleanPredicand)
	}
	return []string{}
}
