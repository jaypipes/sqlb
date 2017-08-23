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
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }

    val := &Value{value: "foo"}

    e := &Expression{
        scanInfo: exprscanTable[EXP_EQUAL],
        elements: []element{c, val},
    }

    exp := "users.name = ?"
    expLen := len(exp)
    expargCount := 1

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)

    // Check that if we reverse the order in which the Expression is constructed,
    // that our scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := &Expression{
        scanInfo: exprscanTable[EXP_EQUAL],
        elements: []element{val, c},
    }

    exp = "? = users.name"
    expLen = len(exp)
    expargCount = 1

    s = erev.size()
    assert.Equal(expLen, s)

    argc = erev.argCount()
    assert.Equal(expargCount, argc)

    args = make([]interface{}, 1)
    b = make([]byte, s)
    written, numArgs = erev.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)
}

func TestEqualFuncValue(t *testing.T) {
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

    eq := Equal(c, "foo")

    exp := "users.name = ?"
    expLen := len(exp)
    expargCount := 1

    s := eq.size()
    assert.Equal(expLen, s)

    argc := eq.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := eq.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)
}

func TestEqualFuncTwoelements(t *testing.T) {
    assert := assert.New(t)

    users := &TableDef{
        name: "users",
        schema: "test",
    }

    userId := &ColumnDef{
        name: "id",
        tdef: users,
    }

    c1 := &Column{
        cdef: userId,
        tbl: users.Table(),
    }

    articles := &TableDef{
        name: "articles",
        schema: "test",
    }

    author := &ColumnDef{
        name: "author",
        tdef: articles,
    }

    c2 := &Column{
        cdef: author,
        tbl: articles.Table(),
    }

    eq := Equal(c1, c2)

    exp := "users.id = articles.author"
    expLen := len(exp)
    expargCount := 0

    s := eq.size()
    assert.Equal(expLen, s)

    argc := eq.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 0)
    b := make([]byte, s)
    written, numArgs := eq.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(0, len(args))

    // Check that if we reverse the order in which the Expression is constructed,
    // that our scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := Equal(c2, c1)

    exp = "articles.author = users.id"
    expLen = len(exp)
    expargCount = 0

    s = erev.size()
    assert.Equal(expLen, s)

    argc = erev.argCount()
    assert.Equal(expargCount, argc)

    args = make([]interface{}, 0)
    b = make([]byte, s)
    written, numArgs = erev.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
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
        tdef: td,
    }

    c := &Column{
        cdef: cd,
        tbl: td.Table(),
    }

    val := &Value{value: "foo"}

    e := &Expression{
        scanInfo: exprscanTable[EXP_NEQUAL],
        elements: []element{c, val},
    }

    exp := "users.name != ?"
    expLen := len(exp)
    expargCount := 1

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)
}

func TestNotEqualFuncValue(t *testing.T) {
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

    eq := NotEqual(c, "foo")

    exp := "users.name != ?"
    expLen := len(exp)
    expargCount := 1

    s := eq.size()
    assert.Equal(expLen, s)

    argc := eq.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := eq.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal("foo", args[0])
    assert.Equal(expargCount, numArgs)
}

func TestNotEqualFuncTwoelements(t *testing.T) {
    assert := assert.New(t)

    users := &TableDef{
        name: "users",
        schema: "test",
    }

    userId := &ColumnDef{
        name: "id",
        tdef: users,
    }

    c1 := &Column{
        cdef: userId,
        tbl: users.Table(),
    }

    articles := &TableDef{
        name: "articles",
        schema: "test",
    }

    author := &ColumnDef{
        name: "author",
        tdef: articles,
    }

    c2 := &Column{
        cdef: author,
        tbl: articles.Table(),
    }

    eq := NotEqual(c1, c2)

    exp := "users.id != articles.author"
    expLen := len(exp)
    expargCount := 0

    s := eq.size()
    assert.Equal(expLen, s)

    argc := eq.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 0)
    b := make([]byte, s)
    written, numArgs := eq.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(0, len(args))

    // Check that if we reverse the order in which the Expression is constructed,
    // that our scan() still functions but merely generates the SQL string with
    // the left and right expression reversed

    erev := NotEqual(c2, c1)

    exp = "articles.author != users.id"
    expLen = len(exp)
    expargCount = 0

    s = erev.size()
    assert.Equal(expLen, s)

    argc = erev.argCount()
    assert.Equal(expargCount, argc)

    args = make([]interface{}, 0)
    b = make([]byte, s)
    written, numArgs = erev.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(0, len(args))
}

func TestInSingle(t *testing.T) {
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

    e := In(c, "foo")

    exp := "users.name IN (?)"
    expLen := len(exp)
    expargCount := 1

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 1)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(expargCount, len(args))
}

func TestInMulti(t *testing.T) {
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

    e := In(c, "foo", "bar", 1)

    exp := "users.name IN (?, ?, ?)"
    expLen := len(exp)
    expargCount := 3

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, 3)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(expargCount, len(args))
}

func TestAnd(t *testing.T) {
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

    ea := &Expression{
        scanInfo: exprscanTable[EXP_NEQUAL],
        elements: []element{c, &Value{value: "foo"}},
    }

    eb := &Expression{
        scanInfo: exprscanTable[EXP_NEQUAL],
        elements: []element{c, &Value{value: "bar"}},
    }
    e := And(ea, eb)

    exp := "users.name != ? AND users.name != ?"
    expLen := len(exp)
    expargCount := 2

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(expargCount, len(args))
    assert.Equal("foo", args[0])
    assert.Equal("bar", args[1])
}

func TestOr(t *testing.T) {
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

    ea := &Expression{
        scanInfo: exprscanTable[EXP_EQUAL],
        elements: []element{c, &Value{value: "foo"}},
    }

    eb := &Expression{
        scanInfo: exprscanTable[EXP_EQUAL],
        elements: []element{c, &Value{value: "bar"}},
    }
    e := Or(ea, eb)

    exp := "users.name = ? OR users.name = ?"
    expLen := len(exp)
    expargCount := 2

    s := e.size()
    assert.Equal(expLen, s)

    argc := e.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := e.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
    assert.Equal(expargCount, len(args))
    assert.Equal("foo", args[0])
    assert.Equal("bar", args[1])
}
