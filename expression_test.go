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

func TestEqualFuncLiteral(t *testing.T) {
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

    eq := Equal(c, "foo")

    exp := "name = ?"
    expLen := len(exp)
    expArgCount := 1

    s := eq.Size()
    assert.Equal(expLen, s)

    argc := eq.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := eq.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}

func TestEqualFuncTwoElements(t *testing.T) {
    assert := assert.New(t)

    users := &TableDef{
        name: "users",
        schema: "test",
    }

    userId := &ColumnDef{
        name: "id",
        table: users,
    }

    c1 := &Column{
        def: userId,
    }

    articles := &TableDef{
        name: "articles",
        schema: "test",
    }

    author := &ColumnDef{
        name: "author",
        table: articles,
    }

    c2 := &Column{
        def: author,
    }

    eq := Equal(c1, c2)

    exp := "id = author"
    expLen := len(exp)
    expArgCount := 0

    s := eq.Size()
    assert.Equal(expLen, s)

    argc := eq.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 0)
    b := make([]byte, s)
    written, numArgs := eq.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal(0, len(args))

    // Check that if we reverse the order in which the Expression is constructed,
    // that our Scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := Equal(c2, c1)

    exp = "author = id"
    expLen = len(exp)
    expArgCount = 0

    s = erev.Size()
    assert.Equal(expLen, s)

    argc = erev.ArgCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, 0)
    b = make([]byte, s)
    written, numArgs = erev.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal(0, len(args))
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

func TestNotEqualFuncLiteral(t *testing.T) {
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

    eq := NotEqual(c, "foo")

    exp := "name != ?"
    expLen := len(exp)
    expArgCount := 1

    s := eq.Size()
    assert.Equal(expLen, s)

    argc := eq.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := eq.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expArgCount, numArgs)
}

func TestNotEqualFuncTwoElements(t *testing.T) {
    assert := assert.New(t)

    users := &TableDef{
        name: "users",
        schema: "test",
    }

    userId := &ColumnDef{
        name: "id",
        table: users,
    }

    c1 := &Column{
        def: userId,
    }

    articles := &TableDef{
        name: "articles",
        schema: "test",
    }

    author := &ColumnDef{
        name: "author",
        table: articles,
    }

    c2 := &Column{
        def: author,
    }

    eq := NotEqual(c1, c2)

    exp := "id != author"
    expLen := len(exp)
    expArgCount := 0

    s := eq.Size()
    assert.Equal(expLen, s)

    argc := eq.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, 0)
    b := make([]byte, s)
    written, numArgs := eq.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal(0, len(args))

    // Check that if we reverse the order in which the Expression is constructed,
    // that our Scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := NotEqual(c2, c1)

    exp = "author != id"
    expLen = len(exp)
    expArgCount = 0

    s = erev.Size()
    assert.Equal(expLen, s)

    argc = erev.ArgCount()
    assert.Equal(expArgCount, argc)

    args = make([]interface{}, 0)
    b = make([]byte, s)
    written, numArgs = erev.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal(0, len(args))
}