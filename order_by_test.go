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
            elements: []Element{
                &SortColumn{el: cd},
            },
        },
    }

    exp := " ORDER BY users.name"
    expLen := len(exp)
    expArgCount := 0

    s := ob.Size()
    assert.Equal(expLen, s)

    argc := ob.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
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
            elements: []Element{
                &SortColumn{el: cd, desc: true},
            },
        },
    }

    exp := " ORDER BY users.name DESC"
    expLen := len(exp)
    expArgCount := 0

    s := ob.Size()
    assert.Equal(expLen, s)

    argc := ob.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
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
            elements: []Element{
                &SortColumn{el: cd1},
                &SortColumn{el: cd2},
            },
        },
    }

    exp := " ORDER BY users.name, users.email"
    expLen := len(exp)
    expArgCount := 0

    s := ob.Size()
    assert.Equal(expLen, s)

    argc := ob.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
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
            elements: []Element{
                &SortColumn{el: cd1},
                &SortColumn{el: cd2, desc: true},
            },
        },
    }

    exp := " ORDER BY users.name, users.email DESC"
    expLen := len(exp)
    expArgCount := 0

    s := ob.Size()
    assert.Equal(expLen, s)

    argc := ob.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}
