package sqlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateQuery(t *testing.T) {
	assert := assert.New(t)

	m := testFixtureMeta()
	users := m.Table("users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		q     *UpdateQuery
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "Values missing",
			q:    Update(users, nil),
			qe:   ERR_UPDATE_NO_VALUES,
		},
		{
			name: "Target table missing",
			q:    Update(nil, map[string]interface{}{"name": "foo"}),
			qe:   ERR_UPDATE_NO_TARGET,
		},
		{
			name: "Unknown column",
			q:    Update(users, map[string]interface{}{"unknown": 1}),
			qe:   ERR_UPDATE_UNKNOWN_COLUMN,
		},
		{
			name:  "UPDATE no WHERE",
			q:     Update(users, map[string]interface{}{"name": "foo"}),
			qs:    "UPDATE users SET name = ?",
			qargs: []interface{}{"foo"},
		},
		{
			name: "UPDATE simple WHERE",
			q: Update(users, map[string]interface{}{"name": "bar"}).Where(
				Equal(colUserName, "foo"),
			),
			qs:    "UPDATE users SET name = ? WHERE users.name = ?",
			qargs: []interface{}{"bar", "foo"},
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
