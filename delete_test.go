//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/errors"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/scanner"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDeleteQuery(t *testing.T) {
	assert := assert.New(t)

	m := testutil.Meta()
	users := T(m, "users")
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
			qe:   errors.NoTargetTable,
		},
		{
			name: "DELETE all rows",
			q:    Delete(users),
			qs:   "DELETE FROM users",
		},
		//{
		//	name: "Table.Delete() variant",
		//	q:    users.Delete(),
		//	qs:   "DELETE FROM users",
		//},
		{
			name:  "DELETE simple WHERE",
			q:     Delete(users).Where(expression.Equal(colUserName, "foo")),
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
		scan := scanner.DefaultScanner
		qs, qargs := scan.StringArgs(test.q)
		assert.Equal(len(test.qargs), len(qargs))
		assert.Equal(test.qs, qs)
	}
}
