package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestC(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    c := m.Table("users").C("name")

    exp := "users.name"
    expLen := len(exp)
    s := c.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnWithTableAlias(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    c := m.Table("users").As("u").C("name")

    exp := "u.name"
    expLen := len(exp)
    s := c.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnAlias(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    c := m.Table("users").C("name").As("user_name")

    exp := "users.name AS user_name"
    expLen := len(exp)
    s := c.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
