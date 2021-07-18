//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"testing"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInsertQuery(t *testing.T) {
	assert := assert.New(t)

	sc := testutil.Schema()
	users := T(sc, "users")

	tests := []struct {
		name  string
		q     *InsertQuery
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "Values missing",
			q:    Insert(users, nil),
			qe:   ERR_INSERT_NO_VALUES,
		},
		{
			name: "Unknown column",
			q:    Insert(users, map[string]interface{}{"unknown": 1}),
			qe:   ERR_INSERT_UNKNOWN_COLUMN,
		},
		{
			name:  "Simple INSERT",
			q:     Insert(users, map[string]interface{}{"id": 1}),
			qs:    "INSERT INTO users (id) VALUES (?)",
			qargs: []interface{}{1},
		},
		//{
		//	name:  "INSERT using Table.Insert() adapter",
		//	q:     users.Insert(map[string]interface{}{"id": 1}),
		//	qs:    "INSERT INTO users (id) VALUES (?)",
		//	qargs: []interface{}{1},
		//},
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
