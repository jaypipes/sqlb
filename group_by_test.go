//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/types"
	"github.com/stretchr/testify/assert"
)

type groupByClauseTest struct {
	c     *groupByClause
	qs    string
	qargs []interface{}
}

func TestGroupByClause(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []groupByClauseTest{
		// Single column
		groupByClauseTest{
			c: &groupByClause{
				cols: []types.Projection{colUserName},
			},
			qs: " GROUP BY users.name",
		},
		// Multiple columns
		groupByClauseTest{
			c: &groupByClause{
				cols: []types.Projection{colUserName, colUserId},
			},
			qs: " GROUP BY users.name, users.id",
		},
		// Aliased column should NOT output alias in GROUP BY
		groupByClauseTest{
			c: &groupByClause{
				cols: []types.Projection{colUserName.As("user_name")},
			},
			qs: " GROUP BY users.name",
		},
	}
	for _, test := range tests {
		expLen := len(test.qs)
		s := test.c.Size(defaultScanner)
		assert.Equal(expLen, s)

		expArgc := len(test.qargs)
		assert.Equal(expArgc, test.c.ArgCount())

		b := make([]byte, s)
		curArg := 0
		written := test.c.Scan(defaultScanner, b, test.qargs, &curArg)

		assert.Equal(written, s)
		assert.Equal(test.qs, string(b))
	}
}
