package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExpressionEqual(t *testing.T) {
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

    lit := &Literal{value: "foo"}

    e := &Expression{
        op: OP_EQUAL,
        left: c,
        right: lit,
    }

    exp := "name = ?"
    expLen := len(exp)
    expArgCount := 1

    s := e.Size()
    assert.Equal(expLen, s)

    argc := e.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := e.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)

    // Check that if we reverse the order in which the Expression is constructed,
    // that our Scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := &Expression{
        op: OP_EQUAL,
        left: lit,
        right: c,
    }

    exp = "? = name"
    expLen = len(exp)
    expArgCount = 1

    s = erev.Size()
    assert.Equal(expLen, s)

    argc = erev.ArgCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, 1)
    b = make([]byte, s)
    written, numArgs = erev.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}

func TestExpressionNotEqual(t *testing.T) {
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

    lit := &Literal{value: "foo"}

    e := &Expression{
        op: OP_NEQUAL,
        left: c,
        right: lit,
    }

    exp := "name != ?"
    expLen := len(exp)
    expArgCount := 1

    s := e.Size()
    assert.Equal(expLen, s)

    argc := e.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := e.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}
