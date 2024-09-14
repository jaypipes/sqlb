//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/grammar"

// StringValueExpressionFromAny evaluates the supplied interface argument and
// returns a *StringValueExpression if the supplied argument can be converted
// into a StringValueExpression, or nil if the conversion cannot be done.
func StringValueExpressionFromAny(subject interface{}) *grammar.StringValueExpression {
	switch v := subject.(type) {
	case *grammar.StringValueExpression:
		return v
	case grammar.StringValueExpression:
		return &v
	case *grammar.CharacterValueExpression:
		return &grammar.StringValueExpression{
			Character: v,
		}
	case grammar.CharacterValueExpression:
		return &grammar.StringValueExpression{
			Character: &v,
		}
	case *grammar.CharacterFactor:
		return &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: v,
			},
		}
	case grammar.CharacterFactor:
		return &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &v,
			},
		}
	case *grammar.BlobValueExpression:
		return &grammar.StringValueExpression{
			Blob: v,
		}
	case grammar.BlobValueExpression:
		return &grammar.StringValueExpression{
			Blob: &v,
		}
	case *grammar.BlobFactor:
		return &grammar.StringValueExpression{
			Blob: &grammar.BlobValueExpression{
				Factor: v,
			},
		}
	case grammar.BlobFactor:
		return &grammar.StringValueExpression{
			Blob: &grammar.BlobValueExpression{
				Factor: &v,
			},
		}
	}
	v := CharacterPrimaryFromAny(subject)
	if v != nil {
		return &grammar.StringValueExpression{
			Character: &grammar.CharacterValueExpression{
				Factor: &grammar.CharacterFactor{
					Primary: *v,
				},
			},
		}
	}
	return nil
}

// CharacterValueExpressionFromAny evaluates the supplied interface argument and
// returns a *CharacterValueExpression if the supplied argument can be converted
// into a CharacterValueExpression, or nil if the conversion cannot be done.
func CharacterValueExpressionFromAny(subject interface{}) *grammar.CharacterValueExpression {
	switch v := subject.(type) {
	case *grammar.CharacterValueExpression:
		return v
	case grammar.CharacterValueExpression:
		return &v
	case *grammar.CharacterFactor:
		return &grammar.CharacterValueExpression{
			Factor: v,
		}
	case grammar.CharacterFactor:
		return &grammar.CharacterValueExpression{
			Factor: &v,
		}
	}
	v := CharacterPrimaryFromAny(subject)
	if v != nil {
		return &grammar.CharacterValueExpression{
			Factor: &grammar.CharacterFactor{
				Primary: *v,
			},
		}
	}
	return nil
}

// CharacterPrimaryFromAny evaluates the supplied interface argument and
// returns a *CharacterPrimary if the supplied argument can be converted into a
// CharacterPrimary, or nil if the conversion cannot be done.
func CharacterPrimaryFromAny(subject interface{}) *grammar.CharacterPrimary {
	switch v := subject.(type) {
	case *grammar.CharacterPrimary:
		return v
	case grammar.CharacterPrimary:
		return &v
	case *grammar.StringValueFunction:
		return &grammar.CharacterPrimary{
			Function: v,
		}
	case grammar.StringValueFunction:
		return &grammar.CharacterPrimary{
			Function: &v,
		}
	}
	v := ValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.CharacterPrimary{
			Primary: v,
		}
	}
	return nil
}

// ReferredFromStringValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied StringValueExpression.
func ReferredFromStringValueExpression(
	cve *grammar.StringValueExpression,
) []string {
	return []string{}
}
