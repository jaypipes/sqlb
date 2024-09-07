//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api_test

import (
	"testing"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")

	tests := []struct {
		name  string
		t     *api.Table
		qs    string
		qargs []interface{}
		qe    error
	}{
		{
			name: "No target table",
			qe:   api.TableRequired,
		},
		{
			name: "DELETE all rows",
			t:    users,
			qs:   "DELETE FROM users",
		},
	}
	for _, tt := range tests {
		got, err := api.Delete(tt.t)
		if tt.qe != nil {
			assert.Equal(tt.qe, err)
			continue
		}
		assert.Nil(err)
		b := builder.New()
		qs, qargs := b.StringArgs(got)
		assert.Equal(len(tt.qargs), len(qargs))
		assert.Equal(tt.qs, qs)
	}
}

func TestTableDeleteWhere(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")

	q := users.Delete()
	q.Where(api.Equal(colUserName, "foo"))
	b := builder.New()
	expqargs := []interface{}{"foo"}
	expqs := "DELETE FROM users WHERE users.name = ?"
	qs, qargs := b.StringArgs(q)
	assert.Equal(expqargs, qargs)
	assert.Equal(expqs, qs)
}
