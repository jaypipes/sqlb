//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package function_test

import (
	"strings"
	"testing"

	"github.com/jaypipes/sqlb"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/jaypipes/sqlb/types"
	"github.com/stretchr/testify/assert"
)

func TestTrimFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := sqlb.T(m, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		el    scanner.Element
		qs    map[types.Dialect]string
		qargs []interface{}
	}{
		{
			name: "TRIM(column) or BTRIM(column)",
			el:   function.Trim(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "TRIM(users.name)",
				types.DialectPostgreSQL: "BTRIM(users.name)",
			},
		},
		{
			name: "LTRIM(column) or TRIM(LEADING FROM column)",
			el:   function.LTrim(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "LTRIM(users.name)",
				types.DialectPostgreSQL: "TRIM(LEADING FROM users.name)",
			},
		},
		{
			name: "RTRIM(column) or TRIM(TRAILING FROM column)",
			el:   function.RTrim(colUserName),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "RTRIM(users.name)",
				types.DialectPostgreSQL: "TRIM(TRAILING FROM users.name)",
			},
		},
		{
			name: "TRIM(remstr FROM column) OR BTRIM(column, chars)",
			el:   function.TrimChars(colUserName, "xyz"),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "TRIM(? FROM users.name)",
				types.DialectPostgreSQL: "BTRIM(users.name, $1)",
			},
			qargs: []interface{}{"xyz"},
		},
		{
			name: "TRIM(LEADING remstr FROM column)",
			el:   function.LTrimChars(colUserName, "xyz"),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "TRIM(LEADING ? FROM users.name)",
				types.DialectPostgreSQL: "TRIM(LEADING $1 FROM users.name)",
			},
			qargs: []interface{}{"xyz"},
		},
		{
			name: "TRIM(TRAILING remstr FROM column)",
			el:   function.RTrimChars(colUserName, "xyz"),
			qs: map[types.Dialect]string{
				types.DialectMySQL:      "TRIM(TRAILING ? FROM users.name)",
				types.DialectPostgreSQL: "TRIM(TRAILING $1 FROM users.name)",
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
			sc := &scanner.Scanner{
				Dialect: dialect,
			}
			expLen := len(qs)
			size := test.el.Size(sc)
			size += scanner.InterpolationLength(dialect, argc)
			assert.Equal(expLen, size)

			var b strings.Builder
			b.Grow(size)
			args := make([]interface{}, argc)
			curArg := 0
			test.el.Scan(sc, &b, args, &curArg)

			assert.Equal(qs, b.String())
			if expArgc > 0 {
				assert.Equal(args, test.qargs)
			}
		}
	}
}
