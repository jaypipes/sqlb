package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

var (
    users = &TableDef{
        name: "users",
        schema: "test",
    }

    colUserId = &ColumnDef{
        name: "id",
        tdef: users,
    }

    colUserName = &ColumnDef{
        name: "name",
        tdef: users,
    }

    articles = &TableDef{
        name: "articles",
        schema: "test",
    }

    colArticleId = &ColumnDef{
        name: "id",
        tdef: articles,
    }

    colArticleAuthor = &ColumnDef{
        name: "author",
        tdef: articles,
    }
)

func init() {
    users.cdefs = []*ColumnDef{colUserId, colUserName}
    articles.cdefs = []*ColumnDef{colArticleId, colArticleAuthor}
}

func TestJoinFuncGenerics(t *testing.T) {
    // Test that the sqlb.Join() func can take a *Table or *TableDef and zero
    // or more *Expression struct pointers and returns a *JoinClause struct
    // pointers. Essentially, we're testing the Selection generic interface here
    assert := assert.New(t)

    cond := Equal(colArticleAuthor, colUserId)

    joins := []*JoinClause{
        Join(articles, users, cond),
        Join(articles.Table(), users.Table(), cond),
    }

    for _, j := range joins {
        exp := " JOIN users ON author = id"
        expLen := len(exp)
        expArgCount := 0

        s := j.Size()
        assert.Equal(expLen, s)

        argc := j.ArgCount()
        assert.Equal(expArgCount, argc)

        args := make([]interface{}, expArgCount)
        b := make([]byte, s)
        written, numArgs := j.Scan(b, args)

        assert.Equal(s, written)
        assert.Equal(exp, string(b))
        assert.Equal(expArgCount, numArgs)
    }
}

func TestJoinClauseInnerOnEqualSingle(t *testing.T) {
    assert := assert.New(t)

    j := &JoinClause{
        left: articles.Table(),
        right: users.Table(),
        onExprs: []*Expression{
            Equal(colArticleAuthor, colUserId),
        },
    }

    exp := " JOIN users ON author = id"
    expLen := len(exp)
    expArgCount := 0

    s := j.Size()
    assert.Equal(expLen, s)

    argc := j.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := j.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestJoinClauseOnMethod(t *testing.T) {
    assert := assert.New(t)

    j := &JoinClause{
        left: articles.Table(),
        right: users.Table(),
    }
    j.On(Equal(colArticleAuthor, colUserId))

    exp := " JOIN users ON author = id"
    expLen := len(exp)
    expArgCount := 0

    s := j.Size()
    assert.Equal(expLen, s)

    argc := j.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := j.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestJoinClauseAliasedInnerOnEqualSingle(t *testing.T) {
    assert := assert.New(t)

    atbl := articles.Table().As("a")
    utbl := users.Table().As("u")

    aliasAuthorCol := atbl.Column("author")
    assert.NotNil(aliasAuthorCol)

    aliasIdCol := utbl.Column("id")
    assert.NotNil(aliasIdCol)

    j := &JoinClause{
        left: atbl,
        right: utbl,
        onExprs: []*Expression{
            Equal(aliasAuthorCol, aliasIdCol),
        },
    }

    exp := " JOIN users AS u ON a.author = u.id"
    expLen := len(exp)
    expArgCount := 0

    s := j.Size()
    assert.Equal(expLen, s)

    argc := j.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := j.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
}

func TestJoinClauseInnerOnEqualMulti(t *testing.T) {
    assert := assert.New(t)

    j := &JoinClause{
        left: articles.Table(),
        right: users.Table(),
        onExprs: []*Expression{
            Equal(colArticleAuthor, colUserId),
            Equal(colUserName, "foo"),
        },
    }

    exp := " JOIN users ON author = id AND name = ?"
    expLen := len(exp)
    expArgCount := 1

    s := j.Size()
    assert.Equal(expLen, s)

    argc := j.ArgCount()
    assert.Equal(expArgCount, argc)

    args := make([]interface{}, expArgCount)
    b := make([]byte, s)
    written, numArgs := j.Scan(b, args)

    assert.Equal(s, written)
    assert.Equal(exp, string(b))
    assert.Equal(expArgCount, numArgs)
    assert.Equal("foo", args[0])
}
