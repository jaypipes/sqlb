//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function_test

import (
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTrimFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		el    api.Element
		qs    map[api.Dialect]string
		qargs []interface{}
	}{
		{
			name: "TRIM(column) or BTRIM(column)",
			el:   function.Trim(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "TRIM(users.name)",
				api.DialectPostgreSQL: "BTRIM(users.name)",
			},
		},
		{
			name: "LTRIM(column) or TRIM(LEADING FROM column)",
			el:   function.LTrim(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "LTRIM(users.name)",
				api.DialectPostgreSQL: "TRIM(LEADING FROM users.name)",
			},
		},
		{
			name: "RTRIM(column) or TRIM(TRAILING FROM column)",
			el:   function.RTrim(colUserName),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "RTRIM(users.name)",
				api.DialectPostgreSQL: "TRIM(TRAILING FROM users.name)",
			},
		},
		{
			name: "TRIM(remstr FROM column) OR BTRIM(column, chars)",
			el:   function.TrimChars(colUserName, "xyz"),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "TRIM(? FROM users.name)",
				api.DialectPostgreSQL: "BTRIM(users.name, $1)",
			},
			qargs: []interface{}{"xyz"},
		},
		{
			name: "TRIM(LEADING remstr FROM column)",
			el:   function.LTrimChars(colUserName, "xyz"),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "TRIM(LEADING ? FROM users.name)",
				api.DialectPostgreSQL: "TRIM(LEADING $1 FROM users.name)",
			},
			qargs: []interface{}{"xyz"},
		},
		{
			name: "TRIM(TRAILING remstr FROM column)",
			el:   function.RTrimChars(colUserName, "xyz"),
			qs: map[api.Dialect]string{
				api.DialectMySQL:      "TRIM(TRAILING ? FROM users.name)",
				api.DialectPostgreSQL: "TRIM(TRAILING $1 FROM users.name)",
			},
			qargs: []interface{}{"xyz"},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.el.ArgCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			b := builder.New(api.WithDialect(dialect))

			gotqs, gotargs := b.StringArgs(test.el)

			assert.Equal(qs, gotqs)
			if len(test.qargs) > 0 {
				assert.Equal(test.qargs, gotargs)
			}
		}
	}
}
