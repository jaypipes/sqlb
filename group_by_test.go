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
        table: td,
    }

    ob := &GroupByClause{
        cols: &List{
            elements: []Element{cd},
        },
    }

    exp := " GROUP BY name"
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

func TestGroupByClauseMulti(t *testing.T) {
    assert := assert.New(t)

    td := &TableDef{
        name: "users",
        schema: "test",
    }

    cd1 := &ColumnDef{
        name: "name",
        table: td,
    }

    cd2 := &ColumnDef{
        name: "email",
        table: td,
    }

    ob := &GroupByClause{
        cols: &List{
            elements: []Element{cd1, cd2},
        },
    }

    exp := " GROUP BY name, email"
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
