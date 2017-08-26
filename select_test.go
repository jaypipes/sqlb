package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type selClauseTest struct {
    sel *selectClause
    qs string
    qargs []interface{}
}

func TestSelectClause(t *testing.T) {
    assert := assert.New(t)

    tests := []selClauseTest{
        // TableDef and ColumnDef
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users},
                projected: &List{elements: []element{colUserName}},
            },
            qs: "SELECT users.name FROM users",
        },
        // Table and ColumnDef
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users.Table()},
                projected: &List{elements: []element{colUserName}},
            },
            qs: "SELECT users.name FROM users",
        },
        // TableDef and Column
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users},
                projected: &List{
                    elements: []element{
                        colUserName.Column(),
                    },
                },
            },
            qs: "SELECT users.name FROM users",
        },
        // aliased Table and Column
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users.As("u")},
                projected: &List{
                    elements: []element{
                        users.As("u").Column("name"),
                    },
                },
            },
            qs: "SELECT u.name FROM users AS u",
        },
        // TableDef and mutiple ColumnDef
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users},
                projected: &List{
                    elements: []element{
                        colUserId, colUserName,
                    },
                },
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // TableDef and mixed Column and ColumnDef
        selClauseTest{
            sel: &selectClause{
                selections: []selection{users},
                projected: &List{
                    elements: []element{
                        colUserId, colUserName.Column(),
                    },
                },
            },
            qs: "SELECT users.id, users.name FROM users",
        },
    }
    for _, test := range tests {
        expLen := len(test.qs)
        s := test.sel.size()
        assert.Equal(expLen, s)

        expArgc := len(test.qargs)
        assert.Equal(expArgc, test.sel.argCount())

        b := make([]byte, s)
        written, _ := test.sel.scan(b, test.qargs)

        assert.Equal(written, s)
        assert.Equal(test.qs, string(b))
    }
}

func TestWhereSingleEqual(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).Where(Equal(colUserName, "foo"))

    exp := "SELECT users.name FROM users WHERE users.name = ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectLimit(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).Limit(10)

    exp := "SELECT users.name FROM users LIMIT ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectLimitWithOffset(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).LimitWithOffset(10, 5)

    exp := "SELECT users.name FROM users LIMIT ? OFFSET ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectOrderByAsc(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).OrderBy(colUserName.Asc())

    exp := "SELECT users.name FROM users ORDER BY users.name"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectOrderByMultiAscDesc(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).OrderBy(colUserId.Asc(), colUserName.Desc())

    exp := "SELECT users.name FROM users ORDER BY users.id, users.name DESC"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectStringArgs(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).Where(In(colUserName, "foo", "bar"))

    expStr := "SELECT users.name FROM users WHERE users.name IN (?, ?)"
    expLen := len(expStr)
    expArgCount := 2
    expArgs := []interface{}{"foo", "bar"}

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())

    actStr, actArgs := sel.StringArgs()

    assert.Equal(expStr, actStr)
    assert.Equal(expArgs, actArgs)
}

func TestSelectGroupByAsc(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).GroupBy(colUserName)

    exp := "SELECT users.name FROM users GROUP BY users.name"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectGroupByMultiAsc(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).GroupBy(colUserId, colUserName)

    exp := "SELECT users.name FROM users GROUP BY users.id, users.name"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestSelectGroupOrderLimit(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).GroupBy(colUserName).OrderBy(colUserName.Desc()).Limit(10)

    exp := "SELECT users.name FROM users GROUP BY users.name ORDER BY users.name DESC LIMIT ?"
    expLen := len(exp)
    expArgCount := 1

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}
