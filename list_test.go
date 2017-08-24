package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    cl := &List{elements: []element{c}}

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

    c1 := &Column{
        cdef: colUserId,
        tbl: users.Table(),
    }

    c2:= &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    cl := &List{elements: []element{c1, c2}}

    exp := "users.id, users.name"
    expLen := len(exp)
    s := cl.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := cl.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
