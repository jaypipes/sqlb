package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestAsTemplated(t *testing.T) {
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

    assert.Equal("", c.alias)

    As(c, "n")

    assert.Equal("n", c.alias)

    c = c.As("name")

    assert.Equal("name", c.alias)
}

func TestAsMethod(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := &Table{
        def: td,
    }

    assert.Equal("", t1.alias)

    t1 = t1.As("t")

    assert.Equal("t", t1.alias)

    cd := &ColumnDef{
        name: "name",
        table: td,
    }

    c := &Column{
        def: cd,
    }

    assert.Equal("", c.alias)

    c = c.As("n")

    assert.Equal("n", c.alias)
}
