//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jaypipes/sqlb/core/fn"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/internal/inspect"
)

// NewTable returns a new Table with the supplied properties
func NewTable(
	meta *Meta, // The metadata from the RDBMS
	name string, // The true name of the table
	cnames ...string, // column names
) *Table {
	t := &Table{
		m:    meta,
		name: name,
	}
	cols := make(map[string]types.Projection, len(cnames))
	for _, cname := range cnames {
		c := &Column{
			t:    t,
			name: cname,
		}
		cols[cname] = c
	}
	t.columns = cols
	return t
}

// Table describes metadata about a table in a database.
type Table struct {
	// Meta is a pointer at the metadata collection for the database.
	m *Meta
	// Name is the name of the table in the database.
	name string
	// Columns is a map of Column structs, keyed by the column's actual name
	// (not alias).
	columns map[string]types.Projection
	// Alias is any alias/correlation name given to this Table for use in a
	// SELECT statement
	alias string
}

// Meta returns the metadata associated with the underlying RDBMS
func (t *Table) Meta() *Meta {
	return t.m
}

// AddColumn adds a new Column to the Table. The supplied argument can be
// either a *Column or a string. If the argument is a string, a new Column with
// that name is created. If a same-named Column already existed for the Table,
// it is overwritten with the supplied Column.
func (t *Table) AddColumn(c interface{}) {
	if t.columns == nil {
		t.columns = map[string]types.Projection{}
	}
	if col, ok := c.(*Column); ok {
		t.columns[col.name] = col
	} else {
		cname := c.(string)
		t.columns[cname] = &Column{name: cname, t: t}
	}
}

// Name returns the true name of the Table, no alias
func (t *Table) Name() string {
	return t.name
}

// Alias returns the aliased name of the Table
func (t *Table) Alias() string {
	return t.alias
}

// AliasOrName returns the aliased name of the Table or the real name if the
// Table is not aliased
func (t *Table) AliasOrName() string {
	if t.alias != "" {
		return t.alias
	}
	return t.name
}

// Projections returns a slice of Projection things referenced by the
// Selectable. The slice should be sorted by the Projection's name.
func (t *Table) Projections() []types.Projection {
	cols := make([]types.Projection, 0, len(t.columns))
	for _, c := range t.columns {
		cols = append(cols, c)
	}
	slices.SortFunc(cols, func(a, b types.Projection) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

// As returns a copy of the Table, aliased to the supplied name
func (t *Table) As(alias string) *Table {
	at := &Table{
		m:     t.m,
		name:  t.name,
		alias: alias,
	}
	// Build a copy of the table's columns and point those columns to the new
	// aliased table
	atCols := make(map[string]types.Projection, len(t.columns))
	for k, c := range t.columns {
		atc := &Column{
			name: c.Name(),
			t:    at,
		}
		atCols[k] = atc
	}
	at.columns = atCols
	return at
}

// Column returns a pointer to a Column with a name matching the supplied
// string, or nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *Table) Column(name string) types.Projection {
	if c, ok := t.columns[name]; ok {
		return c
	}
	for _, c := range t.columns {
		if strings.EqualFold(c.Name(), name) {
			return c
		}
	}
	return nil
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *Table) C(name string) types.Projection {
	return t.Column(name)
}

// QuerySpecification returns the object as a `*grammar.QuerySpecification`
func (t *Table) QuerySpecification() *grammar.QuerySpecification {
	sels := []grammar.SelectSublist{}
	for _, c := range t.Projections() {
		sels = append(sels, grammar.SelectSublist{DerivedColumn: c.DerivedColumn()})
	}
	return &grammar.QuerySpecification{
		SelectList: grammar.SelectList{
			Sublists: sels,
		},
		TableExpression: grammar.TableExpression{
			FromClause: grammar.FromClause{
				TableReferences: []grammar.TableReference{*t.TableReference()},
			},
		},
	}
}

// TablePrimary returns the object as a `*grammar.TablePrimary`
func (t *Table) TablePrimary() *grammar.TablePrimary {
	tname := t.Name()
	tp := &grammar.TablePrimary{
		TableName: &tname,
	}
	if t.alias != "" {
		tp.Correlation = &grammar.Correlation{
			Name: t.Alias(),
		}
	}
	return tp
}

// TableReference returns the object as a `*grammar.TableReference`
func (t *Table) TableReference() *grammar.TableReference {
	return &grammar.TableReference{Primary: t.TablePrimary()}
}

// Insert returns an InstanceStatement that produces an INSERT SQL statement
// for the table and map of column name to value for that column to insert,
func (t *Table) Insert(
	values map[string]interface{},
) (*grammar.InsertStatement, error) {
	if t == nil {
		return nil, types.TableRequired
	}
	if len(values) == 0 {
		return nil, types.NoValues
	}

	// Make sure all keys in the map point to actual columns in the target
	// table.
	cols := make([]string, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, types.UnknownColumn
		}
		cols[x] = c.Name()
		vals[x] = v
		x++
	}

	return &grammar.InsertStatement{
		TableName: t.name,
		Columns:   cols,
		Values:    vals,
	}, nil
}

