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
    colUserName := users.Column("name")
    colUserId := users.Column("id")
    colArticleId := articles.Column("id")
    colArticleAuthor := articles.Column("author")

    tests := []selClauseTest{
        // TableDef and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserName},
            },
            qs: "SELECT users.name FROM users",
        },
        // Table and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users.Table()},
                projections: []projection{colUserName},
            },
            qs: "SELECT users.name FROM users",
        },
        // TableDef and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserName.Column()},
            },
            qs: "SELECT users.name FROM users",
        },
        // aliased Table and Column
        selClauseTest{
            c: &selectClause{
                selections: []selection{users.As("u")},
                projections: []projection{
                    users.As("u").Column("name"),
                },
            },
            qs: "SELECT u.name FROM users AS u",
        },
        // TableDef and multiple ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserId, colUserName},
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // TableDef and mixed Column and ColumnDef
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserId, colUserName.Column()},
            },
            qs: "SELECT users.id, users.name FROM users",
        },
        // Simple WHERE
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserName},
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
                projections: []projection{colUserName},
                limit: &limitClause{limit: 10},
            },
            qs: "SELECT users.name FROM users LIMIT ?",
            qargs: []interface{}{10},
        },
        // Simple ORDER BY
        selClauseTest{
            c: &selectClause{
                selections: []selection{users},
                projections: []projection{colUserName},
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
                projections: []projection{colUserName},
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
                projections: []projection{colUserName},
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
                projections: []projection{colUserName},
            },
            qs: "(SELECT users.name FROM users) AS u",
        },
        // Single join
        selClauseTest{
            c: &selectClause{
                selections: []selection{articles},
                projections: []projection{colArticleId, colUserName.As("author")},
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
