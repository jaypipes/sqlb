package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type orderByTest struct {
	c     *orderByClause
	qs    string
	qargs []interface{}
}

func TestOrderBy(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserId := users.C("id")
	colUserName := users.C("name")

	tests := []orderByTest{
		// column asc
		orderByTest{
			c: &orderByClause{
				scols: []*sortColumn{colUserName.Asc()},
			},
			qs: " ORDER BY users.name",
		},
		// column desc
		orderByTest{
			c: &orderByClause{
				scols: []*sortColumn{colUserName.Desc()},
			},
			qs: " ORDER BY users.name DESC",
		},
		// Aliased column should NOT output alias in ORDER BY
		orderByTest{
			c: &orderByClause{
				scols: []*sortColumn{colUserName.As("user_name").Desc()},
			},
			qs: " ORDER BY users.name DESC",
		},
		// multi column mixed
		orderByTest{
			c: &orderByClause{
				scols: []*sortColumn{colUserName.Asc(), colUserId.Desc()},
			},
			qs: " ORDER BY users.name, users.id DESC",
		},
		// sort by a function
		orderByTest{
			c: &orderByClause{
				scols: []*sortColumn{Count(users).Desc()},
			},
			qs: " ORDER BY COUNT(*) DESC",
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
