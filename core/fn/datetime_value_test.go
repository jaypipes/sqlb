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
	"github.com/jaypipes/sqlb/grammar"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDatetimeValueFunctionCurrentDate(t *testing.T) {
	tests := []struct {
		name string
		exp  *grammar.DatetimeValueFunction
	}{
		{
			name: "no args",
			exp: &grammar.DatetimeValueFunction{
				CurrentDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.CurrentDate()
			assert.Equal(tt.exp, got.DatetimeValueFunction)
		})
	}
}

func TestSelectCurrentDateFunction(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")
	colCreatedOn := users.C("created_on")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name: "CURRENT_DATE func no args",
			q: expr.Select(colUserId).Where(
				expr.Equal(colCreatedOn, fn.CurrentDate()),
			),
			qs: "SELECT users.id FROM users WHERE users.created_on = CURRENT_DATE()",
		},
		{
			name: "CURRENT_DATE func using alias",
			q:    expr.Select(colUserId, fn.CurrentDate().As("now")),
			qs:   "SELECT users.id, CURRENT_DATE() AS now FROM users",
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

func TestDatetimeValueFunctionCurrentTime(t *testing.T) {
	p10 := uint(10)
	tests := []struct {
		name      string
		exp       *grammar.DatetimeValueFunction
		precision *uint
	}{
		{
			name: "no args",
			exp: &grammar.DatetimeValueFunction{
				CurrentTime: &grammar.CurrentTimeFunction{},
			},
		},
		{
			name: "with precision",
			exp: &grammar.DatetimeValueFunction{
				CurrentTime: &grammar.CurrentTimeFunction{
					Precision: &p10,
				},
			},
			precision: &p10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.CurrentTime()
			if tt.precision != nil {
				got = got.Precision(*tt.precision)
			}
			assert.Equal(tt.exp, got.DatetimeValueFunction)
		})
	}
}

func TestSelectCurrentTimeFunction(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")
	colCreatedOn := users.C("created_on")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name: "CURRENT_TIME func no args",
			q: expr.Select(colUserId).Where(
				expr.Equal(colCreatedOn, fn.CurrentTime()),
			),
			qs: "SELECT users.id FROM users WHERE users.created_on = CURRENT_TIME()",
		},
		{
			name: "CURRENT_TIME func precision arg",
			q: expr.Select(colUserId).Where(
				expr.Equal(colCreatedOn, fn.CurrentTime().Precision(10)),
			),
			qs: "SELECT users.id FROM users WHERE users.created_on = CURRENT_TIME(10)",
		},
		{
			name: "CURRENT_TIME func using alias",
			q:    expr.Select(colUserId, fn.CurrentTime().As("now")),
			qs:   "SELECT users.id, CURRENT_TIME() AS now FROM users",
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

func TestDatetimeValueFunctionCurrentTimestamp(t *testing.T) {
	p10 := uint(10)
	tests := []struct {
		name      string
		exp       *grammar.DatetimeValueFunction
		precision *uint
	}{
		{
			name: "no args",
			exp: &grammar.DatetimeValueFunction{
				CurrentTimestamp: &grammar.CurrentTimestampFunction{},
			},
		},
		{
			name: "with precision",
			exp: &grammar.DatetimeValueFunction{
				CurrentTimestamp: &grammar.CurrentTimestampFunction{
					Precision: &p10,
				},
			},
			precision: &p10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got := fn.CurrentTimestamp()
			if tt.precision != nil {
				got = got.Precision(*tt.precision)
			}
			assert.Equal(tt.exp, got.DatetimeValueFunction)
		})
	}
}

func TestSelectCurrentTimestampFunction(t *testing.T) {
	m := testutil.M()
	users := m.T("users")
	colUserId := users.C("id")
	colCreatedOn := users.C("created_on")

	tests := []struct {
		name  string
		q     *expr.Selection
		qs    string
		qargs []interface{}
	}{
		{
			name: "CURRENT_TIMESTAMP func no args",
			q: expr.Select(colUserId).Where(
				expr.Equal(colCreatedOn, fn.CurrentTimestamp()),
			),
			qs: "SELECT users.id FROM users WHERE users.created_on = CURRENT_TIMESTAMP()",
		},
		{
			name: "CURRENT_TIMESTAMP func precision arg",
			q: expr.Select(colUserId).Where(
				expr.Equal(colCreatedOn, fn.CurrentTimestamp().Precision(10)),
			),
			qs: "SELECT users.id FROM users WHERE users.created_on = CURRENT_TIMESTAMP(10)",
		},
		{
			name: "CURRENT_TIMESTAMP func using alias",
			q:    expr.Select(colUserId, fn.CurrentTimestamp().As("now")),
			qs:   "SELECT users.id, CURRENT_TIMESTAMP() AS now FROM users",
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
