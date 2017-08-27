package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAsFunc(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    c := m.Table("users").Column("name")

    assert.Equal("", c.alias)
    As(c, "n")
    assert.Equal("n", c.alias)
    c = c.As("name")
    assert.Equal("name", c.alias)
}

func TestAsMethod(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")

    assert.Equal("", users.alias)
    u := users.As("u")
    assert.Equal("u", u.alias)
    c := m.Table("users").Column("name")
    assert.Equal("", c.alias)
    c = c.As("n")
    assert.Equal("n", c.alias)
}
