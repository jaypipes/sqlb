package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAsTemplated(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    assert.Equal("", c.alias)

    As(c, "n")

    assert.Equal("n", c.alias)

    c = c.As("name")

    assert.Equal("name", c.alias)
}

func TestAsMethod(t *testing.T) {
    assert := assert.New(t)

    t1 := &Table{
        tdef: users,
    }

    assert.Equal("", t1.alias)

    t1 = t1.As("t")

    assert.Equal("t", t1.alias)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    assert.Equal("", c.alias)

    c = c.As("n")

    assert.Equal("n", c.alias)
}
