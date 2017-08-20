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
        table: td,
    }

    m := Max(cd).As("max_created_on")

    exp := "MAX(created_on) AS max_created_on"
    expLen := len(exp)
    expArgCount := 0

    s := m.Size()
    assert.Equal(expLen, s)

    argc := m.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.Scan(b, args)

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
        table: td,
    }

    m := Max(cd)

    exp := "MAX(created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.Size()
    assert.Equal(expLen, s)

    argc := m.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.Scan(b, args)

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
        table: td,
    }

    m := cd.Max()

    exp := "MAX(created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.Size()
    assert.Equal(expLen, s)

    argc := m.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        def: cd,
    }
    m = c.Max()

    s = m.Size()
    assert.Equal(expLen, s)

    argc = m.ArgCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = m.Scan(b, args)

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
        table: td,
    }

    m := Min(cd)

    exp := "MIN(created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.Size()
    assert.Equal(expLen, s)

    argc := m.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.Scan(b, args)

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
        table: td,
    }

    m := cd.Min()

    exp := "MIN(created_on)"
    expLen := len(exp)
    expArgCount := 0

    s := m.Size()
    assert.Equal(expLen, s)

    argc := m.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := m.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)

    // Test with Column not ColumnDef
    c := &Column{
        def: cd,
    }
    m = c.Min()

    s = m.Size()
    assert.Equal(expLen, s)

    argc = m.ArgCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, expArgCount)
    b = make([]byte, s)
    written, numArgs = m.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}
