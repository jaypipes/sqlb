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
