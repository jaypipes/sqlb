package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type functionTest struct {
    c *sqlFunc
    qs string
    qargs []interface{}
}

func TestFunctions(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    colUserName := users.Column("name")

    tests := []functionTest{
        // MAX column
        functionTest{
            c: Max(colUserName),
            qs: "MAX(users.name)",
        },
        // function aliased
        functionTest{
            c: Max(colUserName).As("max_name"),
            qs: "MAX(users.name) AS max_name",
        },
        // MIN column
        functionTest{
            c: Min(colUserName),
            qs: "MIN(users.name)",
        },
        // SUM column
        functionTest{
            c: Sum(colUserName),
            qs: "SUM(users.name)",
        },
        // AVG column
        functionTest{
            c: Avg(colUserName),
            qs: "AVG(users.name)",
        },
        // COUNT(*)
        functionTest{
            c: Count(),
            qs: "COUNT(*)",
        },
        // COUNT(DISTINCT column)
        functionTest{
            c: CountDistinct(colUserName),
            qs: "COUNT(DISTINCT users.name)",
        },
        // Ensure AS alias does not appear for internal projection
        functionTest{
            c: CountDistinct(colUserName.As("user_name")),
            qs: "COUNT(DISTINCT users.name)",
        },
        // CAST(col AS type))
        functionTest{
            c: Cast(colUserName, SQL_TYPE_TEXT),
            qs: "CAST(users.name AS TEXT)",
        },
        // TRIM column
        functionTest{
            c: Trim(colUserName),
            qs: "TRIM(users.name)",
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
