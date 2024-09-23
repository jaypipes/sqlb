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
			Primary: v,
		}
	case *grammar.CommonValueExpression:
		return &grammar.RowValuePredicand{
			Common: v,
		}
	case types.CommonValueExpressionConverter:
		return &grammar.RowValuePredicand{
			Common: v.CommonValueExpression(),
		}
	case *grammar.BooleanPredicand:
		return &grammar.RowValuePredicand{
			Boolean: v,
		}
	}
	// We could have a simple literal passed to us. See if we can convert it
	// into a simple row value predicand with a non-parenthesized value
	// expression primary
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.RowValuePredicand{
			Primary: v,
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
	if rvp.Primary != nil {
		return ReferredFromNonParenthesizedValueExpressionPrimary(rvp.Primary)
	} else if rvp.Common != nil {
		return ReferredFromCommonValueExpression(rvp.Common)
	} else if rvp.Boolean != nil {
		return ReferredFromBooleanPredicand(rvp.Boolean)
	}
	return []string{}
}
