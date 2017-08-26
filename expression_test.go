package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type expressionTest struct {
    c *Expression
    qs string
    qargs []interface{}
}

func TestExpressions(t *testing.T) {
    assert := assert.New(t)

    tests := []expressionTest{
        // equal value
        expressionTest{
            c: Equal(colUserName, "foo"),
            qs: "users.name = ?",
            qargs: []interface{}{"foo"},
        },
        // reverse args equal
        expressionTest{
            c: Equal("foo", colUserName),
            qs: "? = users.name",
            qargs: []interface{}{"foo"},
        },
        // equal columns
        expressionTest{
            c: Equal(colUserId, colArticleAuthor),
            qs: "users.id = articles.author",
        },
        // not equal value
        expressionTest{
            c: NotEqual(colUserName, "foo"),
            qs: "users.name != ?",
            qargs: []interface{}{"foo"},
        },
        // in single value
        expressionTest{
            c: In(colUserName, "foo"),
            qs: "users.name IN (?)",
            qargs: []interface{}{"foo"},
        },
        // in multi value
        expressionTest{
            c: In(colUserName, "foo", "bar", 1),
            qs: "users.name IN (?, ?, ?)",
            qargs: []interface{}{"foo", "bar", 1},
        },
        // AND expression
        expressionTest{
            c: And(
                NotEqual(colUserName, "foo"),
                NotEqual(colUserName, "bar"),
            ),
            qs: "users.name != ? AND users.name != ?",
            qargs: []interface{}{"foo", "bar"},
        },
        // OR expression
        expressionTest{
            c: Or(
                Equal(colUserName, "foo"),
                Equal(colUserName, "bar"),
            ),
            qs: "users.name = ? OR users.name = ?",
            qargs: []interface{}{"foo", "bar"},
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
