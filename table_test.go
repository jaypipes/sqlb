package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
    assert := assert.New(t)

    exp := "users"
    expLen := len(exp)
    s := users.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := users.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableAlias(t *testing.T) {
    assert := assert.New(t)

    t1 := users.As("u")

    exp := "users AS u"
    expLen := len(exp)
    s := t1.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := t1.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableColumnDefs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
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

    c := users.Column("name")

    assert.Equal(users, c.cdef.tdef)
    assert.Equal("name", c.cdef.name)

    // Check an unknown column name returns nil
    unknown := users.Column("unknown")
    assert.Nil(unknown)
}

func TestTableColumn(t *testing.T) {
    assert := assert.New(t)

    tbl := &Table{
        tdef: users,
    }

    c := tbl.Column("name")

    assert.Equal(users, c.cdef.tdef)
    assert.Equal("name", c.cdef.name)

    // Check an unknown column name returns nil
    unknown := tbl.Column("unknown")
    assert.Nil(unknown)
}