// DeleteAll returns a `*grammar.DeleteStatementSearched` that will produce a
// DELETE SQL statement **with no WHERE clause**.
func (t *Table) DeleteAll() *grammar.DeleteStatementSearched {
	return &grammar.DeleteStatementSearched{
		TableName: t.name,
	}
}

// DeleteWhere returns the `*grammar.DeleteStatementSearched` adapted with a
// supplied search condition. The argument must be coercible into a Boolean
// Value Expression.
func (t *Table) Delete(
	expr interface{},
) *grammar.DeleteStatementSearched {
	bve := inspect.BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	return &grammar.DeleteStatementSearched{
		TableName: t.name,
		WhereClause: &grammar.WhereClause{
			SearchCondition: *bve,
		},
	}
}

// UpdateAll returns a `grammar.UpdateStatementSearched` that will produce an
// UPDATE SQL statement **with no WHERE clause**.
//
// The supplied map of values is keyed by the column name the value will be
// updated to.
func (t *Table) UpdateAll(
	values map[string]interface{},
) (*grammar.UpdateStatementSearched, error) {
	if len(values) == 0 {
		return nil, types.NoValues
	}
	cols := make([]string, len(values))
	vals := make([]interface{}, len(values))
	x := 0
	for k, v := range values {
		c := t.C(k)
		if c == nil {
			return nil, types.UnknownColumn
		}
		cols[x] = k
		vals[x] = v
		x++
	}
	return &grammar.UpdateStatementSearched{
		TableName: t.name,
		Columns:   cols,
		Values:    vals,
	}, nil
}

// Update returns a `*grammar.UpdateStatementSearched` adapted with the
// supplied search condition. The first argument must be coercible into a
// Boolean Value Expression.
func (t *Table) Update(
	expr interface{},
	values map[string]interface{},
) *grammar.UpdateStatementSearched {
	bve := inspect.BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	us, err := t.UpdateAll(values)
	if err != nil {
		panic(err)
	}
	us.WhereClause = &grammar.WhereClause{
		SearchCondition: *bve,
	}
	return us
}

// Count returns an AggregateFunction representing a COUNT(*) or a
// COUNT(<column>) on the table. The function accepts zero or one arguments. If
// no arguments are passed, the result is an AggregateFunction that produces a
// COUNT(*) against the table. If one argument is supplied, it should be either
// a string name of a column in the table, a Column in the table or something
// that can be coerced into a RowValueExpression.
func (t *Table) Count(args ...interface{}) types.Projection {
	if len(args) > 1 {
		panic("Count expects either zero or one argument")
	}
	if len(args) == 0 {
		return fn.CountStar(t)
	}
	arg := args[0]
	switch arg := arg.(type) {
	case string:
		c := t.C(arg)
		if c == nil {
			msg := fmt.Sprintf(
				"attempted Count() on unknown column %s",
				arg,
			)
			panic(msg)
		}
		return fn.Aggregate(c, grammar.ComputationalOperationCount)
	}
	return fn.Aggregate(arg, grammar.ComputationalOperationCount)
}
