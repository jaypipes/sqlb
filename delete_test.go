//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDeleteQuery(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")
	colUserName := users.C("name")

	tests := []struct {
		name  string
		q     *DeleteQuery
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "No target table",
			q:    Delete(nil),
			qe:   ERR_DELETE_NO_TARGET,
		},
		{
			name: "DELETE all rows",
			q:    Delete(users),
			qs:   "DELETE FROM users",
		},
		{
			name: "Table.Delete() variant",
			q:    users.Delete(),
			qs:   "DELETE FROM users",
		},
		{
			name:  "DELETE simple WHERE",
			q:     Delete(users).Where(ast.Equal(colUserName, "foo")),
			qs:    "DELETE FROM users WHERE users.name = ?",
			qargs: []interface{}{"foo"},
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
