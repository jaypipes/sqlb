package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestInsertStatement(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")
    colUserName := users.C("name")
    colUserId := users.C("id")

    tests := []struct{
        name string
        s *insertStatement
        qs string
        qargs []interface{}
    }{
        {
            name: "Simple INSERT",
            s: &insertStatement{
                table: users,
                columns: []*Column{colUserId, colUserName},
                values: []interface{}{nil, "foo"},
            },
            qs: "INSERT INTO users (id, name) VALUES (?, ?)",
            qargs: []interface{}{nil, "foo"},
        },
        {
            name: "Ensure no aliasing in table names",
            s: &insertStatement{
                table: users.As("u"),
                columns: []*Column{colUserId, colUserName},
                values: []interface{}{nil, "foo"},
            },
            qs: "INSERT INTO users (id, name) VALUES (?, ?)",
            qargs: []interface{}{nil, "foo"},
        },
        {
            name: "Ensure no aliasing in column names",
            s: &insertStatement{
                table: users,
                columns: []*Column{colUserId.As("user_id"), colUserName},
                values: []interface{}{nil, "foo"},
            },
            qs: "INSERT INTO users (id, name) VALUES (?, ?)",
            qargs: []interface{}{nil, "foo"},
        },
    }
    for _, test := range tests {
        expLen := len(test.qs)
        s := test.s.size()
        assert.Equal(expLen, s)

        expArgc := len(test.qargs)
        assert.Equal(expArgc, test.s.argCount())

        b := make([]byte, s)
        written, _ := test.s.scan(b, test.qargs)

        assert.Equal(written, s)
        assert.Equal(test.qs, string(b))
    }
}
