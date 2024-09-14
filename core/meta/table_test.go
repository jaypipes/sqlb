//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta_test

import (
	"testing"

	"github.com/jaypipes/sqlb/core/expr"
	"github.com/jaypipes/sqlb/core/meta"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTableDelete(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")

	q := users.Delete(expr.Equal(colUserName, "foo"))
	b := builder.New()
	expqargs := []interface{}{"foo"}
	expqs := "DELETE FROM users WHERE users.name = ?"
	qs, qargs := b.StringArgs(q)
	assert.Equal(expqargs, qargs)
	assert.Equal(expqs, qs)
}

func TestTableInsert(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")

	tests := []struct {
		name   string
		t      *meta.Table
		values map[string]interface{}
		qs     string
		qargs  []interface{}
		qe     error
	}{
		{
			name:   "Values missing",
			t:      users,
			values: nil,
			qe:     types.NoValues,
		},
		{
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     types.UnknownColumn,
		},
		{
			name:   "Simple INSERT",
			t:      users,
			values: map[string]interface{}{"id": 1},
			qs:     "INSERT INTO users (id) VALUES (?)",
			qargs:  []interface{}{1},
		},
	}
	for _, tt := range tests {
		got, err := tt.t.Insert(tt.values)
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

func TestTableUpdateAll(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")

	tests := []struct {
		name   string
		values map[string]interface{}
		qs     string
		qargs  []interface{}
		qe     error
	}{
		{
			name:   "Values missing",
			values: nil,
			qe:     types.NoValues,
		},
		{
			name:   "Unknown column",
			values: map[string]interface{}{"unknown": 1},
			qe:     types.UnknownColumn,
		},
		{
			name:   "UPDATE no WHERE",
			values: map[string]interface{}{"name": "foo"},
			qs:     "UPDATE users SET name = ?",
			qargs:  []interface{}{"foo"},
		},
	}
	for _, tt := range tests {
		got, err := users.UpdateAll(tt.values)
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

func TestTableUpdate(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")

	values := map[string]interface{}{"name": "bar"}
	q := users.Update(expr.Equal(colUserName, "foo"), values)
	b := builder.New()
	expqargs := []interface{}{"bar", "foo"}
	expqs := "UPDATE users SET name = ? WHERE users.name = ?"
	qs, qargs := b.StringArgs(q)
	assert.Equal(expqargs, qargs)
	assert.Equal(expqs, qs)
}
