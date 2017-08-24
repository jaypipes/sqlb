package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

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

    c := &Column{
        cdef: colUserName,
        tbl: users.As("u"),
    }

    exp := "u.name"
    expLen := len(exp)
    s := c.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnDefSorts(t *testing.T) {
    assert := assert.New(t)

    sc := colUserName.Asc()

    exp := "users.name"
    expLen := len(exp)
    s := sc.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := sc.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))

    sc = colUserName.Desc()

    exp = "users.name DESC"
    expLen = len(exp)
    s = sc.size()
    assert.Equal(expLen, s)

    b = make([]byte, s)
    written, _ = sc.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnSorts(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    sc := c.Asc()

    exp := "users.name"
    expLen := len(exp)
    s := sc.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := sc.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))

    sc = c.Desc()

    exp = "users.name DESC"
    expLen = len(exp)
    s = sc.size()
    assert.Equal(expLen, s)

    b = make([]byte, s)
    written, _ = sc.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnAlias(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
        alias: "user_name",
    }

    exp := "users.name AS user_name"
    expLen := len(exp)
    s := c.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnAs(t *testing.T) {
    assert := assert.New(t)

    c := colUserName.As("n")
    assert.Equal("n", c.alias)
    assert.Equal(colUserName, c.cdef)
}
