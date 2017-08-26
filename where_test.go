package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

type whereClauseTest struct {
    c *whereClause
    qs string
    qargs []interface{}
}

func TestWhereClause(t *testing.T) {
    assert := assert.New(t)

    tests := []whereClauseTest{
        // Empty clause
        whereClauseTest{
            c: &whereClause{},
            qs: "",
        },
        // Single expression
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    Equal(colUserName, "foo"),
                },
            },
            qs: " WHERE users.name = ?",
            qargs: []interface{}{"foo"},
        },
        // Single AND expression
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    And(
                        NotEqual(colUserName, "foo"),
                        NotEqual(colUserName, "bar"),
                    ),
                },
            },
            qs: " WHERE (users.name != ? AND users.name != ?)",
            qargs: []interface{}{"foo", "bar"},
        },
        // Multiple unary expressions should be AND'd together
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    NotEqual(colUserName, "foo"),
                    NotEqual(colUserName, "bar"),
                },
            },
            qs: " WHERE users.name != ? AND users.name != ?",
            qargs: []interface{}{"foo", "bar"},
        },
        // Single OR expression
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    Or(
                        Equal(colUserName, "foo"),
                        Equal(colUserName, "bar"),
                    ),
                },
            },
            qs: " WHERE (users.name = ? OR users.name = ?)",
            qargs: []interface{}{"foo", "bar"},
        },
        // An OR and another unary expression
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    Or(
                        Equal(colUserName, "foo"),
                        Equal(colUserName, "bar"),
                    ),
                    NotEqual(colUserName, "baz"),
                },
            },
            qs: " WHERE (users.name = ? OR users.name = ?) AND users.name != ?",
            qargs: []interface{}{"foo", "bar", "baz"},
        },
        // Two AND expressions OR'd together
        whereClauseTest{
            c: &whereClause{
                filters: []*Expression{
                    Or(
                        And(
                            NotEqual(colUserName, "foo"),
                            NotEqual(colUserName, "bar"),
                        ),
                        And(
                            NotEqual(colUserName, "baz"),
                            Equal(colUserId, 1),
                        ),
                    ),
                },
            },
            qs: " WHERE ((users.name != ? AND users.name != ?) OR (users.name != ? AND users.id = ?))",
            qargs: []interface{}{"foo", "bar", "baz", 1},
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
