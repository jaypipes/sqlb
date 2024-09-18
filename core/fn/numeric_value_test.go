//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn_test

import (
	"testing"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/core/fn"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNumericValueFunctionCharacterLength(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.LengthExpression
		expRefersTo types.Relation
	}{
		{
			name:    "column",
			subject: colUserId,
			exp: &grammar.LengthExpression{
				Character: &grammar.CharacterLengthExpression{
					Subject: grammar.StringValueExpression{
						Character: &grammar.CharacterValueExpression{
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
					},
				},
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.CharacterLength(tt.subject)
			assert.Equal(tt.exp, got.LengthExpression)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.CharacterLength(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.CharacterLength(users)
	})
}

func TestSelectCharacterLengthFunction(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name:  "char_length literal",
			q:     expr.Select(colUserId, fn.CharacterLength("42")),
			qs:    "SELECT users.id, CHAR_LENGTH(?) FROM users",
			qargs: []interface{}{"42"},
		},
		{
			name: "char_length column",
			q:    expr.Select(fn.CharacterLength(colUserId)),
			qs:   "SELECT CHAR_LENGTH(users.id) FROM users",
		},
		{
			name: "char_length column with using",
			q:    expr.Select(fn.CharacterLength(colUserId).Using(grammar.CharacterLengthUnitsOctets)),
			qs:   "SELECT CHAR_LENGTH(users.id USING OCTETS) FROM users",
		},
		{
			name: "char_length column using alias",
			q:    expr.Select(fn.CharacterLength(colUserId).As("clen")),
			qs:   "SELECT CHAR_LENGTH(users.id) AS clen FROM users",
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
