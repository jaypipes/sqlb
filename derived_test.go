//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type derivedTest struct {
	c     *derivedTable
	qs    string
	qargs []interface{}
}

func TestDerived(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []derivedTest{
		// Simple one-column sub-SELECT
		derivedTest{
			c: &derivedTable{
				from: &selectStatement{
					projs: []projection{
						colUserName,
					},
					selections: []selection{
						users,
					},
				},
				alias: "u",
			},
			qs: "(SELECT users.name FROM users) AS u",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.size(defaultScanner)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.argCount())

		b := make([]byte, s)
		curArg := 0
		written := test.c.scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
