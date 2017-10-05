package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestTableMeta(t *testing.T) {
    assert := assert.New(t)

    m := NewMeta("mysql", "test")
    td := m.Table("users")
    assert.Nil(td)
    td = m.NewTable("users")
    assert.NotNil(td)
    assert.Equal(td.meta, m)

    assert.Equal(td, m.Table("users"))

    cd := td.C("id")
    assert.Nil(cd)

    cd = td.NewColumn("id")
    assert.NotNil(cd)

    assert.Equal(cd, td.C("id"))
}

func TestTable(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")

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

    m := testFixtureMeta()
    u := m.Table("users").As("u")

    exp := "users AS u"
    expLen := len(exp)
    s := u.size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := u.scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestTableColumns(t *testing.T) {
    assert := assert.New(t)

    td := &Table{
        name: "users",
    }

    cols := []*Column{
         &Column{
            name: "id",
            tbl: td,
        },
        &Column{
            name: "email",
            tbl: td,
        },
    }
    td.columns = cols

    defs := td.columns

    assert.Equal(2, len(defs))
    for _, def := range defs {
        assert.Equal(td, def.tbl)
    }

    // Check stable order of insertion from above...
    assert.Equal(defs[0].name, "id")
    assert.Equal(defs[1].name, "email")
}

func TestTableC(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")

    c := users.C("name")

    assert.Equal(users, c.tbl)
    assert.Equal("name", c.name)

    // Check an unknown column name returns nil
    unknown := users.C("unknown")
    assert.Nil(unknown)
}
