package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type queryTest struct {
    q *Query
    qs string
    qargs []interface{}
}

func TestQuery(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    colUserName := users.Column("name")

    tests := []queryTest{
        // Simple FROM
        queryTest{
            q: Select(users),
            qs: "SELECT users.id, users.name FROM users",
        },
        // Simple WHERE
        queryTest{
            q: Select(users).Where(Equal(colUserName, "foo")),
            qs: "SELECT users.id, users.name FROM users WHERE users.name = ?",
            qargs: []interface{}{"foo"},
        },
        // Simple GROUP BY
        queryTest{
            q: Select(users).GroupBy(colUserName),
            qs: "SELECT users.id, users.name FROM users GROUP BY users.name",
        },
        // Simple ORDER BY
        queryTest{
            q: Select(users).OrderBy(colUserName.Desc()),
            qs: "SELECT users.id, users.name FROM users ORDER BY users.name DESC",
        },
        // Simple LIMIT
        queryTest{
            q: Select(users).Limit(10),
            qs: "SELECT users.id, users.name FROM users LIMIT ?",
            qargs: []interface{}{10},
        },
        // Simple LIMIT with OFFSET
        queryTest{
            q: Select(users).LimitWithOffset(10, 20),
            qs: "SELECT users.id, users.name FROM users LIMIT ? OFFSET ?",
            qargs: []interface{}{10, 20},
        },
    }
    for _, test := range tests {
        qs, qargs := test.q.StringArgs()
        assert.Equal(len(test.qargs), len(qargs))
        assert.Equal(test.qs, qs)
    }
}

func TestModifyingQueryUpdatesBuffer(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.TableDef("users")

    q := Select(users)

    qs, qargs := q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users", qs)
    assert.Nil(qargs)

    // Modify the underlying SELECT and verify string and args changed
    q.Where(Equal(users.Column("id"), 1))
    qs, qargs = q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
    assert.Equal([]interface{}{1}, qargs)
}
