package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type orderByTest struct {
    c *orderByClause
    qs string
    qargs []interface{}
}

func TestOrderBy(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.TableDef("users")
    colUserId := users.ColumnDef("id")
    colUserName := users.ColumnDef("name")

    tests := []orderByTest{
        // column def asc
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{colUserName.Asc()},
            },
            qs: " ORDER BY users.name",
        },
        // column asc
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{colUserName.Column().Asc()},
            },
            qs: " ORDER BY users.name",
        },
        // column def desc
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{colUserName.Desc()},
            },
            qs: " ORDER BY users.name DESC",
        },
        // column desc
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{colUserName.Column().Desc()},
            },
            qs: " ORDER BY users.name DESC",
        },
        // Aliased column should NOT output alias in ORDER BY
        orderByTest{
            c: &orderByClause{
               scols: []*sortColumn{colUserName.As("user_name").Desc()},
            },
            qs: " ORDER BY users.name DESC",
        },
        // multi column mixed
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{colUserName.Asc(), colUserId.Column().Desc()},
            },
            qs: " ORDER BY users.name, users.id DESC",
        },
        // sort by a function
        orderByTest{
            c: &orderByClause{
                scols: []*sortColumn{Count().Desc()},
            },
            qs: " ORDER BY COUNT(*) DESC",
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
