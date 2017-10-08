package sqlb

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDeleteQuery(t *testing.T) {
    assert := assert.New(t)

    m := testFixtureMeta()
    users := m.Table("users")

    tests := []struct{
        name string
        q *DeleteQuery
        qs string
        qargs []interface{}
        qe error
    }{
        {
            name: "No target table",
            q: Delete(nil),
            qe: ERR_DELETE_NO_TARGET,
        },
        {
            name: "DELETE all rows",
            q: Delete(users),
            qs: "DELETE FROM users",
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
