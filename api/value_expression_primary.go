//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"slices"

	"github.com/jaypipes/sqlb/grammar"
)

// ValueExpressionPrimaryFromAny evaluates the supplied interface argument and
// returns a *ValueExpressionPrimary if the supplied argument can be converted
// into a ValueExpressionPrimary, or nil if the conversion cannot be done.
func ValueExpressionPrimaryFromAny(
	subject interface{},
) *grammar.ValueExpressionPrimary {
	switch v := subject.(type) {
	case *grammar.ValueExpressionPrimary:
		return v
	case grammar.ValueExpressionPrimary:
		return &v
	}
	v := NonParenthesizedValueExpressionPrimaryFromAny(subject)
	if v != nil {
		return &grammar.ValueExpressionPrimary{
			Primary: v,
		}
	}
	return nil
}

// NonParenthesizedValueExpressionPrimaryFromAny evaluates the supplied
// interface argument and returns a *NonParenthesizedValueExpressionPrimary if
// the supplied argument can be converted into a
// NonParenthesizedValueExpressionPrimary, or nil if the conversion cannot be
// done.
func NonParenthesizedValueExpressionPrimaryFromAny(
	subject interface{},
) *grammar.NonParenthesizedValueExpressionPrimary {
	switch v := subject.(type) {
	case *grammar.NonParenthesizedValueExpressionPrimary:
		return v
	case grammar.NonParenthesizedValueExpressionPrimary:
		return &v
	case *grammar.UnsignedValueSpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			UnsignedValue: v,
		}
	case grammar.UnsignedValueSpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			UnsignedValue: &v,
		}
	case *grammar.ColumnReference:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ColumnReference: v,
		}
	case grammar.ColumnReference:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ColumnReference: &v,
		}
	case *Column:
		tname := v.TableName()
		cr := &grammar.ColumnReference{
			BasicIdentifierChain: &grammar.IdentifierChain{
				Identifiers: []string{
					tname, v.name,
				},
			},
		}
		if v.alias != "" {
			cr.Correlation = &grammar.Correlation{
				Name: v.alias,
			}
		}
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ColumnReference: cr,
		}
	case *grammar.SetFunctionSpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			SetFunction: v,
		}
	case grammar.SetFunctionSpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			SetFunction: &v,
		}
	case *grammar.Subquery:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: v,
		}
	case grammar.Subquery:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &v,
		}
	case *grammar.QueryExpression:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: *v,
			},
		}
	case grammar.QueryExpression:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: v,
			},
		}
	case *grammar.QueryExpressionBody:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: *v,
				},
			},
		}
	case grammar.QueryExpressionBody:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: v,
				},
			},
		}
	case *grammar.QuerySpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: grammar.QueryExpressionBody{
						NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
							NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									SimpleTable: &grammar.SimpleTable{
										QuerySpecification: v,
									},
								},
							},
						},
					},
				},
			},
		}
	case grammar.QuerySpecification:
		return &grammar.NonParenthesizedValueExpressionPrimary{
			ScalarSubquery: &grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: grammar.QueryExpressionBody{
						NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
							NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									SimpleTable: &grammar.SimpleTable{
										QuerySpecification: &v,
									},
								},
							},
						},
					},
				},
			},
		}
	}
	// We could have a simple literal passed to us. See if we can convert it
	// into a simple row value predicand with a non-parenthesized value
	// expression primary
	v := UnsignedValueSpecificationFromAny(subject)
	if v != nil {
		return &grammar.NonParenthesizedValueExpressionPrimary{
			UnsignedValue: v,
		}
	}
	return nil
}

// ReferredFromNonParenthesizedValueExpressionPrimary returns a slice of string
// names of any tables or derived tables (subqueries in the FROM clause) that
// are referenced within a supplied NonParenthesizedValueExpressionPrimary.
func ReferredFromNonParenthesizedValueExpressionPrimary(
	p *grammar.NonParenthesizedValueExpressionPrimary,
) []string {
	if p.ColumnReference != nil {
		if len(p.ColumnReference.BasicIdentifierChain.Identifiers) > 0 {
			return slices.Clone(p.ColumnReference.BasicIdentifierChain.Identifiers[:len(p.ColumnReference.BasicIdentifierChain.Identifiers)-1])
		}
	}
	return []string{}
}
