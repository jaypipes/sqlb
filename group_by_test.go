package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGroupByClauseSingle(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd := &ColumnDef{
        name: "name",
        tdef: td,
    }

    ob := &GroupByClause{
        cols: &List{
            elements: []element{cd},
        },
    }

    exp := " GROUP BY users.name"
    expLen := len(exp)
    expArgCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestGroupByClauseMulti(t *testing.T) {
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

    ob := &GroupByClause{
        cols: &List{
            elements: []element{cd1, cd2},
        },
    }

    exp := " GROUP BY users.name, users.email"
    expLen := len(exp)
    expArgCount := 0

    s := ob.size()
    assert.Equal(expLen, s)

    argc := ob.argCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := ob.scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}
