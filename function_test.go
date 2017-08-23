package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFuncWithAlias(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    m := Max(cd).As("max_created_on")

    exp := "MAX(users.created_on) AS max_created_on"
    expLen := len(exp)
    expArgCount := 0

    s := m.size()
    assert.Equal(expLen, s)

    argc := m.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncMax(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    m := Max(cd)

    exp := "MAX(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.size()
    assert.Equal(expLen, s)

    argc := m.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncMaxColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    m := cd.Max()

    exp := "MAX(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.size()
    assert.Equal(expLen, s)

    argc := m.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }
    m = c.Max()

    s = m.size()
    assert.Equal(expLen, s)

    argc = m.argCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncMin(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    m := Min(cd)

    exp := "MIN(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.size()
    assert.Equal(expLen, s)

    argc := m.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncMinColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    m := cd.Min()

    exp := "MIN(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.size()
    assert.Equal(expLen, s)

    argc := m.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }
    m = c.Min()

    s = m.size()
    assert.Equal(expLen, s)

    argc = m.argCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = m.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncSum(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    f := Sum(cd)

    exp := "SUM(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := f.size()
    assert.Equal(expLen, s)

    argc := f.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncSumColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    f := cd.Sum()

    exp := "SUM(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := f.size()
    assert.Equal(expLen, s)

    argc := f.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }
    f = c.Sum()

    s = f.size()
    assert.Equal(expLen, s)

    argc = f.argCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncAvg(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    f := Avg(cd)

    exp := "AVG(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := f.size()
    assert.Equal(expLen, s)

    argc := f.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestFuncAvgColumn(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "created_on",
        tdef: td,
    }

    f := cd.Avg()

    exp := "AVG(users.created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := f.size()
    assert.Equal(expLen, s)

    argc := f.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }
    f = c.Avg()

    s = f.size()
    assert.Equal(expLen, s)

    argc = f.argCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = f.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}
