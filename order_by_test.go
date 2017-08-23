package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestOrderByClauseSingleAsc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &SortColumn{el: cd},
            },
        },
    }

    exp := " ORDER BY users.name"
    expLen := len(exp)
    expargCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
}

func TestOrderByClauseSingleDesc(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &SortColumn{el: cd, desc: true},
            },
        },
    }

    exp := " ORDER BY users.name DESC"
    expLen := len(exp)
    expargCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
}

func TestOrderByClauseMultiAsc(t *testing.T) {
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

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &SortColumn{el: cd1},
                &SortColumn{el: cd2},
            },
        },
    }

    exp := " ORDER BY users.name, users.email"
    expLen := len(exp)
    expargCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
}

func TestOrderByClauseMultiAscDesc(t *testing.T) {
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

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &SortColumn{el: cd1},
                &SortColumn{el: cd2, desc: true},
            },
        },
    }

    exp := " ORDER BY users.name, users.email DESC"
    expLen := len(exp)
    expargCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expargCount, argc)

    args := make([]interface{}, expargCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expargCount, numArgs)
}
