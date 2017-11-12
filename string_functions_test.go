package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringFunctions(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		el    element
		qs    map[Dialect]string
		qargs []interface{}
	}{
		{
			name: "TRIM(column)",
			el:   Trim(colUserName),
			qs: map[Dialect]string{
				DIALECT_MYSQL:      "TRIM(users.name)",
				DIALECT_POSTGRESQL: "BTRIM(users.name)",
			},
		},
	}
	for _, test := range tests {
		expArgc := len(test.qargs)
		argc := test.el.argCount()
		assert.Equal(expArgc, argc)

		// Test each SQL dialect output
		for dialect, qs := range test.qs {
			test.el.setDialect(dialect)
			expLen := len(qs)
			size := test.el.size()
			size += interpolationLength(dialect, argc)
			assert.Equal(expLen, size)

			b := make([]byte, size)
			curArg := 0
			written := test.el.scan(b, test.qargs, &curArg)

			assert.Equal(written, size)
			assert.Equal(qs, string(b))
		}
	}
}
