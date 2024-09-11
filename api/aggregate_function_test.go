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

func TestAggregateFunctionCount(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name   string
		target []interface{}
		exp    *api.AggregateFunction
	}{
		{
			name: "count star",
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					CountStar: &struct{}{},
				},
			},
		},
		{
			name:   "column",
			target: []interface{}{colUserId},
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation: grammar.ComputationalOperationCount,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := api.Count(tt.target...)
			assert.Equal(tt.exp, got)
		})
	}

	// Count expects either zero or one argument
	assert.Panics(t, func() {
		_ = api.Count(1, 2)
	})

	// Argument must be coercible into a ValueExpression
	assert.Panics(t, func() {
		_ = api.Count(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a ValueExpression
		_ = api.Count(users)
	})
}

func TestAggregateFunctionAvgMinMaxSum(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name   string
		target interface{}
		f      func(interface{}) *api.AggregateFunction
		exp    *api.AggregateFunction
	}{
		{
			name:   "avg column",
			target: colUserId,
			f:      api.Avg,
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation: grammar.ComputationalOperationAvg,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
		{
			name:   "min column",
			target: colUserId,
			f:      api.Min,
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation: grammar.ComputationalOperationMin,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
		{
			name:   "max column",
			target: colUserId,
			f:      api.Max,
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation: grammar.ComputationalOperationMax,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
		{
			name:   "sum column",
			target: colUserId,
			f:      api.Sum,
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation: grammar.ComputationalOperationSum,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.f(tt.target)
			assert.Equal(tt.exp, got)
		})
	}

	// Count expects either zero or one argument
	assert.Panics(t, func() {
		_ = api.Count(1, 2)
	})

	// Argument must be coercible into a ValueExpression
	assert.Panics(t, func() {
		_ = api.Count(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a ValueExpression
		_ = api.Count(users)
	})
}

func TestAggregateFunctionDistinct(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")

	tests := []struct {
		name   string
		target interface{}
		f      func(interface{}) *api.AggregateFunction
		exp    *api.AggregateFunction
	}{
		{
			name:   "avg distinct column",
			target: colUserId,
			f:      api.Avg,
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					GeneralSetFunction: &grammar.GeneralSetFunction{
						Operation:  grammar.ComputationalOperationAvg,
						Quantifier: grammar.SetQuantifierDistinct,
						ValueExpression: grammar.ValueExpression{
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
				Referred: users,
			},
		},
		{
			name: "count star is not affected with distinct",
			f:    func(_ interface{}) *api.AggregateFunction { return api.Count() },
			exp: &api.AggregateFunction{
				AggregateFunction: &grammar.AggregateFunction{
					CountStar: &struct{}{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := tt.f(tt.target).Distinct()
			assert.Equal(tt.exp, got)
		})
	}

	// Count expects either zero or one argument
	assert.Panics(t, func() {
		_ = api.Count(1, 2)
	})

	// Argument must be coercible into a ValueExpression
	assert.Panics(t, func() {
		_ = api.Count(struct{}{})
	})
	assert.Panics(t, func() {
		// A Table is not coercible into a ValueExpression
		_ = api.Count(users)
	})
}

func TestSelectAggregateFunction(t *testing.T) {
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
			name: "count star on Selection",
			q:    api.Select(users).Count(),
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name: "count star on Table",
			q:    api.Select(users.Count()),
			qs:   "SELECT COUNT(*) FROM users",
		},
		{
			name: "count star api.Count",
			q:    api.Select(users, api.Count()),
			qs:   "SELECT users.id, users.name, COUNT(*) FROM users",
		},
		{
			name: "count distinct column",
			q:    api.Select(api.Count(colUserId).Distinct()),
			qs:   "SELECT COUNT(DISTINCT users.id) FROM users",
		},
		{
			name: "count column",
			q:    api.Select(api.Count(colUserId)),
			qs:   "SELECT COUNT(users.id) FROM users",
		},
		{
			name: "avg column",
			q:    api.Select(api.Avg(colUserId)),
			qs:   "SELECT AVG(users.id) FROM users",
		},
		{
			name: "min column",
			q:    api.Select(api.Min(colUserId)),
			qs:   "SELECT MIN(users.id) FROM users",
		},
		{
			name: "max column",
			q:    api.Select(api.Max(colUserId)),
			qs:   "SELECT MAX(users.id) FROM users",
		},
		{
			name: "sum column",
			q:    api.Select(api.Sum(colUserId)),
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
