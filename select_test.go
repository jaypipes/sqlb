package sqlb

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/assert"
)

type selectQueryTest struct {
    q *SelectQuery
    qs string
    qargs []interface{}
    qe error
}

func TestSelectQuery(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    articles := m.TableDef("articles")
    article_states := m.TableDef("article_states")
    colUserId := users.Column("id")
    colUserName := users.Column("name")
    colArticleId := articles.Column("id")
    colArticleAuthor := articles.Column("author")
    colArticleState := articles.ColumnDef("state")
    colArticleStateId := article_states.ColumnDef("id")
    colArticleStateName := article_states.ColumnDef("name")

    tests := []selectQueryTest{
        // Simple FROM
        selectQueryTest{
            q: Select(users),
            qs: "SELECT users.id, users.name FROM users",
        },
        // Simple WHERE
        selectQueryTest{
            q: Select(users).Where(Equal(colUserName, "foo")),
            qs: "SELECT users.id, users.name FROM users WHERE users.name = ?",
            qargs: []interface{}{"foo"},
        },
        // Simple GROUP BY
        selectQueryTest{
            q: Select(users).GroupBy(colUserName),
            qs: "SELECT users.id, users.name FROM users GROUP BY users.name",
        },
        // Simple ORDER BY
        selectQueryTest{
            q: Select(users).OrderBy(colUserName.Desc()),
            qs: "SELECT users.id, users.name FROM users ORDER BY users.name DESC",
        },
        // Simple LIMIT
        selectQueryTest{
            q: Select(users).Limit(10),
            qs: "SELECT users.id, users.name FROM users LIMIT ?",
            qargs: []interface{}{10},
        },
        // Simple LIMIT with OFFSET
        selectQueryTest{
            q: Select(users).LimitWithOffset(10, 20),
            qs: "SELECT users.id, users.name FROM users LIMIT ? OFFSET ?",
            qargs: []interface{}{10, 20},
        },
        // Simple sub-SELECT
        selectQueryTest{
            q: Select(users).As("u"),
            // This is incorrect. Should be:
            // qs: "SELECT u.id, u.name FROM (SELECT users.id, users.name FROM users) AS u",
            qs: "SELECT users.id, users.name FROM (SELECT users.id, users.name FROM users) AS u",
        },
        // Simple INNER JOIN
        selectQueryTest{
            q: Select(colArticleId, colUserName.As("author")).Join(users, Equal(colArticleAuthor, colUserId)),
            qs: "SELECT articles.id, users.name AS author FROM articles JOIN users ON articles.author = users.id",
        },
        // Two JOINs using Join() method
        selectQueryTest{
            q: Select(
                colArticleId,
                colUserName.As("author"),
                colArticleStateName.As("state"),
            ).Join(users, Equal(colArticleAuthor, colUserId),
            ).Join(article_states, Equal(colArticleState, colArticleStateId)),
            qs: "SELECT articles.id, users.name AS author, article_states.name AS state FROM articles JOIN users ON articles.author = users.id JOIN article_states ON articles.state = article_states.id",
        },
    }
    for _, test := range tests {
        if test.qe != nil {
            assert.Equal(test.qe, test.q.Error())
            continue
        } else if test.q.Error() != nil {
            qe := test.q.Error()
            assert.Fail(qe.Error())
            continue
        }
        qs, qargs := test.q.StringArgs()
        assert.Equal(len(test.qargs), len(qargs))
        assert.Equal(test.qs, qs)
    }
}

func TestModifyingSelectQueryUpdatesBuffer(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.TableDef("users")

    q := Select(users)

    qs, qargs := q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users", qs)
    assert.Nil(qargs)

    // Modify the underlying SELECT and verify string and args changed
    q.Where(Equal(users.Column("id"), 1))
    qs, qargs = q.StringArgs()
    assert.Equal("SELECT users.id, users.name FROM users WHERE users.id = ?", qs)
    assert.Equal([]interface{}{1}, qargs)
}

func TestSelectQueryErrors(t *testing.T) {
    assert := assert.New(t)

    q := &SelectQuery{}

    assert.False(q.IsValid()) // Doesn't have a selectClause yet...
    assert.Nil(q.Error()) // But there is no error set yet...

    m := testFixtureMeta()
    users := m.TableDef("users")

    q = Select(users)

    assert.True(q.IsValid())
    assert.Nil(q.Error())

    q.e = fmt.Errorf("Cannot determine left side of JOIN expression.")
    assert.False(q.IsValid())
    assert.NotNil(q.Error())
}
