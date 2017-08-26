package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestOrderByClauseSingleAsc(t *testing.T) {
    assert := assert.New(t)

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &sortColumn{el: colUserName},
            },
        },
    }

    exp := " ORDER BY users.name"
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

func TestOrderByClauseSingleDesc(t *testing.T) {
    assert := assert.New(t)

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &sortColumn{el: colUserName, desc: true},
            },
        },
    }

    exp := " ORDER BY users.name DESC"
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

func TestOrderByClauseMultiAsc(t *testing.T) {
    assert := assert.New(t)

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &sortColumn{el: colUserId},
                &sortColumn{el: colUserName},
            },
        },
    }

    exp := " ORDER BY users.id, users.name"
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

func TestOrderByClauseMultiAscDesc(t *testing.T) {
    assert := assert.New(t)

    ob := &OrderByClause{
        cols: &List{
            elements: []element{
                &sortColumn{el: colUserId},
                &sortColumn{el: colUserName, desc: true},
            },
        },
    }

    exp := " ORDER BY users.id, users.name DESC"
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
