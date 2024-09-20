//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package fn

import (
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

// BaseFunction is the struct from which all other projectable functions
// derive. It has base implementations of the non-polymorphic
// `types.Projection` interface methods.
type BaseFunction struct {
	// refs is any Table or DerivedTable that is refs from
	// the aggregate function
	ref types.Relation
	// alias is the aggregate function as an aliased projection
	// (e.g. COUNT(*) AS counter)
	alias string
}

// Name returns the thing's alias, or an empty string if not aliased
func (f *BaseFunction) Name() string {
	return f.alias
}

// Alias returns the thing's alias, or an empty string if not aliased
func (f *BaseFunction) Alias() string {
	return f.alias
}

// AliasOrName returns the thing's alias or its name if not aliased
func (f *BaseFunction) AliasOrName() string {
	return f.alias
}

// References returns a slice of tables or derived tables that are referenced
// by the Projection
func (f *BaseFunction) References() types.Relation {
	return f.ref
}

// ColumnReference returns the object as a `*grammar.ColumnReference`
func (f *BaseFunction) ColumnReference() *grammar.ColumnReference {
	ids := []string{}
	ref := f.References()
	if ref != nil {
		ids = append(ids, ref.AliasOrName())
	}
	ids = append(ids, f.Name())
	cr := &grammar.ColumnReference{
		BasicIdentifierChain: &grammar.IdentifierChain{
			Identifiers: ids,
		},
	}
	return cr
}

// Asc returns a SortSpecification indicating the Function should be used in an
// ORDER BY clause in ASCENDING sort order
func (f *BaseFunction) Asc() grammar.SortSpecification {
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: f.ColumnReference(),
				},
			},
		},
	}
}

// Desc returns a SortSpecification indicating the Column should used in an
// ORDER BY clause in DESCENDING sort order
func (f *BaseFunction) Desc() grammar.SortSpecification {
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: f.ColumnReference(),
				},
			},
		},
		Order: grammar.OrderSpecificationDesc,
	}
}
