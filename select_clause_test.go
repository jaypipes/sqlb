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
    users := m.TableDef("users")
    articles := m.TableDef("articles")
    article_states := m.TableDef("article_states")
    colUserName := users.ColumnDef("name")
    colUserId := users.ColumnDef("id")
    colArticleId := articles.ColumnDef("id")
    colArticleAuthor := articles.ColumnDef("author")
    colArticleState := articles.ColumnDef("state")
    colArticleStateId := article_states.ColumnDef("id")
    colArticleStateName := article_states.ColumnDef("name")

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
        // TableDef and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName},
            },
            qs: "SELECT users.name FROM users",
        },
        // Table and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users.Table()},
                projs: []projection{colUserName},
            },
            qs: "SELECT users.name FROM users",
        },
        // TableDef and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserName.Column()},
            },
            qs: "SELECT users.name FROM users",
        },
        // aliased Table and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users.As("u")},
                projs: []projection{
                    users.As("u").Column("name"),
                },
            },
            qs: "SELECT u.name FROM users AS u",
        },
        // TableDef and multiple ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserId, colUserName},
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // TableDef and mixed Column and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{colUserId, colUserName.Column()},
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
        // Simple sub-SELECT
        selClauseTest{
            c: &selectClause{
                alias: "u",
                selections: []selection{users},
                projs: []projection{colUserName},
            },
            qs: "(SELECT users.name FROM users) AS u",
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
                        onExprs: []*Expression{
                            Equal(colArticleAuthor, colUserId),
                        },
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
                        onExprs: []*Expression{
                            Equal(colArticleAuthor, colUserId),
                        },
                    },
                    &joinClause{
                        left: articles,
                        right: article_states,
                        onExprs: []*Expression{
                            Equal(colArticleState, colArticleStateId),
                        },
                    },
                },
            },
            qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
        },
        // Simple COUNT(*) on a table
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projs: []projection{Count()},
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
