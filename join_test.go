package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type joinClauseTest struct {
    c *joinClause
    qs string
    qargs []interface{}
}

func TestJoinClause(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.TableDef("users")
    articles := m.TableDef("articles")
    colUserId := users.Column("id")
    colArticleAuthor := articles.Column("author")

    auCond := Equal(colArticleAuthor, colUserId)
    uaCond := Equal(colUserId, colArticleAuthor)

    join := &joinClause{
        left: articles.Table(),
        right: users.Table(),
        onExprs: []*Expression{},
    }
    tests := []joinClauseTest{
        // articles to users table defs
        joinClauseTest{
            c: Join(articles, users, auCond),
            qs: " JOIN users ON articles.author = users.id",
        },
        // users to articles table defs
        joinClauseTest{
            c: Join(users, articles, uaCond),
            qs: " JOIN articles ON users.id = articles.author",
        },
        // articles to users tables
        joinClauseTest{
            c: Join(articles.Table(), users.Table(), auCond),
            qs: " JOIN users ON articles.author = users.id",
        },
        // joinClause.On() method
        joinClauseTest{
            c: join.On(auCond),
            qs: " JOIN users ON articles.author = users.id",
        },
        // join an aliased table to non-aliased table
        joinClauseTest{
            c: &joinClause{
                left: articles.As("a"),
                right: users.Table(),
                onExprs: []*Expression{
                    Equal(articles.As("a").Column("author"), colUserId),
                },
            },
            qs: " JOIN users ON a.author = users.id",
        },
        // join a non-aliased table to aliased table
        joinClauseTest{
            c: &joinClause{
                left: articles,
                right: users.As("u"),
                onExprs: []*Expression{
                    Equal(colArticleAuthor, users.As("u").Column("id")),
                },
            },
            qs: " JOIN users AS u ON articles.author = u.id",
        },
        // aliased projections should not include "AS alias" in output
        joinClauseTest{
            c: &joinClause{
                left: articles,
                right: users,
                onExprs: []*Expression{
                    Equal(colArticleAuthor, colUserId.As("user_id")),
                },
            },
            qs: " JOIN users ON articles.author = users.id",
        },
    }
    for _, test := range tests {
        expLen := len(test.qs)
        s := test.c.size()
        assert.Equal(expLen, s)

        expArgc := len(test.qargs)
        assert.Equal(expArgc, test.c.argCount())

        b := make([]byte, s)
        written, _ := test.c.scan(b, test.qargs)

        assert.Equal(written, s)
        assert.Equal(test.qs, string(b))
    }
}
