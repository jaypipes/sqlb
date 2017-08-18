package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSelectSingleColumn(t *testing.T) {
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

    sel := Select(c)

    exp := "SELECT name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.Size())
    assert.Equal(exp, sel.String())
}

func TestSelectMultiColumnsSingleTable(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        table: td,
    }

    c1 := &Column{
        def: cd1,
    }

    cd2 := &ColumnDef{
        name: "email",
        table: td,
    }

    c2 := &Column{
        def: cd2,
    }

    sel := Select(c1, c2)

    exp := "SELECT name, email FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.Size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromColumnDef(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        table: td,
    }

    sel := Select(cd)

    exp := "SELECT name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.Size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromColumnDefAndColumn(t *testing.T) {
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

    c2 := &Column{
        def: cd2,
    }

    sel := Select(cd1, c2)

    exp := "SELECT name, email FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.Size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromTableDef(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := []*ColumnDef{
        &ColumnDef{
            name: "name",
            table: td,
        },
        &ColumnDef{
            name: "email",
            table: td,
        },
    }
    td.columns = cdefs

    sel := Select(td)

    exp := "SELECT name, email FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.Size())
    assert.Equal(exp, sel.String())
}

func TestWhereSingleEqual(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        table: td,
    }

    sel := Select(cd).Where(Equal(cd, "foo"))

    exp := "SELECT name FROM users WHERE name = ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}
