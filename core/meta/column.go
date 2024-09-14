//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package meta

import (
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

func NewColumn(
	t types.Relation,
	name string,
) *Column {
	return &Column{
		t:    t,
		name: name,
	}
}

// Column describes a column in a Table
type Column struct {
	// Table is a pointer to the Table or DerivedTable housing this Column
	t types.Relation
	// Name is the name of the Column in the Table
	name string
	// Alias is an optional alias for the column (when a user uses the As()
	// method to alias a column in a SELECT statement)
	alias string
}

// Name returns the true name of the Column (not any alias)
func (c *Column) Name() string {
	return c.name
}

// As returns a copy of the Column aliased as the supplied name
func (c *Column) As(
	alias string,
) types.Projection {
	return &Column{
		t:     c.t,
		name:  c.name,
		alias: alias,
	}
}

// References returns a slice of tables or derived tables that are referenced
// by the Projection
func (c *Column) References() types.Relation {
	return c.t
}

// ColumnReference returns the object as a `*grammar.ColumnReference`
func (c *Column) ColumnReference() *grammar.ColumnReference {
	cr := &grammar.ColumnReference{
		BasicIdentifierChain: &grammar.IdentifierChain{
			Identifiers: []string{c.t.AliasOrName(), c.name},
		},
	}
	if c.alias != "" {
		cr.Correlation = &grammar.Correlation{
			Name: c.alias,
		}
	}
	return cr
}

// DerivedColumn returns the object as a `*grammar.DerivedColumn`
func (c *Column) DerivedColumn() *grammar.DerivedColumn {
	dc := &grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: c.ColumnReference(),
				},
			},
		},
	}
	return dc
}

// Asc returns a SortSpecification indicating the Column should used in an
// ORDER BY clause in ASCENDING sort order
func (c *Column) Asc() grammar.SortSpecification {
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: c.ColumnReference(),
				},
			},
		},
	}
}

// Desc returns a SortSpecification indicating the Column should used in an
// ORDER BY clause in DESCENDING sort order
func (c *Column) Desc() grammar.SortSpecification {
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: c.ColumnReference(),
				},
			},
		},
		Order: grammar.OrderSpecificationDesc,
	}
}

/*
func (c *Column) Reverse() api.Projection {
	return function.Reverse(c)
}

func (c *Column) Ascii() api.Projection {
	return function.Ascii(c)
}

func (c *Column) Max() api.Projection {
	return function.Max(c)
}

func (c *Column) Min() api.Projection {
	return function.Min(c)
}

func (c *Column) Sum() api.Projection {
	return function.Sum(c)
}

func (c *Column) Avg() api.Projection {
	return function.Avg(c)
}

func (c *Column) CharLength() api.Projection {
	return function.CharLength(c)
}

func (c *Column) BitLength() api.Projection {
	return function.BitLength(c)
}

func (c *Column) Trim() api.Projection {
	f := function.Trim(c)
	return f
}

func (c *Column) LTrim() api.Projection {
	f := function.LTrim(c)
	return f
}

func (c *Column) RTrim() api.Projection {
	f := function.RTrim(c)
	return f
}

func (c *Column) TrimChars(chars string) api.Projection {
	f := function.TrimChars(c, chars)
	return f
}

func (c *Column) LTrimChars(chars string) api.Projection {
	f := function.LTrimChars(c, chars)
	return f
}

func (c *Column) RTrimChars(chars string) api.Projection {
	f := function.RTrimChars(c, chars)
	return f
}
*/
