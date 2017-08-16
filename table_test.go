package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := &Table{
        def: td,
    }

    exp := "users"
    expLen := len(exp)
    s := t1.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := t1.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableAlias(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := &Table{
        def: td,
        alias: "u",
    }

    exp := "users AS u"
    expLen := len(exp)
    s := t1.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := t1.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}
