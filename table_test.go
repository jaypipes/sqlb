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
        tdef: td,
    }

    exp := "users"
    expLen := len(exp)
    s := t1.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := t1.Scan(b, nil)

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
        tdef: td,
        alias: "u",
    }

    exp := "users AS u"
    expLen := len(exp)
    s := t1.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := t1.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableColumnDefs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := []*ColumnDef{
         &ColumnDef{
            name: "id",
            tdef: td,
        },
        &ColumnDef{
            name: "email",
            tdef: td,
        },
    }
    td.cdefs = cdefs

    defs := td.cdefs

    assert.Equal(2, len(defs))
    for _, def := range defs {
        assert.Equal(td, def.tdef)
    }

    // Check stable order of insertion from above...
    assert.Equal(defs[0].name, "id")
    assert.Equal(defs[1].name, "email")
}

func TestTableDefColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := []*ColumnDef{
        &ColumnDef{
            name: "id",
            tdef: td,
        },
        &ColumnDef{
            name: "email",
            tdef: td,
        },
    }
    td.cdefs = cdefs

    c := td.Column("email")

    assert.Equal(td, c.cdef.tdef)
    assert.Equal("email", c.cdef.name)

    // Check an unknown column name returns nil
    unknown := td.Column("unknown")
    assert.Nil(unknown)
}

func TestTableColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := []*ColumnDef{
        &ColumnDef{
            name: "id",
            tdef: td,
        },
        &ColumnDef{
            name: "email",
            tdef: td,
        },
    }
    td.cdefs = cdefs

    tbl := &Table{
        tdef: td,
    }

    c := tbl.Column("email")

    assert.Equal(td, c.cdef.tdef)
    assert.Equal("email", c.cdef.name)

    // Check an unknown column name returns nil
    unknown := tbl.Column("unknown")
    assert.Nil(unknown)
}

func TestTableAs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := td.As("u")
    assert.Equal("u", t1.alias)
    assert.Equal(td, t1.tdef)
}
