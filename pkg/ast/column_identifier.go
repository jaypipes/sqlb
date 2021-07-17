//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/schema"
	"github.com/jaypipes/sqlb/pkg/types"
)

type ColumnIdentifier struct {
	tbl   *TableIdentifier
	Alias string
	Name  string
}

// Schema returns a pointer to the underlying Schema
func (c *ColumnIdentifier) Schema() *schema.Schema {
	return c.tbl.Schema()
}

func (c *ColumnIdentifier) From() types.Selection {
	return c.tbl
}

func (c *ColumnIdentifier) DisableAliasScan() func() {
	origAlias := c.Alias
	c.Alias = ""
	return func() { c.Alias = origAlias }
}

func (c *ColumnIdentifier) ArgCount() int {
	return 0
}

func (c *ColumnIdentifier) Size(scanner types.Scanner) int {
	size := 0
	if c.tbl.Alias != "" {
		size += len(c.tbl.Alias)
	} else {
		size += len(c.tbl.Name)
	}
	size += len(grammar.Symbols[grammar.SYM_PERIOD])
	size += len(c.Name)
	if c.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(c.Alias)
	}
	return size
}

func (c *ColumnIdentifier) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if c.tbl.Alias != "" {
		bw += copy(b[bw:], c.tbl.Alias)
	} else {
		bw += copy(b[bw:], c.tbl.Name)
	}
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_PERIOD])
	bw += copy(b[bw:], c.Name)
	if c.Alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], c.Alias)
	}
	return bw
}

func (c *ColumnIdentifier) As(alias string) *ColumnIdentifier {
	return &ColumnIdentifier{
		Alias: alias,
		Name:  c.Name,
		tbl:   c.tbl,
	}
}

func (c *ColumnIdentifier) Desc() *SortColumn {
	return NewSortColumn(c, true)
}

func (c *ColumnIdentifier) Asc() *SortColumn {
	return NewSortColumn(c, false)
}

func (c *ColumnIdentifier) Reverse() types.Projection {
	return Reverse(c)
}

func (c *ColumnIdentifier) Ascii() types.Projection {
	return Ascii(c)
}

func (c *ColumnIdentifier) Max() types.Projection {
	return Max(c)
}

func (c *ColumnIdentifier) Min() types.Projection {
	return Min(c)
}

func (c *ColumnIdentifier) Sum() types.Projection {
	return Sum(c)
}

func (c *ColumnIdentifier) Avg() types.Projection {
	return Avg(c)
}

func (c *ColumnIdentifier) CharLength() types.Projection {
	return CharLength(c)
}

func (c *ColumnIdentifier) BitLength() types.Projection {
	return BitLength(c)
}

func (c *ColumnIdentifier) Trim() types.Projection {
	f := Trim(c)
	return f
}

func (c *ColumnIdentifier) LTrim() types.Projection {
	f := LTrim(c)
	return f
}

func (c *ColumnIdentifier) RTrim() types.Projection {
	f := LTrim(c)
	return f
}

func (c *ColumnIdentifier) TrimChars(chars string) types.Projection {
	f := TrimChars(c, chars)
	return f
}

func (c *ColumnIdentifier) LTrimChars(chars string) types.Projection {
	f := LTrimChars(c, chars)
	return f
}

func (c *ColumnIdentifier) RTrimChars(chars string) types.Projection {
	f := RTrimChars(c, chars)
	return f
}
