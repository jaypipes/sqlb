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
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
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
        tdef: td,
    }

    c1 := &Column{
        cdef: cd1,
        tbl: td.Table(),
    }

    cd2 := &ColumnDef{
        name: "email",
        tdef: td,
    }

    c2 := &Column{
        cdef: cd2,
        tbl: td.Table(),
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
        tdef: td,
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
        tdef: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        tdef: td,
    }

    c2 := &Column{
        cdef: cd2,
        tbl: td.Table(),
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
            tdef: td,
        },
        &ColumnDef{
            name: "email",
            tdef: td,
        },
    }
    td.cdefs = cdefs

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
        tdef: td,
    }

    sel := Select(cd).Where(Equal(cd, "foo"))

    exp := "SELECT name FROM users WHERE name = ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestWhereSingleAnd(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).Where(And(NotEqual(cd, "foo"), NotEqual(cd, "bar")))

    exp := "SELECT name FROM users WHERE name != ? AND name != ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestWhereSingleIn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).Where(In(cd, "foo", "bar"))

    exp := "SELECT name FROM users WHERE name IN (?, ?)"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestWhereMultiNotEqual(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).Where(NotEqual(cd, "foo"))
    sel = sel.Where(NotEqual(cd, "bar"))

    exp := "SELECT name FROM users WHERE name != ? AND name != ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestWhereMultiInAndEqual(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        tdef: td,
    }

    cd2 := &ColumnDef{
        name: "is_author",
        tdef: td,
    }

    sel := Select(cd1).Where(And(In(cd1, "foo", "bar"), Equal(cd2, 1)))

    exp := "SELECT name FROM users WHERE name IN (?, ?) AND is_author = ?"
    expLen := len(exp)
    expArgCount := 3

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectLimit(t *testing.T) {
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

    sel := Select(c).Limit(10)

    exp := "SELECT name FROM users LIMIT ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectLimitWithOffset(t *testing.T) {
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

    sel := Select(c).LimitWithOffset(10, 5)

    exp := "SELECT name FROM users LIMIT ? OFFSET ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectOrderByAsc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).OrderBy(cd.Asc())

    exp := "SELECT name FROM users ORDER BY name"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectOrderByMultiAscDesc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        tdef: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        tdef: td,
    }

    sel := Select(cd1).OrderBy(cd1.Asc(), cd2.Desc())

    exp := "SELECT name FROM users ORDER BY name, email DESC"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectStringArgs(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).Where(In(cd, "foo", "bar"))

    expStr := "SELECT name FROM users WHERE name IN (?, ?)"
    expLen := len(expStr)
    expArgCount := 2
    expArgs := []interface{}{"foo", "bar"}

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())

    actStr, actArgs := sel.StringArgs()

    assert.Equal(expStr, actStr)
    assert.Equal(expArgs, actArgs)
}

func TestSelectGroupByAsc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).GroupBy(cd)

    exp := "SELECT name FROM users GROUP BY name"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectGroupByMultiAscDesc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        tdef: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        tdef: td,
    }

    sel := Select(cd1).GroupBy(cd1, cd2)

    exp := "SELECT name FROM users GROUP BY name, email"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}

func TestSelectGroupOrderLimit(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    sel := Select(cd).GroupBy(cd).OrderBy(cd.Desc()).Limit(10)

    exp := "SELECT name FROM users GROUP BY name ORDER BY name DESC LIMIT ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.Size())
    assert.Equal(expArgCount, sel.ArgCount())
    assert.Equal(exp, sel.String())
}
