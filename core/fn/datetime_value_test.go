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
