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
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }

    exp := "users.name"
    expLen := len(exp)
    s := c.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnWithTableAlias(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.As("u"),
    }

    exp := "u.name"
    expLen := len(exp)
    s := c.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnDefSorts(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sc := cd.Asc()

    exp := "users.name"
    expLen := len(exp)
    s := sc.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := sc.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))

    sc = cd.Desc()

    exp = "users.name DESC"
    expLen = len(exp)
    s = sc.Size()
    assert.Equal(expLen, s)

    b = make([]byte, s)
    written, _ = sc.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnSorts(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }

    sc := c.Asc()

    exp := "users.name"
    expLen := len(exp)
    s := sc.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := sc.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))

    sc = c.Desc()

    exp = "users.name DESC"
    expLen = len(exp)
    s = sc.Size()
    assert.Equal(expLen, s)

    b = make([]byte, s)
    written, _ = sc.Scan(b, nil)

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
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
        alias: "user_name",
    }

    exp := "users.name AS user_name"
    expLen := len(exp)
    s := c.Size()
    assert.Equal(expLen, s)

    b := make([]byte, s)
    written, _ := c.Scan(b, nil)

    assert.Equal(written, s)
    assert.Equal(exp, string(b))
}

func TestColumnAs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    c := cd.As("n")
    assert.Equal("n", c.alias)
    assert.Equal(cd, c.cdef)
}
