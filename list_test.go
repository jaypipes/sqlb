package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestListSingle(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }

    cl := &List{elements: []Element{c}}

    exp := "name"
    expLen := len(exp)
    s := cl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := cl.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestListMulti(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        tdef: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        tdef: td,
    }

    c1 := &Column{
        cdef: cd1,
        tbl: td.Table(),
    }

    c2:= &Column{
        cdef: cd2,
        tbl: td.Table(),
    }

    cl := &List{elements: []Element{c1, c2}}

    exp := "name, email"
    expLen := len(exp)
    s := cl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := cl.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
