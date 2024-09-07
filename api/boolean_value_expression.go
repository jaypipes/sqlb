//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"
	"slices"

	"github.com/jaypipes/sqlb/grammar"
)

// BooleanValueExpressionFromAny evaluates the supplied interface argument and
// returns a *BooleanValueExpression if the supplied argument can be converted
// into a BooleanValueExpression, or nil if the conversion cannot be done.
func BooleanValueExpressionFromAny(
	subject interface{},
) *grammar.BooleanValueExpression {
	switch v := subject.(type) {
	case *grammar.BooleanValueExpression:
		return v
	case grammar.BooleanValueExpression:
		return &v
	case *grammar.BooleanTerm:
		return &grammar.BooleanValueExpression{
			Unary: v,
		}
	case grammar.BooleanTerm:
		return &grammar.BooleanValueExpression{
			Unary: &v,
		}
	}
	// predicates like "A = B" are themselves boolean value expressions...
	pred := PredicateFromAny(subject)
	if pred != nil {
		return &grammar.BooleanValueExpression{
			Unary: &grammar.BooleanTerm{
				Unary: &grammar.BooleanFactor{
					Test: grammar.BooleanTest{
						Primary: grammar.BooleanPrimary{
							Predicate: pred,
						},
					},
				},
			},
		}
	}
	return nil
}

// BooleanTermFromAny evaluates the supplied interface argument and returns a
// *BooleanTerm if the supplied argument can be converted into a BooleanTerm,
// or nil if the conversion cannot be done.
func BooleanTermFromAny(
	subject interface{},
) *grammar.BooleanTerm {
	switch v := subject.(type) {
	case *grammar.BooleanTerm:
		return v
	case grammar.BooleanTerm:
		return &v
	case *grammar.BooleanFactor:
		return &grammar.BooleanTerm{
			Unary: v,
		}
	case grammar.BooleanFactor:
		return &grammar.BooleanTerm{
			Unary: &v,
		}
	case *grammar.BooleanPrimary:
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: *v,
				},
			},
		}
	case grammar.BooleanPrimary:
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: v,
				},
			},
		}
	case *grammar.Predicate:
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: grammar.BooleanPrimary{
						Predicate: v,
					},
				},
			},
		}
	case *grammar.BooleanPredicand:
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: grammar.BooleanPrimary{
						BooleanPredicand: v,
					},
				},
			},
		}
	case grammar.BooleanPredicand:
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: grammar.BooleanPrimary{
						BooleanPredicand: &v,
					},
				},
			},
		}
	}
	// predicates like "A = B" are themselves boolean primaries...
	pred := PredicateFromAny(subject)
	if pred != nil {
		return &grammar.BooleanTerm{
			Unary: &grammar.BooleanFactor{
				Test: grammar.BooleanTest{
					Primary: grammar.BooleanPrimary{
						Predicate: pred,
					},
				},
			},
		}
	}
	return nil
}

// BooleanFactorFromAny evaluates the supplied interface argument and returns a
// *BooleanFactor if the supplied argument can be converted into a
// BooleanFactor, or nil if the conversion cannot be done.
func BooleanFactorFromAny(
	subject interface{},
) *grammar.BooleanFactor {
	switch v := subject.(type) {
	case *grammar.BooleanFactor:
		return v
	case grammar.BooleanFactor:
		return &v
	case *grammar.BooleanPrimary:
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: *v,
			},
		}
	case grammar.BooleanPrimary:
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: v,
			},
		}
	case *grammar.Predicate:
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: grammar.BooleanPrimary{
					Predicate: v,
				},
			},
		}
	case *grammar.BooleanPredicand:
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: grammar.BooleanPrimary{
					BooleanPredicand: v,
				},
			},
		}
	case grammar.BooleanPredicand:
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: grammar.BooleanPrimary{
					BooleanPredicand: &v,
				},
			},
		}
	}
	// predicates like "A = B" are themselves boolean primaries...
	pred := PredicateFromAny(subject)
	if pred != nil {
		return &grammar.BooleanFactor{
			Test: grammar.BooleanTest{
				Primary: grammar.BooleanPrimary{
					Predicate: pred,
				},
			},
		}
	}
	return nil
}

