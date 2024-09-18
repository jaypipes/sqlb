//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package inspect

import "github.com/jaypipes/sqlb/grammar"

// BlobValueExpressionFromAny evaluates the supplied interface argument and
// returns a *BlobValueExpression if the supplied argument can be converted
// into a BlobValueExpression, or nil if the conversion cannot be done.
func BlobValueExpressionFromAny(subject interface{}) *grammar.BlobValueExpression {
	switch v := subject.(type) {
	case *grammar.BlobValueExpression:
		return v
	case grammar.BlobValueExpression:
		return &v
	case *grammar.BlobFactor:
		return &grammar.BlobValueExpression{
			Factor: v,
		}
	case grammar.BlobFactor:
		return &grammar.BlobValueExpression{
			Factor: &v,
		}
	}
	v := BlobPrimaryFromAny(subject)
	if v != nil {
		return &grammar.BlobValueExpression{
			Factor: &grammar.BlobFactor{
				Primary: *v,
			},
		}
	}
	return nil
}

// BlobPrimaryFromAny evaluates the supplied interface argument and
// returns a *BlobPrimary if the supplied argument can be converted into a
// BlobPrimary, or nil if the conversion cannot be done.
func BlobPrimaryFromAny(subject interface{}) *grammar.BlobPrimary {
	switch v := subject.(type) {
	case *grammar.BlobPrimary:
		return v
	case grammar.BlobPrimary:
		return &v
	case *grammar.StringValueFunction:
		return &grammar.BlobPrimary{
			Function: v,
		}
	case grammar.StringValueFunction:
		return &grammar.BlobPrimary{
			Function: &v,
		}
	}
	v := ValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.BlobPrimary{
			Primary: v,
		}
	}
	return nil
}

// ReferredFromBlobValueExpression returns a slice of string names of any
// tables or derived tables (subqueries in the FROM clause) that are referenced
// within a supplied BlobValueExpression.
func ReferredFromBlobValueExpression(
	cve *grammar.BlobValueExpression,
) []string {
	return []string{}
}
