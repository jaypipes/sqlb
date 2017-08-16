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

func TestTableListSingle(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := &Table{
        def: td,
    }

    tl := &TableList{tables: []*Table{t1}}

    exp := "users"
    expLen := len(exp)
    s := tl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := tl.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableListMulti(t *testing.T) {
    assert := assert.New(t)

    td1 := &TableDef{
        name: "users",
        schema: "test",
    }

    td2 := &TableDef{
        name: "articles",
        schema: "test",
    }

    t1 := &Table{
        def: td1,
    }

    t2 := &Table{
        def: td2,
    }

    tl := &TableList{tables: []*Table{t1, t2}}

    exp := "users, articles"
    expLen := len(exp)
    s := tl.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written := tl.Scan(b)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableColumnDefs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := map[string]*ColumnDef{
        "id": &ColumnDef{
            name: "id",
            table: td,
        },
        "email": &ColumnDef{
            name: "email",
            table: td,
        },
    }
    td.columns = cdefs

    defs := td.ColumnDefs()

    assert.Equal(2, len(defs))
    for _, def := range defs {
        assert.Equal(td, def.table)
    }
}

func TestTableColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cdefs := map[string]*ColumnDef{
        "id": &ColumnDef{
            name: "id",
            table: td,
        },
        "email": &ColumnDef{
            name: "email",
            table: td,
        },
    }
    td.columns = cdefs

    defs := td.ColumnDefs()

    assert.Equal(2, len(defs))
    for _, def := range defs {
        assert.Equal(td, def.table)
    }
}

func TestTableAs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    t1 := td.As("u")
    assert.Equal("u", t1.alias)
    assert.Equal(td, t1.def)
}
