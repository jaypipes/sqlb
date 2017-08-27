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

    tests := []orderByTest{
        // column def asc
        orderByTest{
            c: &orderByClause{
                cols: &List{
                    elements: []element{
                        colUserName.Asc(),
                    },
                },
            },
            qs: " ORDER BY users.name",
        },
        // column asc
        orderByTest{
            c: &orderByClause{
                cols: &List{
                    elements: []element{
                        colUserName.Column().Asc(),
                    },
                },
            },
            qs: " ORDER BY users.name",
        },
        // column def desc
        orderByTest{
            c: &orderByClause{
                cols: &List{
                    elements: []element{
                        colUserName.Desc(),
                    },
                },
            },
            qs: " ORDER BY users.name DESC",
        },
        // column desc
        orderByTest{
            c: &orderByClause{
                cols: &List{
                    elements: []element{
                        colUserName.Column().Desc(),
                    },
                },
            },
            qs: " ORDER BY users.name DESC",
        },
        // multi column mixed
        orderByTest{
            c: &orderByClause{
                cols: &List{
                    elements: []element{
                        colUserName.Column(),
                        colUserId.Desc(),
                    },
                },
            },
            qs: " ORDER BY users.name, users.id DESC",
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
