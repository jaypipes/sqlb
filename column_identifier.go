//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type ColumnIdentifier struct {
	tbl   *TableIdentifier
	alias string
	name  string
}

func (c *ColumnIdentifier) From() types.Selection {
	return c.tbl
}

func (c *ColumnIdentifier) DisableAliasScan() func() {
	origAlias := c.alias
	c.alias = ""
	return func() { c.alias = origAlias }
}

func (c *ColumnIdentifier) ArgCount() int {
	return 0
}

func (c *ColumnIdentifier) Size(scanner types.Scanner) int {
	size := 0
	if c.tbl.alias != "" {
		size += len(c.tbl.alias)
	} else {
		size += len(c.tbl.name)
	}
	size += len(grammar.Symbols[grammar.SYM_PERIOD])
	size += len(c.name)
	if c.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(c.alias)
	}
	return size
}

func (c *ColumnIdentifier) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if c.tbl.alias != "" {
		bw += copy(b[bw:], c.tbl.alias)
	} else {
		bw += copy(b[bw:], c.tbl.name)
	}
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_PERIOD])
	bw += copy(b[bw:], c.name)
	if c.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], c.alias)
	}
	return bw
}

func (c *ColumnIdentifier) setAlias(alias string) {
	c.alias = alias
}

func (c *ColumnIdentifier) As(alias string) *ColumnIdentifier {
	return &ColumnIdentifier{
		alias: alias,
		name:  c.name,
		tbl:   c.tbl,
	}
}

func (c *ColumnIdentifier) Desc() *ast.SortColumn {
	return ast.NewSortColumn(c, true)
}

func (c *ColumnIdentifier) Asc() *ast.SortColumn {
	return ast.NewSortColumn(c, false)
}

func (c *ColumnIdentifier) Reverse() *ast.Function {
	return ast.Reverse(c)
}

func (c *ColumnIdentifier) Ascii() *ast.Function {
	return ast.Ascii(c)
}

func (c *ColumnIdentifier) Max() *ast.Function {
	return ast.Max(c)
}

func (c *ColumnIdentifier) Min() *ast.Function {
	return ast.Min(c)
}

func (c *ColumnIdentifier) Sum() *ast.Function {
	return ast.Sum(c)
}

func (c *ColumnIdentifier) Avg() *ast.Function {
	return ast.Avg(c)
}

func (c *ColumnIdentifier) CharLength() *ast.Function {
	return ast.CharLength(c)
}

func (c *ColumnIdentifier) BitLength() *ast.Function {
	return ast.BitLength(c)
}
