package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
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

    exp := "name"
    expLen := len(exp)
    s := c.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := c.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnAlias(t *testing.T) {
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
        alias: "user_name",
    }

    exp := "name AS user_name"
    expLen := len(exp)
    s := c.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := c.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnListSingle(t *testing.T) {
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

    cl := &ColumnList{columns: []*Column{c}}

    exp := "name"
    expLen := len(exp)
    s := cl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := cl.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnListMulti(t *testing.T) {
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

    cl := &ColumnList{columns: []*Column{c1, c2}}

    exp := "name, email"
    expLen := len(exp)
    s := cl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := cl.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
