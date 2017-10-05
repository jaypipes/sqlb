package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type selClauseTest struct {
    c *selectClause
    qs string
    qargs []interface{}
}

func TestSelectClause(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    articles := m.Table("articles")
    article_states := m.Table("article_states")
    colUserName := users.C("name")
    colUserId := users.C("id")
    colArticleId := articles.C("id")
    colArticleAuthor := articles.C("author")
    colArticleState := articles.C("state")
    colArticleStateId := article_states.C("id")
    colArticleStateName := article_states.C("name")

    tests := []selClauseTest{
        // a literal value
        selClauseTest{
            c: &selectClause{
                projs: []projection{&value{val: 1}},
            },
            qs: "SELECT ?",
            qargs: []interface{}{1},
        },
        // a literal value aliased
        selClauseTest{
            c: &selectClause{
                projs: []projection{
                    &value{alias: "foo", val: 1},
                },
            },
            qs: "SELECT ? AS foo",
            qargs: []interface{}{1},
        },
        // two literal values
        selClauseTest{
            c: &selectClause{
                projs: []projection{
                    &value{val: 1},
                    &value{val: 1},
                },
            },
            qs: "SELECT ?, ?",
            qargs: []interface{}{1, 2},
        },
        // Table and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
            },
            qs: "SELECT users.name FROM users",
        },
        // aliased Table and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users.As("u")},
                projs: []projection{
                    users.As("u").C("name"),
                },
            },
            qs: "SELECT u.name FROM users AS u",
        },
        // Table and multiple Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserId, colUserName},
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // Table and mixed Column and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserId, colUserName},
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // Simple WHERE
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
                where: &whereClause{
                    filters: []*Expression{
                        Equal(colUserName, "foo"),
                    },
                },
            },
            qs: "SELECT users.name FROM users WHERE users.name = ?",
            qargs: []interface{}{"foo"},
        },
        // Simple LIMIT
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
                limit: &limitClause{limit: 10},
            },
            qs: "SELECT users.name FROM users LIMIT ?",
            qargs: []interface{}{10},
        },
        // Simple ORDER BY
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
                orderBy: &orderByClause{
                    scols: []*sortColumn{colUserName.Desc()},
                },
            },
            qs: "SELECT users.name FROM users ORDER BY users.name DESC",
        },
        // Simple GROUP BY
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
                groupBy: &groupByClause{
                    cols: []projection{colUserName},
                },
            },
            qs: "SELECT users.name FROM users GROUP BY users.name",
        },
        // GROUP BY, ORDER BY and LIMIT
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
                groupBy: &groupByClause{
                    cols: []projection{colUserName},
                },
                orderBy: &orderByClause{
                    scols: []*sortColumn{colUserName.Desc()},
                },
                limit: &limitClause{limit: 10},
            },
            qs: "SELECT users.name FROM users GROUP BY users.name ORDER BY users.name DESC LIMIT ?",
            qargs: []interface{}{10},
        },
        // Single join
        selClauseTest{
            c: &selectClause{
                selections: []selection{articles},
                projs: []projection{colArticleId, colUserName.As("author")},
                joins: []*joinClause{
                    &joinClause{
                        left: articles,
                        right: users,
                        on: Equal(colArticleAuthor, colUserId),
                    },
                },
            },
            qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
        },
        // Multiple joins
        selClauseTest{
            c: &selectClause{
                selections: []selection{articles},
                projs: []projection{colArticleId, colUserName.As("author"), colArticleStateName.As("state")},
                joins: []*joinClause{
                    &joinClause{
                        left: articles,
                        right: users,
                        on: Equal(colArticleAuthor, colUserId),
                    },
                    &joinClause{
                        left: articles,
                        right: article_states,
                        on: Equal(colArticleState, colArticleStateId),
                    },
                },
            },
            qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
        },
        // Simple COUNT(*) on a table
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{Count(users)},
            },
            qs: "SELECT COUNT(*) FROM users",
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
