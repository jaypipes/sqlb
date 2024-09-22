//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr_test

import (
	"testing"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	articles := m.T("articles")
	colUserId := users.C("id")
	colArticleAuthor := articles.C("author")

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		exp  *grammar.ComparisonPredicate
	}{
		{
			name: "two columns",
			a:    colUserId,
			b:    colArticleAuthor,
			exp: &grammar.ComparisonPredicate{
				A: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
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
				B: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
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
			},
		},
		{
			name: "column and numeric literal",
			a:    colUserId,
			b:    1,
			exp: &grammar.ComparisonPredicate{
				A: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
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
				B: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValue: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								UnsignedNumeric: &grammar.UnsignedNumericLiteral{
									Value: 1,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "string literal and column",
			a:    "foo",
			b:    colUserId,
			exp: &grammar.ComparisonPredicate{
				A: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
						UnsignedValue: &grammar.UnsignedValueSpecification{
							UnsignedLiteral: &grammar.UnsignedLiteral{
								General: &grammar.GeneralLiteral{
									Value: "foo",
								},
							},
						},
					},
				},
				B: grammar.RowValuePredicand{
					Primary: &grammar.NonParenthesizedValueExpressionPrimary{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			got := expr.Equal(tt.a, tt.b)
			assert.Equal(tt.exp, got)
		})
	}
}
