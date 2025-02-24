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

func TestStringValueFunctionCharacterSubstring(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		from        interface{}
		exp         *grammar.CharacterSubstringFunction
		expRefersTo types.Relation
	}{
		{
			name:    "column and unsigned literal",
			subject: colUserId,
			from:    42,
			exp: &grammar.CharacterSubstringFunction{
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
												UnsignedNumeric: &grammar.UnsignedNumericLiteral{
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
			expRefersTo: users,
		},
		{
			name:    "column and another column",
			subject: colUserId,
			from:    colUserId,
			exp: &grammar.CharacterSubstringFunction{
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
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Substring(tt.subject, tt.from)
			assert.Equal(tt.exp, got.CharacterSubstringFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Substring(struct{}{}, 1)
	})
	// Second argument must be coercible into a NumericValueExpression
	assert.Panics(t, func() {
		_ = fn.Substring(colUserId, struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.Substring(users, 1)
	})
}

func TestSelectSubstringFunction(t *testing.T) {
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
			name:  "substring column with literal from",
			q:     expr.Select(fn.Substring(colUserId, 42)),
			qs:    "SELECT SUBSTRING(users.id FROM ?) FROM users",
			qargs: []interface{}{42},
		},
		{
			name:  "substring column with literal from using alias",
			q:     expr.Select(fn.Substring(colUserId, 42).As("subber")),
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

func TestStringValueFunctionRegexSubstring(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		similar     interface{}
		escape      interface{}
		exp         *grammar.RegexSubstringFunction
		expRefersTo types.Relation
	}{
		{
			name:    "column and two string literals",
			subject: colUserId,
			similar: "$[a-z]",
			escape:  "/",
			exp: &grammar.RegexSubstringFunction{
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
				Similar: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
								Primary: &grammar.NonParenthesizedValueExpressionPrimary{
									UnsignedValue: &grammar.UnsignedValueSpecification{
										UnsignedLiteral: &grammar.UnsignedLiteral{
											General: &grammar.GeneralLiteral{
												Value: "$[a-z]",
											},
										},
									},
								},
							},
						},
					},
				},
				Escape: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
								Primary: &grammar.NonParenthesizedValueExpressionPrimary{
									UnsignedValue: &grammar.UnsignedValueSpecification{
										UnsignedLiteral: &grammar.UnsignedLiteral{
											General: &grammar.GeneralLiteral{
												Value: "/",
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
			got := fn.RegexSubstring(tt.subject, tt.similar, tt.escape)
			assert.Equal(tt.exp, got.RegexSubstringFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.RegexSubstring(struct{}{}, "2", "3")
	})
	// Second argument must be coercible into a CharacerValueExpression
	assert.Panics(t, func() {
		_ = fn.RegexSubstring(colUserId, struct{}{}, "3")
	})
	// Third argument must be coercible into a CharacerValueExpression
	assert.Panics(t, func() {
		_ = fn.RegexSubstring(colUserId, "2", struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.RegexSubstring(users, "1", "2")
	})
}

func TestSelectRegexSubstringFunction(t *testing.T) {
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
			name:  "substring column with string literals",
			q:     expr.Select(fn.RegexSubstring(colUserId, "$[a-z]", "/")),
			qs:    "SELECT SUBSTRING(users.id SIMILAR ? ESCAPE ?) FROM users",
			qargs: []interface{}{"$[a-z]", "/"},
		},
		{
			name:  "substring column with string literals using alias",
			q:     expr.Select(fn.RegexSubstring(colUserId, "$[a-z]", "/").As("subber")),
			qs:    "SELECT SUBSTRING(users.id SIMILAR ? ESCAPE ?) AS subber FROM users",
			qargs: []interface{}{"$[a-z]", "/"},
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

func TestStringValueFunctionFold(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		foldCase    grammar.FoldCase
		exp         *grammar.FoldFunction
		expRefersTo types.Relation
	}{
		{
			name:     "UPPER column",
			subject:  colUserId,
			foldCase: grammar.FoldCaseUpper,
			exp: &grammar.FoldFunction{
				Case: grammar.FoldCaseUpper,
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
			},
			expRefersTo: users,
		},
		{
			name:     "LOWER column",
			subject:  colUserId,
			foldCase: grammar.FoldCaseLower,
			exp: &grammar.FoldFunction{
				Case: grammar.FoldCaseLower,
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
			},
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Fold(tt.subject, tt.foldCase)
			assert.Equal(tt.exp, got.FoldFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Fold(struct{}{}, grammar.FoldCaseUpper)
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.Fold(users, grammar.FoldCaseUpper)
	})
}

func TestSelectFoldFunction(t *testing.T) {
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
			name: "upper column",
			q:    expr.Select(fn.Upper(colUserId)),
			qs:   "SELECT UPPER(users.id) FROM users",
		},
		{
			name: "lower column with alias",
			q:    expr.Select(fn.Lower(colUserId).As("lowered")),
			qs:   "SELECT LOWER(users.id) AS lowered FROM users",
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

func TestStringValueFunctionTranscoding(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name            string
		subject         interface{}
		transcodingName string
		exp             *grammar.TranscodingFunction
		expRefersTo     types.Relation
	}{
		{
			name:            "CONVERT column",
			subject:         colUserId,
			transcodingName: "utf8",
			exp: &grammar.TranscodingFunction{
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
				Using: grammar.SchemaQualifiedName{
					Identifiers: grammar.IdentifierChain{
						Identifiers: []string{"utf8"},
					},
				},
			},
			expRefersTo: users,
		},
		{
			name:            "CONVERT string literal",
			subject:         "foo",
			transcodingName: "utf8",
			exp: &grammar.TranscodingFunction{
				Subject: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
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
						},
					},
				},
				Using: grammar.SchemaQualifiedName{
					Identifiers: grammar.IdentifierChain{
						Identifiers: []string{"utf8"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Convert(tt.subject, tt.transcodingName)
			assert.Equal(tt.exp, got.TranscodingFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Convert(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.Convert(users, "utf8")
	})
}

func TestSelectTranscodingFunction(t *testing.T) {
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
			name: "convert column",
			q:    expr.Select(fn.Convert(colUserId, "utf8")),
			qs:   "SELECT CONVERT(users.id USING utf8) FROM users",
		},
		{
			name: "convert column with alias",
			q:    expr.Select(fn.Convert(colUserId, "utf8").As("converted")),
			qs:   "SELECT CONVERT(users.id USING utf8) AS converted FROM users",
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

func TestStringValueFunctionCharacterTransliteration(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name            string
		subject         interface{}
		transcodingName string
		exp             *grammar.CharacterTransliterationFunction
		expRefersTo     types.Relation
	}{
		{
			name:            "TRANSLATE column",
			subject:         colUserId,
			transcodingName: "utf8",
			exp: &grammar.CharacterTransliterationFunction{
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
				Using: grammar.SchemaQualifiedName{
					Identifiers: grammar.IdentifierChain{
						Identifiers: []string{"utf8"},
					},
				},
			},
			expRefersTo: users,
		},
		{
			name:            "TRANSLATE string literal",
			subject:         "foo",
			transcodingName: "utf8",
			exp: &grammar.CharacterTransliterationFunction{
				Subject: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
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
						},
					},
				},
				Using: grammar.SchemaQualifiedName{
					Identifiers: grammar.IdentifierChain{
						Identifiers: []string{"utf8"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Translate(tt.subject, tt.transcodingName)
			assert.Equal(tt.exp, got.CharacterTransliterationFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Translate(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.Translate(users, "utf8")
	})
}

func TestSelectTransliterationFunction(t *testing.T) {
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
			name: "convert column",
			q:    expr.Select(fn.Translate(colUserId, "utf8")),
			qs:   "SELECT TRANSLATE(users.id USING utf8) FROM users",
		},
		{
			name: "convert column with alias",
			q:    expr.Select(fn.Translate(colUserId, "utf8").As("translated")),
			qs:   "SELECT TRANSLATE(users.id USING utf8) AS translated FROM users",
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

func TestStringValueFunctionTrim(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		chars       interface{}
		spec        grammar.TrimSpecification
		exp         *grammar.TrimFunction
		expRefersTo types.Relation
	}{
		{
			name:    "TRIM column with literal chars",
			subject: colUserId,
			chars:   "\n",
			exp: &grammar.TrimFunction{
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
				Character: &grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
								Primary: &grammar.NonParenthesizedValueExpressionPrimary{
									UnsignedValue: &grammar.UnsignedValueSpecification{
										UnsignedLiteral: &grammar.UnsignedLiteral{
											General: &grammar.GeneralLiteral{
												Value: "\n",
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
		{
			name:    "TRIM string literal no characters",
			subject: "foo",
			exp: &grammar.TrimFunction{
				Subject: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
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
						},
					},
				},
			},
		},
		{
			name:    "TRIM string literal with literal chars trailing only",
			subject: "foo",
			chars:   "\n",
			spec:    grammar.TrimSpecificationTrailing,
			exp: &grammar.TrimFunction{
				Specification: grammar.TrimSpecificationTrailing,
				Subject: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
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
						},
					},
				},
				Character: &grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
								Primary: &grammar.NonParenthesizedValueExpressionPrimary{
									UnsignedValue: &grammar.UnsignedValueSpecification{
										UnsignedLiteral: &grammar.UnsignedLiteral{
											General: &grammar.GeneralLiteral{
												Value: "\n",
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
			var got *fn.TrimFunction
			if tt.chars != nil {
				got = fn.Trim(tt.subject, tt.chars, tt.spec)
			} else {
				got = fn.TrimSpace(tt.subject)
			}
			assert.Equal(tt.exp, got.TrimFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Convert(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = fn.Convert(users, "utf8")
	})
}

func TestSelectTrimFunction(t *testing.T) {
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
			name: "trim space column",
			q:    expr.Select(fn.TrimSpace(colUserId)),
			qs:   "SELECT TRIM(users.id) FROM users",
		},
		{
			name: "leading trim space column",
			q:    expr.Select(fn.LTrimSpace(colUserId)),
			qs:   "SELECT TRIM(LEADING users.id) FROM users",
		},
		{
			name: "trailing trim space column",
			q:    expr.Select(fn.RTrimSpace(colUserId)),
			qs:   "SELECT TRIM(TRAILING users.id) FROM users",
		},
		{
			name:  "trim column",
			q:     expr.Select(fn.Trim(colUserId, "\n", grammar.TrimSpecificationBoth)),
			qs:    "SELECT TRIM(? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name:  "leading trim column",
			q:     expr.Select(fn.LTrim(colUserId, "\n")),
			qs:    "SELECT TRIM(LEADING ? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name:  "trailing trim column",
			q:     expr.Select(fn.RTrim(colUserId, "\n")),
			qs:    "SELECT TRIM(TRAILING ? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name: "trim column with alias",
			q:    expr.Select(fn.TrimSpace(colUserId).As("trimmed")),
			qs:   "SELECT TRIM(users.id) AS trimmed FROM users",
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

func TestStringValueFunctionCharacterOverlay(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		subject     interface{}
		placing     interface{}
		from        interface{}
		exp         *grammar.CharacterOverlayFunction
		expRefersTo types.Relation
	}{
		{
			name:    "column, character literal and unsigned literal",
			subject: colUserId,
			placing: "replace",
			from:    1,
			exp: &grammar.CharacterOverlayFunction{
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
				Placing: grammar.CharacterValueExpression{
					Factor: &grammar.CharacterFactor{
						Primary: grammar.CharacterPrimary{
							Primary: &grammar.ValueExpressionPrimary{
								Primary: &grammar.NonParenthesizedValueExpressionPrimary{
									UnsignedValue: &grammar.UnsignedValueSpecification{
										UnsignedLiteral: &grammar.UnsignedLiteral{
											General: &grammar.GeneralLiteral{
												Value: "replace",
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
												UnsignedNumeric: &grammar.UnsignedNumericLiteral{
													Value: 1,
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
		{
			name:    "column, another column and another column",
			subject: colUserId,
			placing: colUserId,
			from:    colUserId,
			exp: &grammar.CharacterOverlayFunction{
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
				Placing: grammar.CharacterValueExpression{
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
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Overlay(tt.subject, tt.placing, tt.from)
			assert.Equal(tt.exp, got.CharacterOverlayFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Overlay(struct{}{}, "replace", 1)
	})
	// Second argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = fn.Overlay(colUserId, struct{}{}, 1)
	})
	// Third argument must be coercible into a NumericValueExpression
	assert.Panics(t, func() {
		_ = fn.Overlay(users, "replace", struct{}{})
	})
}

func TestSelectOverlayFunction(t *testing.T) {
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
			name:  "overlay column with literal from",
			q:     expr.Select(fn.Overlay(colUserId, "replace", 42)),
			qs:    "SELECT OVERLAY(users.id PLACING ? FROM ?) FROM users",
			qargs: []interface{}{"replace", 42},
		},
		{
			name:  "overlay column with literal from using alias",
			q:     expr.Select(fn.Overlay(colUserId, "replace", 42).As("overlayer")),
			qs:    "SELECT OVERLAY(users.id PLACING ? FROM ?) AS overlayer FROM users",
			qargs: []interface{}{"replace", 42},
		},
		{
			name:  "overlay column placing column with literal from",
			q:     expr.Select(fn.Overlay(colUserId, colUserId, 42)),
			qs:    "SELECT OVERLAY(users.id PLACING users.id FROM ?) FROM users",
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
