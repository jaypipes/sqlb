package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.TableDef("users")

    sel := Select(users)
    q := Query{sel: sel}

    qs, qargs := q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users", qs)
    assert.Nil(qargs)

    // Modify the underlying SELECT and verify string and args changed
    q.Where(Equal(users.Column("id"), 1))
    qs, qargs = q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
    assert.Equal([]interface{}{1}, qargs)
}
