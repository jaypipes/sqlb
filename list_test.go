package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    colUserName := users.Column("name")

    cl := &List{elements: []element{colUserName}}

    exp := "users.name"
    expLen := len(exp)
    s := cl.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := cl.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestListMulti(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    colUserId := users.Column("id")
    colUserName := users.Column("name")

    cl := &List{elements: []element{colUserId, colUserName}}

    exp := "users.id, users.name"
    expLen := len(exp)
    s := cl.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := cl.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
