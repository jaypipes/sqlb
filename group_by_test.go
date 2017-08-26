package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestgroupByClauseSingle(t *testing.T) {
    assert := assert.New(t)

    ob := &groupByClause{
        cols: &List{
            elements: []element{colUserName},
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

func TestgroupByClauseMulti(t *testing.T) {
    assert := assert.New(t)

    ob := &groupByClause{
        cols: &List{
            elements: []element{colUserId, colUserName},
        },
    }

    exp := " GROUP BY users.id, users.name"
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
