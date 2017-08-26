package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSelectSingleColumn(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    sel := Select(c)

    exp := "SELECT users.name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
}

func TestSelectSingleColumnWithTableAlias(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.As("u"),
    }

    sel := Select(c)

    exp := "SELECT u.name FROM users AS u"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
}

func TestSelectMultiColumnsSingleTable(t *testing.T) {
    assert := assert.New(t)

    c1 := &Column{
        cdef: colUserId,
        tbl: users.Table(),
    }

    c2 := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    sel := Select(c1, c2)

    exp := "SELECT users.id, users.name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromColumnDef(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName)

    exp := "SELECT users.name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromColumnDefAndColumn(t *testing.T) {
    assert := assert.New(t)

    c := &Column{
        cdef: colUserName,
        tbl: users.Table(),
    }

    sel := Select(colUserId, c)

    exp := "SELECT users.id, users.name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
}

func TestSelectFromTableDef(t *testing.T) {
    assert := assert.New(t)

    sel := Select(users)

    exp := "SELECT users.id, users.name FROM users"
    expLen := len(exp)

    assert.Equal(expLen, sel.size())
    assert.Equal(exp, sel.String())
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

func TestWhereSingleAnd(t *testing.T) {
    assert := assert.New(t)

    cond := And(
        NotEqual(colUserName, "foo"),
        NotEqual(colUserName, "bar"),
    )
    sel := Select(colUserName).Where(cond)

    exp := "SELECT users.name FROM users WHERE users.name != ? AND users.name != ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestWhereSingleIn(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).Where(In(colUserName, "foo", "bar"))

    exp := "SELECT users.name FROM users WHERE users.name IN (?, ?)"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestWhereMultiNotEqual(t *testing.T) {
    assert := assert.New(t)

    sel := Select(colUserName).Where(NotEqual(colUserName, "foo"))
    sel = sel.Where(NotEqual(colUserName, "bar"))

    exp := "SELECT users.name FROM users WHERE users.name != ? AND users.name != ?"
    expLen := len(exp)
    expArgCount := 2

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}

func TestWhereMultiInAndEqual(t *testing.T) {
    assert := assert.New(t)

    colUserIsAuthor := &ColumnDef{
        name: "is_author",
        tdef: users,
    }

    cond := And(In(colUserName, "foo", "bar"), Equal(colUserIsAuthor, 1))
    sel := Select(colUserName).Where(cond)

    exp := "SELECT users.name FROM users WHERE users.name IN (?, ?) AND users.is_author = ?"
    expLen := len(exp)
    expArgCount := 3

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

func TestSelectJoinSingle(t *testing.T) {
    assert := assert.New(t)

    j := &JoinClause{
        left: users,
        right: articles,
        onExprs: []*Expression{
            Equal(colUserId, colArticleAuthor),
        },
    }

    sel := Select(j)

    exp := "SELECT users.id, users.name, articles.id, articles.author FROM users JOIN articles ON users.id = articles.author"
    expLen := len(exp)
    expArgCount := 0

    assert.Equal(expLen, sel.size())
    assert.Equal(expArgCount, sel.argCount())
    assert.Equal(exp, sel.String())
}
