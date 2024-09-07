//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api_test

import (
	"testing"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPredicateEqual(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	scalarSubquery := api.Select(users).Where(api.Equal(colUserId, 42)).Query()

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		exp  *grammar.ComparisonPredicate
	}{
		{
			name: "two columns",
			a:    colArticleAuthor,
			b:    colUserId,
			exp: &grammar.ComparisonPredicate{
				Operator: grammar.ComparisonOperatorEquals,
				A: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"users",
									"id",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "column and numeric literal value",
			a:    colArticleAuthor,
			b:    42,
			exp: &grammar.ComparisonPredicate{
				Operator: grammar.ComparisonOperatorEquals,
				A: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 42,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "aliased column and numeric literal value",
			a:    colArticleAuthor.As("a"),
			b:    42,
			exp: &grammar.ComparisonPredicate{
				Operator: grammar.ComparisonOperatorEquals,
				A: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									// The alias should not change the
									// underlying identifier chain's
									// composition. An alias for a column
									// should only add a "AS <alias>" to the
									// *projection* in the SELECT list and
									// nowhere else.
									"author",
								},
							},
							Correlation: &grammar.Correlation{
								Name: "a",
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 42,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "column and string literal value",
			a:    colArticleAuthor,
			b:    "foo",
			exp: &grammar.ComparisonPredicate{
				Operator: grammar.ComparisonOperatorEquals,
				A: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								GeneralLiteral: &grammar.GeneralLiteral{
									Value: "foo",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "column and scalar subquery",
			a:    colArticleAuthor,
			b:    scalarSubquery,
			exp: &grammar.ComparisonPredicate{
				Operator: grammar.ComparisonOperatorEquals,
				A: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ScalarSubquery: &grammar.Subquery{
							QueryExpression: grammar.QueryExpression{
								Body: grammar.QueryExpressionBody{
									NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
										NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
											Primary: &grammar.NonJoinQueryPrimary{
												SimpleTable: &grammar.SimpleTable{
													QuerySpecification: scalarSubquery.(*grammar.QuerySpecification),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Equal(tt.a, tt.b)
			assert.Equal(tt.exp, got)
		})
	}
}

func TestPredicateBetween(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	scalarSubquery := api.Select(users).Where(api.Equal(colUserId, 42)).Query()

	tests := []struct {
		name   string
		target interface{}
		start  interface{}
		end    interface{}
		exp    *grammar.BetweenPredicate
	}{
		{
			name:   "column and two numeric literals",
			target: colArticleAuthor,
			start:  1,
			end:    42,
			exp: &grammar.BetweenPredicate{
				Target: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				Start: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 1,
								},
							},
						},
					},
				},
				End: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 42,
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "two columns and numeric literal",
			target: colArticleAuthor,
			start:  colUserId,
			end:    42,
			exp: &grammar.BetweenPredicate{
				Target: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				Start: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"users",
									"id",
								},
							},
						},
					},
				},
				End: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 42,
								},
							},
						},
					},
				},
			},
		},
		{
			name:   "column, numeric literal and scalar subquery",
			target: colArticleAuthor,
			start:  42,
			end:    scalarSubquery,
			exp: &grammar.BetweenPredicate{
				Target: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ColumnReference: &grammar.ColumnReference{
							BasicIdentifierChain: &grammar.IdentifierChain{
								Identifiers: []string{
									"articles",
									"author",
								},
							},
						},
					},
				},
				Start: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValueSpecification: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumericLiteral: &grammar.UnsignedNumericLiteral{
									Value: 42,
								},
							},
						},
					},
				},
				End: grammar.RowValuePredicand{
					NonParenthesizedValueExpressionPrimary: &grammar.NonParenthesizedValueExpressionPrimary{
						ScalarSubquery: &grammar.Subquery{
							QueryExpression: grammar.QueryExpression{
								Body: grammar.QueryExpressionBody{
									NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
										NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
											Primary: &grammar.NonJoinQueryPrimary{
												SimpleTable: &grammar.SimpleTable{
													QuerySpecification: scalarSubquery.(*grammar.QuerySpecification),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Between(tt.target, tt.start, tt.end)
			assert.Equal(tt.exp, got)
		})
	}
}
