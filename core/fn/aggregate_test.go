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

func TestAggregateFunctionCount(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		target      []interface{}
		exp         *grammar.AggregateFunction
		expRefersTo types.Relation
	}{
		{
			name: "count star",
			exp: &grammar.AggregateFunction{
				CountStar: &struct{}{},
			},
		},
		{
			name:   "column",
			target: []interface{}{colUserId},
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: grammar.ComputationalOperationCount,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.Count(tt.target...)
			assert.Equal(tt.exp, got.AggregateFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// Count expects either zero or one argument
	assert.Panics(t, func() {
		_ = fn.Count(1, 2)
	})

	// Argument must be coercible into a ValueExpression
	assert.Panics(t, func() {
		_ = fn.Count(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a ValueExpression
		_ = fn.Count(users)
	})
}

func TestAggregateFunctionAvgMinMaxSum(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		target      interface{}
		f           func(interface{}) *fn.AggregateFunction
		exp         *grammar.AggregateFunction
		expRefersTo types.Relation
	}{
		{
			name:   "avg column",
			target: colUserId,
			f:      fn.Avg,
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: grammar.ComputationalOperationAvg,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
		{
			name:   "min column",
			target: colUserId,
			f:      fn.Min,
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: grammar.ComputationalOperationMin,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
		{
			name:   "max column",
			target: colUserId,
			f:      fn.Max,
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: grammar.ComputationalOperationMax,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
		{
			name:   "sum column",
			target: colUserId,
			f:      fn.Sum,
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation: grammar.ComputationalOperationSum,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.f(tt.target)
			assert.Equal(tt.exp, got.AggregateFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}

	// Argument must be coercible into a ValueExpression
	assert.Panics(t, func() {
		_ = fn.Max(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a ValueExpression
		_ = fn.Max(users)
	})
}

func TestAggregateFunctionDistinct(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name        string
		target      interface{}
		f           func(interface{}) *fn.AggregateFunction
		exp         *grammar.AggregateFunction
		expRefersTo types.Relation
	}{
		{
			name:   "avg distinct column",
			target: colUserId,
			f:      fn.Avg,
			exp: &grammar.AggregateFunction{
				GeneralSet: &grammar.GeneralSetFunction{
					Operation:  grammar.ComputationalOperationAvg,
					Quantifier: grammar.SetQuantifierDistinct,
					Value: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
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
			expRefersTo: users,
		},
		{
			name: "count star is not affected with distinct",
			f:    func(_ interface{}) *fn.AggregateFunction { return fn.Count() },
			exp: &grammar.AggregateFunction{
				CountStar: &struct{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.f(tt.target).Distinct()
			assert.Equal(tt.exp, got.AggregateFunction)
			assert.Equal(tt.expRefersTo, got.References())
		})
	}
}

func TestSelectAggregateFunction(t *testing.T) {
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
			name: "count star on Selection",
			q:    expr.Select(users).Count(),
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name: "count star on Table",
			q:    expr.Select(users.Count()),
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name: "count star expr.Count",
			q:    expr.Select(users, fn.Count()),
			qs:   "SELECT users.created_on, users.id, users.name, COUNT(*) FROM users",
		},
		{
			name: "count distinct column",
			q:    expr.Select(fn.Count(colUserId).Distinct()),
			qs:   "SELECT COUNT(DISTINCT users.id) FROM users",
		},
		{
			name: "count column",
			q:    expr.Select(fn.Count(colUserId)),
			qs:   "SELECT COUNT(users.id) FROM users",
		},
		{
			name: "avg column",
			q:    expr.Select(fn.Avg(colUserId)),
			qs:   "SELECT AVG(users.id) FROM users",
		},
		{
			name: "min column",
			q:    expr.Select(fn.Min(colUserId)),
			qs:   "SELECT MIN(users.id) FROM users",
		},
		{
			name: "max column",
			q:    expr.Select(fn.Max(colUserId)),
			qs:   "SELECT MAX(users.id) FROM users",
		},
		{
			name: "sum column",
			q:    expr.Select(fn.Sum(colUserId)),
			qs:   "SELECT SUM(users.id) FROM users",
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
