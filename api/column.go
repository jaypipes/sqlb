//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import "github.com/jaypipes/sqlb/grammar"

// Column describes a column in a Table
type Column struct {
	// Table is a pointer to the Table or DerivedTable housing this Column
	t interface{}
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
) *Column {
	return &Column{
		t:     c.t,
		name:  c.name,
		alias: alias,
	}
}

// TableName returns the name of the underlying Table or DerivedTable the
// Column belongs to. If the underlying table is a true Table (not a
// DerivedTable) and has an Alias, the Alias is returned instead of the Table's
// Name.
func (c *Column) TableName() string {
	t, ok := c.t.(*Table)
	if ok {
		if t.alias != "" {
			return t.alias
		}
		return t.name
	}
	dt := c.t.(*DerivedTable)
	return dt.name
}

// TableNameNoAlias returns the name of the underlying Table or DerivedTable the
// Column belongs to. If the underlying table is a true Table (not a
// DerivedTable) and has an Alias, it is NOT returned. Instead, the Table's
// real Name is returned.
func (c *Column) TableNameNoAlias() string {
	t, ok := c.t.(*Table)
	if ok {
		return t.name
	}
	dt := c.t.(*DerivedTable)
	return dt.name
}

// TableAlias returns the alias of the underlying Table. If the underlying table is
// a DerivedTable, always returns a nullstring.
func (c *Column) TableAlias() string {
	t, ok := c.t.(*Table)
	if ok {
		return t.alias
	}
	return ""
}

// Asc returns a SortSpecification indicating the Column should used in an
// ORDER BY clause in ASCENDING sort order
func (c *Column) Asc() grammar.SortSpecification {
	cr := &grammar.ColumnReference{
		BasicIdentifierChain: &grammar.IdentifierChain{
			Identifiers: []string{
				c.TableName(), c.name,
			},
		},
	}
	if c.alias != "" {
		cr.Correlation = &grammar.Correlation{
			Name: c.alias,
		}
	}
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: cr,
				},
			},
		},
	}
}

// Desc returns a SortSpecification indicating the Column should used in an
// ORDER BY clause in DESCENDING sort order
func (c *Column) Desc() grammar.SortSpecification {
	cr := &grammar.ColumnReference{
		BasicIdentifierChain: &grammar.IdentifierChain{
			Identifiers: []string{
				c.TableName(), c.name,
			},
		},
	}
	if c.alias != "" {
		cr.Correlation = &grammar.Correlation{
			Name: c.alias,
		}
	}
	return grammar.SortSpecification{
		Key: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					ColumnReference: cr,
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
