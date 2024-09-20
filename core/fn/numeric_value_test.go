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
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/types"
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

func TestNumericValueFunctionOctetLength(t *testing.T) {
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
				Octet: &grammar.OctetLengthExpression{
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
			got := fn.OctetLength(tt.subject)
			assert.Equal(tt.exp, got.LengthExpression)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a OctetValueExpression
	assert.Panics(t, func() {
		_ = fn.OctetLength(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a OctetValueExpression
		_ = fn.OctetLength(users)
	})
}

func TestSelectOctetLengthFunction(t *testing.T) {
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
			name:  "octet_length literal",
			q:     expr.Select(colUserId, fn.OctetLength("42")),
			qs:    "SELECT users.id, OCTET_LENGTH(?) FROM users",
			qargs: []interface{}{"42"},
		},
		{
			name: "octet_length column",
			q:    expr.Select(fn.OctetLength(colUserId)),
			qs:   "SELECT OCTET_LENGTH(users.id) FROM users",
		},
		{
			name: "octet_length column using alias",
			q:    expr.Select(fn.OctetLength(colUserId).As("clen")),
			qs:   "SELECT OCTET_LENGTH(users.id) AS clen FROM users",
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

func TestNumericValueFunctionPosition(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		in          interface{}
		exp         *grammar.PositionExpression
		expRefersTo types.Relation
	}{
		{
			name:    "literal and column",
			subject: "foo",
			in:      colUserId,
			exp: &grammar.PositionExpression{
				Blob: &grammar.BlobPositionExpression{
					Subject: grammar.BlobValueExpression{
						Factor: &grammar.BlobFactor{
							Primary: grammar.BlobPrimary{
								Primary: &grammar.ValueExpressionPrimary{
									Primary: &grammar.NonParenthesizedValueExpressionPrimary{
										UnsignedValue: &grammar.UnsignedValueSpecification{
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
					},
					In: grammar.BlobValueExpression{
						Factor: &grammar.BlobFactor{
							Primary: grammar.BlobPrimary{
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
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Position(tt.subject, tt.in)
			assert.Equal(tt.exp, got.PositionExpression)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a BlobValueExpression or StringValueExpression
	assert.Panics(t, func() {
		_ = fn.Position(struct{}{}, colUserId)
	})
	// Second argument must be coercible into a BlobValueExpression or StringValueExpression
	assert.Panics(t, func() {
		_ = fn.Substring(colUserId, struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a StringValueExpression
		_ = fn.Position(users, colUserId)
	})
}

func TestSelectPositionFunction(t *testing.T) {
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
			name:  "position literal in column",
			q:     expr.Select(fn.Position("42", colUserId)),
			qs:    "SELECT POSITION(? IN users.id) FROM users",
			qargs: []interface{}{"42"},
		},
		{
			name:  "position column in literal",
			q:     expr.Select(fn.Position(colUserId, "42")),
			qs:    "SELECT POSITION(users.id IN ?) FROM users",
			qargs: []interface{}{"42"},
		},
		{
			name:  "position column with using",
			q:     expr.Select(fn.Position("42", colUserId).Using(grammar.CharacterLengthUnitsOctets)),
			qs:    "SELECT POSITION(? IN users.id USING OCTETS) FROM users",
			qargs: []interface{}{"42"},
		},
		{
			name:  "position column using alias",
			q:     expr.Select(fn.Position("42", colUserId).As("pos")),
			qs:    "SELECT POSITION(? IN users.id) AS pos FROM users",
			qargs: []interface{}{"42"},
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

func TestNumericValueFunctionExtract(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		from        interface{}
		what        fn.ExtractField
		exp         *grammar.ExtractExpression
		expRefersTo types.Relation
	}{
		{
			name: "extract second from column",
			from: colUserId,
			what: fn.ExtractFieldSecond,
			exp: &grammar.ExtractExpression{
				From: grammar.ExtractSource{
					Datetime: &grammar.DatetimeValueExpression{
						Unary: &grammar.DatetimeTerm{
							Factor: grammar.DatetimeFactor{
								Primary: grammar.DatetimePrimary{
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
				What: grammar.ExtractField{
					Datetime: &grammar.PrimaryDatetimeField{
						Second: true,
					},
				},
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Extract(tt.from, tt.what)
			assert.Equal(tt.exp, got.ExtractExpression)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectExtractFunction(t *testing.T) {
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
			name: "extract second from column",
			q:    expr.Select(fn.Extract(colUserId, fn.ExtractFieldSecond)),
			qs:   "SELECT EXTRACT(SECOND FROM users.id) FROM users",
		},
		{
			name: "extract month from column",
			q:    expr.Select(fn.Extract(colUserId, fn.ExtractFieldMonth)),
			qs:   "SELECT EXTRACT(MONTH FROM users.id) FROM users",
		},
		{
			name: "extract timezone minute from column using alias",
			q:    expr.Select(fn.Extract(colUserId, fn.ExtractFieldTimezoneMinute).As("tzm")),
			qs:   "SELECT EXTRACT(TIMEZONE_MINUTE FROM users.id) AS tzm FROM users",
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

func TestNumericValueFunctionNaturalLogarithm(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "natural logarithm column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				Natural: &grammar.NaturalLogarithm{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.NaturalLogarithm(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectNaturalLogarithm(t *testing.T) {
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
			name: "natural logarithm on column",
			q:    expr.Select(fn.NaturalLogarithm(colUserId)),
			qs:   "SELECT LN(users.id) FROM users",
		},
		{
			name: "ln on column",
			q:    expr.Select(fn.Ln(colUserId)),
			qs:   "SELECT LN(users.id) FROM users",
		},
		{
			name: "natural logarithm on column using alias",
			q:    expr.Select(fn.Ln(colUserId).As("nlog")),
			qs:   "SELECT LN(users.id) AS nlog FROM users",
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

func TestNumericValueFunctionAbsolute(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "absolute value column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				AbsoluteValue: &grammar.AbsoluteValueExpression{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Absolute(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectAbsolute(t *testing.T) {
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
			name: "absolute value on column",
			q:    expr.Select(fn.Absolute(colUserId)),
			qs:   "SELECT ABS(users.id) FROM users",
		},
		{
			name: "abs on column",
			q:    expr.Select(fn.Abs(colUserId)),
			qs:   "SELECT ABS(users.id) FROM users",
		},
		{
			name: "absolute value on column using alias",
			q:    expr.Select(fn.Abs(colUserId).As("abber")),
			qs:   "SELECT ABS(users.id) AS abber FROM users",
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

func TestNumericValueFunctionExponential(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "exponential column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				Exponential: &grammar.ExponentialFunction{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Exponential(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectExponential(t *testing.T) {
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
			name: "exponential function on column",
			q:    expr.Select(fn.Exponential(colUserId)),
			qs:   "SELECT EXP(users.id) FROM users",
		},
		{
			name: "exp on column",
			q:    expr.Select(fn.Exp(colUserId)),
			qs:   "SELECT EXP(users.id) FROM users",
		},
		{
			name: "exponential function on column using alias",
			q:    expr.Select(fn.Exp(colUserId).As("exper")),
			qs:   "SELECT EXP(users.id) AS exper FROM users",
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

func TestNumericValueFunctionSquareRoot(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "square root column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				SquareRoot: &grammar.SquareRoot{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.SquareRoot(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectSquareRoot(t *testing.T) {
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
			name: "square root function on column",
			q:    expr.Select(fn.SquareRoot(colUserId)),
			qs:   "SELECT SQRT(users.id) FROM users",
		},
		{
			name: "exp on column",
			q:    expr.Select(fn.SquareRoot(colUserId)),
			qs:   "SELECT SQRT(users.id) FROM users",
		},
		{
			name: "square root function on column using alias",
			q:    expr.Select(fn.SqRt(colUserId).As("rooter")),
			qs:   "SELECT SQRT(users.id) AS rooter FROM users",
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

func TestNumericValueFunctionCeiling(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "ceiling column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				Ceiling: &grammar.CeilingFunction{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Ceiling(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectCeiling(t *testing.T) {
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
			name: "ceiling function on column",
			q:    expr.Select(fn.Ceiling(colUserId)),
			qs:   "SELECT CEIL(users.id) FROM users",
		},
		{
			name: "exp on column",
			q:    expr.Select(fn.Ceil(colUserId)),
			qs:   "SELECT CEIL(users.id) FROM users",
		},
		{
			name: "ceiling function on column using alias",
			q:    expr.Select(fn.Ceil(colUserId).As("ceiler")),
			qs:   "SELECT CEIL(users.id) AS ceiler FROM users",
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

func TestNumericValueFunctionFloor(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		exp         *grammar.NumericValueFunction
		expRefersTo types.Relation
	}{
		{
			name:    "ceiling column",
			subject: colUserId,
			exp: &grammar.NumericValueFunction{
				Floor: &grammar.FloorFunction{
					Subject: grammar.NumericValueExpression{
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Floor(tt.subject)
			assert.Equal(tt.exp, got.NumericValueFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectFloor(t *testing.T) {
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
			name: "floor function on column",
			q:    expr.Select(fn.Floor(colUserId)),
			qs:   "SELECT FLOOR(users.id) FROM users",
		},
		{
			name: "floor function on column using alias",
			q:    expr.Select(fn.Floor(colUserId).As("floorer")),
			qs:   "SELECT FLOOR(users.id) AS floorer FROM users",
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
