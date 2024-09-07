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

func TestInsert(t *testing.T) {
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
			name:   "Table missing",
			t:      nil,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.TableRequired,
		},
		{
			name:   "Values missing",
			t:      users,
			values: nil,
			qe:     api.NoValues,
		},
		{
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
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
		got, err := api.Insert(tt.t, tt.values)
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
func TestTableInsert(t *testing.T) {
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
			name:   "Unknown column",
			t:      users,
			values: map[string]interface{}{"unknown": 1},
			qe:     api.UnknownColumn,
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
