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
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	m := testutil.M()
	users := m.T("users")

	tests := []struct {
		name   string
		t      *api.Table
		values map[string]interface{}
		qs     string
		qargs  []interface{}
		qe     error
	}{
		{
			name:   "Values missing",
			t:      users,
			values: nil,
			qe:     api.NoValues,
		},
		{
			name:   "Target table missing",
			t:      nil,
			values: map[string]interface{}{"name": "foo"},
			qe:     api.TableRequired,
		},
		{
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
		},
		{
			name:   "UPDATE no WHERE",
			t:      users,
			values: map[string]interface{}{"name": "foo"},
			qs:     "UPDATE users SET name = ?",
			qargs:  []interface{}{"foo"},
		},
	}
	for _, tt := range tests {
		got, err := api.Update(tt.t, tt.values)
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
			qe:     api.NoValues,
		},
		{
			name:   "Unknown column",
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
		},
		{
			name:   "UPDATE no WHERE",
			values: map[string]interface{}{"name": "foo"},
			qs:     "UPDATE users SET name = ?",
			qargs:  []interface{}{"foo"},
		},
	}
	for _, tt := range tests {
		got, err := users.Update(tt.values)
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

func TestTableUpdateWhere(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.M()
	users := m.T("users")
	colUserName := users.C("name")

	values := map[string]interface{}{"name": "bar"}
	q, err := users.Update(values)
	require.Nil(err)
	q.Where(api.Equal(colUserName, "foo"))
	b := builder.New()
	expqargs := []interface{}{"bar", "foo"}
	expqs := "UPDATE users SET name = ? WHERE users.name = ?"
	qs, qargs := b.StringArgs(q)
	assert.Equal(expqargs, qargs)
	assert.Equal(expqs, qs)
}
