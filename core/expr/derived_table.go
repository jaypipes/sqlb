//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr

import (
	"slices"
	"strings"

	"github.com/jaypipes/sqlb/core/meta"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/core/grammar"
)

// NewDerivedTable returns a new DerivedTable from the supplied
// Selection
func NewDerivedTable(
	name string,
	sel types.Relation,
) *DerivedTable {
	cols := map[string]types.Projection{}
	dt := &DerivedTable{
		name: name,
	}
	// We need to project all columns from the supplied Selection's
	// QuerySpecification to the outer QuerySpecification.
	for _, c := range sel.Projections() {
		cname := c.Name()
		outerCol := meta.NewColumn(dt, cname)
		cols[cname] = outerCol
	}
	dt.columns = cols
	dt.qs = sel.QuerySpecification()
	return dt
}

// DerivedTable describes a subquery in the FROM clause of a SELECT statement.
type DerivedTable struct {
	// name is the name of the subquery in the FROM clause
	name string
	// columns is a map of Column structs, keyed by the column's actual name
	// (not alias)
	columns map[string]types.Projection
	// qs is the QuerySpecification the derived table encapsulates
	qs *grammar.QuerySpecification
}

// QuerySpecification returns the object as a `*grammar.QuerySpecification`
func (t *DerivedTable) QuerySpecification() *grammar.QuerySpecification {
	return t.qs
}

// Name returns the name of the DerivedTable
func (t *DerivedTable) Name() string {
	return t.name
}

// Alias returns the alias of the DerivedTable, which is always the
// DerivedTable's name
func (t *DerivedTable) Alias() string {
	return t.name
}

// AliasOrName returns the alias of the DerivedTable
func (t *DerivedTable) AliasOrName() string {
	return t.name
}

// Projections returns a slice of Projection things referenced by the
// Selectable. The slice is sorted by the Projection's name.
func (t *DerivedTable) Projections() []types.Projection {
	cols := make([]types.Projection, 0, len(t.columns))
	for _, c := range t.columns {
		cols = append(cols, c)
	}
	slices.SortFunc(cols, func(a, b types.Projection) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (t *DerivedTable) C(name string) types.Projection {
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

// Column returns a pointer to a Column with a name or alias matching the
// supplied string, or nil if no such column is known
func (t *DerivedTable) Column(name string) types.Projection {
	return t.C(name)
}

// TablePrimary returns the object as a `*grammar.TablePrimary`
func (t *DerivedTable) TablePrimary() *grammar.TablePrimary {
	tp := &grammar.TablePrimary{
		DerivedTable: &grammar.DerivedTable{
			Subquery: grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: grammar.QueryExpressionBody{
						NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
							NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									SimpleTable: &grammar.SimpleTable{
										QuerySpecification: t.QuerySpecification(),
									},
								},
							},
						},
					},
				},
			},
		},
		// Derived tables are always named/aliased
		Correlation: &grammar.Correlation{
			Name: t.Name(),
		},
	}
	return tp
}

// TableReference returns the object as a `*grammar.TableReference`
func (t *DerivedTable) TableReference() *grammar.TableReference {
	return &grammar.TableReference{
		Primary: t.TablePrimary(),
	}
}
