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
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestStringValueFunctionSubstring(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name    string
		subject interface{}
		from    interface{}
		exp     *api.SubstringFunction
	}{
		{
			name:    "column and unsigned literal",
			subject: colUserId,
			from:    42,
			exp: &api.SubstringFunction{
				SubstringFunction: &grammar.SubstringFunction{
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
								Primary: &grammar.ValueExpressionPrimary{
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
					},
					From: grammar.NumericValueExpression{
						Unary: &grammar.Term{
							Unary: &grammar.Factor{
								Primary: grammar.NumericPrimary{
									Primary: &grammar.ValueExpressionPrimary{
										Primary: &grammar.NonParenthesizedValueExpressionPrimary{
											UnsignedValue: &grammar.UnsignedValueSpecification{
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
						},
					},
				},
				Referred: users,
			},
		},
		{
			name:    "column and another column",
			subject: colUserId,
			from:    colUserId,
			exp: &api.SubstringFunction{
				SubstringFunction: &grammar.SubstringFunction{
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
								Primary: &grammar.ValueExpressionPrimary{
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
					},
					From: grammar.NumericValueExpression{
						Unary: &grammar.Term{
							Unary: &grammar.Factor{
								Primary: grammar.NumericPrimary{
									Primary: &grammar.ValueExpressionPrimary{
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
						},
					},
				},
				Referred: users,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Substring(tt.subject, tt.from)
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.Substring(struct{}{}, 1)
	})
	// First argument must be coercible into a NumericValueExpression
	assert.Panics(t, func() {
		_ = api.Substring(colUserId, struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.Substring(users, 1)
	})
}

func TestSelectSubstringFunction(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		q     *api.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name:  "substring column with literal from",
			q:     api.Select(api.Substring(colUserId, 42)),
			qs:    "SELECT SUBSTRING(users.id FROM ?) FROM users",
			qargs: []interface{}{42},
		},
		{
			name:  "substring column with literal from using alias",
			q:     api.Select(api.Substring(colUserId, 42).As("subber")),
			qs:    "SELECT SUBSTRING(users.id FROM ?) AS subber FROM users",
			qargs: []interface{}{42},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			b := builder.New()

			qs, qargs := b.StringArgs(tt.q.Query())
			assert.Equal(len(tt.qargs), len(qargs))
			assert.Equal(tt.qs, qs)
		})
	}
}