// BooleanPredicandFromAny evaluates the supplied interface argument and
// returns a *BooleanPredicand if the supplied argument can be converted
// into a BooleanPredicand, or nil if the conversion cannot be done.
func BooleanPredicandFromAny(
	v interface{},
) *grammar.BooleanPredicand {
	switch v := v.(type) {
	case *grammar.BooleanPredicand:
		return v
	case grammar.BooleanPredicand:
		return &v
	case *grammar.BooleanValueExpression:
		return &grammar.BooleanPredicand{
			ParenthesizedBooleanValueExpression: v,
		}
	case grammar.BooleanValueExpression:
		return &grammar.BooleanPredicand{
			ParenthesizedBooleanValueExpression: &v,
		}
	case *grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.BooleanPredicand{
			NonParenthesizedValueExpressionPrimary: v,
		}
	case grammar.NonParenthesizedValueExpressionPrimary:
		return &grammar.BooleanPredicand{
			NonParenthesizedValueExpressionPrimary: &v,
		}
	}
	return nil
}

// And accepts two things and returns a BooleanValueExpression ANDing the two
// things together. This boolean value expression can be passed to a Join or
// Where clause.
func And(
	leftAny interface{},
	rightAny interface{},
) grammar.BooleanValueExpression {
	left := BooleanTermFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanTerm",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := BooleanFactorFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanFactor",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return grammar.BooleanValueExpression{
		Unary: &grammar.BooleanTerm{
			AndLeft:  left,
			AndRight: right,
		},
	}
}

// Or accepts two things and returns an Element representing an OR expression
// that can be passed to a Join or Where clause.
func Or(
	leftAny interface{},
	rightAny interface{},
) grammar.BooleanValueExpression {
	left := BooleanValueExpressionFromAny(leftAny)
	if left == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			leftAny, leftAny,
		)
		panic(msg)
	}
	right := BooleanTermFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanTerm",
			rightAny, rightAny,
		)
		panic(msg)
	}
	return grammar.BooleanValueExpression{
		OrLeft:  left,
		OrRight: right,
	}
}

// ReferredFromBooleanValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied BooleanValueExpression.
func ReferredFromBooleanValueExpression(
	bve *grammar.BooleanValueExpression,
) []string {
	if bve.Unary != nil {
		return ReferredFromBooleanTerm(bve.Unary)
	} else {
		found := ReferredFromBooleanValueExpression(bve.OrLeft)
		found = append(found, ReferredFromBooleanTerm(bve.OrRight)...)
		return found
	}
}

// ReferredFromBooleanTerm returns a slice of string names of any tables or
// derived tables (subqueries in the FROM clause) that are referenced within a
// supplied BooleanTerm.
func ReferredFromBooleanTerm(
	bt *grammar.BooleanTerm,
) []string {
	if bt.Unary != nil {
		return ReferredFromBooleanFactor(bt.Unary)
	} else {
		found := ReferredFromBooleanTerm(bt.AndLeft)
		found = append(found, ReferredFromBooleanFactor(bt.AndRight)...)
		return found
	}
}

// ReferredFromBooleanFactor returns a slice of string names of any tables or
// derived tables (subqueries in the FROM clause) that are referenced within a
// supplied BooleanFactor.
func ReferredFromBooleanFactor(
	bf *grammar.BooleanFactor,
) []string {
	tp := bf.Test.Primary
	if tp.Predicate != nil {
		return ReferredFromPredicate(tp.Predicate)
	}
	return ReferredFromBooleanPredicand(tp.BooleanPredicand)
}

// ReferredFromBooleanPredicand returns a slice of string names of any tables or
// derived tables (subqueries in the FROM clause) that are referenced within a
// supplied BooleanPredicand.
func ReferredFromBooleanPredicand(
	bp *grammar.BooleanPredicand,
) []string {
	if bp.ParenthesizedBooleanValueExpression != nil {
		return ReferredFromBooleanValueExpression(bp.ParenthesizedBooleanValueExpression)
	}
	p := bp.NonParenthesizedValueExpressionPrimary
	if p.ColumnReference != nil {
		if len(p.ColumnReference.BasicIdentifierChain.Identifiers) > 0 {
			return slices.Clone(p.ColumnReference.BasicIdentifierChain.Identifiers[:len(p.ColumnReference.BasicIdentifierChain.Identifiers)-1])
		}
	}
	return []string{}
}
