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
	// Second argument must be coercible into a NumericValueExpression
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

func TestStringValueFunctionRegexSubstring(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name    string
		subject interface{}
		similar interface{}
		escape  interface{}
		exp     *api.RegexSubstringFunction
	}{
		{
			name:    "column and two string literals",
			subject: colUserId,
			similar: "$[a-z]",
			escape:  "/",
			exp: &api.RegexSubstringFunction{
				RegexSubstringFunction: &grammar.RegexSubstringFunction{
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
												GeneralLiteral: &grammar.GeneralLiteral{
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
												GeneralLiteral: &grammar.GeneralLiteral{
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
				Referred: users,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.RegexSubstring(tt.subject, tt.similar, tt.escape)
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.RegexSubstring(struct{}{}, "2", "3")
	})
	// Second argument must be coercible into a CharacerValueExpression
	assert.Panics(t, func() {
		_ = api.RegexSubstring(colUserId, struct{}{}, "3")
	})
	// Third argument must be coercible into a CharacerValueExpression
	assert.Panics(t, func() {
		_ = api.RegexSubstring(colUserId, "2", struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.RegexSubstring(users, "1", "2")
	})
}

func TestSelectRegexSubstringFunction(t *testing.T) {
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
			name:  "substring column with string literals",
			q:     api.Select(api.RegexSubstring(colUserId, "$[a-z]", "/")),
			qs:    "SELECT SUBSTRING(users.id SIMILAR ? ESCAPE ?) FROM users",
			qargs: []interface{}{"$[a-z]", "/"},
		},
		{
			name:  "substring column with string literals using alias",
			q:     api.Select(api.RegexSubstring(colUserId, "$[a-z]", "/").As("subber")),
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
		name     string
		subject  interface{}
		foldCase grammar.FoldCase
		exp      *api.FoldFunction
	}{
		{
			name:     "UPPER column",
			subject:  colUserId,
			foldCase: grammar.FoldCaseUpper,
			exp: &api.FoldFunction{
				FoldFunction: &grammar.FoldFunction{
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
				Referred: users,
			},
		},
		{
			name:     "LOWER column",
			subject:  colUserId,
			foldCase: grammar.FoldCaseLower,
			exp: &api.FoldFunction{
				FoldFunction: &grammar.FoldFunction{
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
				Referred: users,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Fold(tt.subject, tt.foldCase)
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.Fold(struct{}{}, grammar.FoldCaseUpper)
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.Fold(users, grammar.FoldCaseUpper)
	})
}

func TestSelectFoldFunction(t *testing.T) {
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
			name: "upper column",
			q:    api.Select(api.Upper(colUserId)),
			qs:   "SELECT UPPER(users.id) FROM users",
		},
		{
			name: "lower column with alias",
			q:    api.Select(api.Lower(colUserId).As("lowered")),
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
		exp             *api.TranscodingFunction
	}{
		{
			name:            "CONVERT column",
			subject:         colUserId,
			transcodingName: "utf8",
			exp: &api.TranscodingFunction{
				TranscodingFunction: &grammar.TranscodingFunction{
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
				Referred: users,
			},
		},
		{
			name:            "CONVERT string literal",
			subject:         "foo",
			transcodingName: "utf8",
			exp: &api.TranscodingFunction{
				TranscodingFunction: &grammar.TranscodingFunction{
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
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
					Using: grammar.SchemaQualifiedName{
						Identifiers: grammar.IdentifierChain{
							Identifiers: []string{"utf8"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Convert(tt.subject, tt.transcodingName)
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.Convert(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.Convert(users, "utf8")
	})
}

func TestSelectTranscodingFunction(t *testing.T) {
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
			name: "convert column",
			q:    api.Select(api.Convert(colUserId, "utf8")),
			qs:   "SELECT CONVERT(users.id USING utf8) FROM users",
		},
		{
			name: "convert column with alias",
			q:    api.Select(api.Convert(colUserId, "utf8").As("converted")),
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

func TestStringValueFunctionTransliteration(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name            string
		subject         interface{}
		transcodingName string
		exp             *api.TransliterationFunction
	}{
		{
			name:            "TRANSLATE column",
			subject:         colUserId,
			transcodingName: "utf8",
			exp: &api.TransliterationFunction{
				TransliterationFunction: &grammar.TransliterationFunction{
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
				Referred: users,
			},
		},
		{
			name:            "TRANSLATE string literal",
			subject:         "foo",
			transcodingName: "utf8",
			exp: &api.TransliterationFunction{
				TransliterationFunction: &grammar.TransliterationFunction{
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
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
					Using: grammar.SchemaQualifiedName{
						Identifiers: grammar.IdentifierChain{
							Identifiers: []string{"utf8"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Translate(tt.subject, tt.transcodingName)
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.Convert(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.Convert(users, "utf8")
	})
}

func TestSelectTransliterationFunction(t *testing.T) {
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
			name: "convert column",
			q:    api.Select(api.Translate(colUserId, "utf8")),
			qs:   "SELECT TRANSLATE(users.id USING utf8) FROM users",
		},
		{
			name: "convert column with alias",
			q:    api.Select(api.Translate(colUserId, "utf8").As("translated")),
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
		name    string
		subject interface{}
		chars   interface{}
		spec    grammar.TrimSpecification
		exp     *api.TrimFunction
	}{
		{
			name:    "TRIM column with literal chars",
			subject: colUserId,
			chars:   "\n",
			exp: &api.TrimFunction{
				TrimFunction: &grammar.TrimFunction{
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
												GeneralLiteral: &grammar.GeneralLiteral{
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
				Referred: users,
			},
		},
		{
			name:    "TRIM string literal no characters",
			subject: "foo",
			exp: &api.TrimFunction{
				TrimFunction: &grammar.TrimFunction{
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
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
				},
			},
		},
		{
			name:    "TRIM string literal with literal chars trailing only",
			subject: "foo",
			chars:   "\n",
			spec:    grammar.TrimSpecificationTrailing,
			exp: &api.TrimFunction{
				TrimFunction: &grammar.TrimFunction{
					Specification: grammar.TrimSpecificationTrailing,
					Subject: grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
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
					Character: &grammar.CharacterValueExpression{
						Factor: &grammar.CharacterFactor{
							Primary: grammar.CharacterPrimary{
								Primary: &grammar.ValueExpressionPrimary{
									Primary: &grammar.NonParenthesizedValueExpressionPrimary{
										UnsignedValue: &grammar.UnsignedValueSpecification{
											UnsignedLiteral: &grammar.UnsignedLiteral{
												GeneralLiteral: &grammar.GeneralLiteral{
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			var got *api.TrimFunction
			if tt.chars != nil {
				got = api.Trim(tt.subject, tt.chars, tt.spec)
			} else {
				got = api.TrimSpace(tt.subject)
			}
			assert.Equal(tt.exp, got)
		})
	}

	// First argument must be coercible into a CharacterValueExpression
	assert.Panics(t, func() {
		_ = api.Convert(struct{}{}, "utf8")
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a CharacterValueExpression
		_ = api.Convert(users, "utf8")
	})
}

func TestSelectTrimFunction(t *testing.T) {
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
			name: "trim space column",
			q:    api.Select(api.TrimSpace(colUserId)),
			qs:   "SELECT TRIM(users.id) FROM users",
		},
		{
			name: "leading trim space column",
			q:    api.Select(api.LTrimSpace(colUserId)),
			qs:   "SELECT TRIM(LEADING users.id) FROM users",
		},
		{
			name: "trailing trim space column",
			q:    api.Select(api.RTrimSpace(colUserId)),
			qs:   "SELECT TRIM(TRAILING users.id) FROM users",
		},
		{
			name:  "trim column",
			q:     api.Select(api.Trim(colUserId, "\n", grammar.TrimSpecificationBoth)),
			qs:    "SELECT TRIM(? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name:  "leading trim column",
			q:     api.Select(api.LTrim(colUserId, "\n")),
			qs:    "SELECT TRIM(LEADING ? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name:  "trailing trim column",
			q:     api.Select(api.RTrim(colUserId, "\n")),
			qs:    "SELECT TRIM(TRAILING ? FROM users.id) FROM users",
			qargs: []interface{}{"\n"},
		},
		{
			name: "trim column with alias",
			q:    api.Select(api.TrimSpace(colUserId).As("trimmed")),
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
