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
        table: td,
    }

    c := &Column{
        def: cd,
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
        table: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        table: td,
    }

    c1 := &Column{
        def: cd1,
    }

    c2:= &Column{
        def: cd2,
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
