//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

// PredicateFromAny evaluates the supplied interface argument and
// returns a *Predicate if the supplied argument can be converted
// into a Predicate, or nil if the conversion cannot be done.
func PredicateFromAny(v interface{}) *grammar.Predicate {
	switch v := v.(type) {
	case *grammar.Predicate:
		return v
	case grammar.Predicate:
		return &v
	case *grammar.ComparisonPredicate:
		return &grammar.Predicate{
			Comparison: v,
		}
	case grammar.ComparisonPredicate:
		return &grammar.Predicate{
			Comparison: &v,
		}
	case *grammar.BetweenPredicate:
		return &grammar.Predicate{
			Between: v,
		}
	case grammar.BetweenPredicate:
		return &grammar.Predicate{
			Between: &v,
		}
	case *grammar.InPredicate:
		return &grammar.Predicate{
			In: v,
		}
	case grammar.InPredicate:
		return &grammar.Predicate{
			In: &v,
		}
	case *grammar.NullPredicate:
		return &grammar.Predicate{
			Null: v,
		}
	case grammar.NullPredicate:
		return &grammar.Predicate{
			Null: &v,
		}
	}
	return nil
}

// ReferredFromPredicate returns a slice of string names of any tables or
// derived tables (subqueries in the FROM clause) that are referenced within a
// supplied Predicate.
func ReferredFromPredicate(
	p *grammar.Predicate,
) []string {
	if p.Comparison != nil {
		found := ReferredFromRowValuePredicand(&p.Comparison.A)
		found = append(found, ReferredFromRowValuePredicand(&p.Comparison.B)...)
		return found
	} else if p.In != nil {
		found := ReferredFromRowValuePredicand(&p.In.Target)
		for _, rve := range p.In.Values {
			found = append(
				found,
				ReferredFromNonParenthesizedValueExpressionPrimary(rve.Primary)...,
			)
		}
		return found
	} else if p.Between != nil {
		found := ReferredFromRowValuePredicand(&p.Between.Target)
		found = append(found, ReferredFromRowValuePredicand(&p.Between.Start)...)
		found = append(found, ReferredFromRowValuePredicand(&p.Between.End)...)
		return found
	} else if p.Null != nil {
		return ReferredFromRowValuePredicand(&p.Null.Target)
	}
	return []string{}
}
