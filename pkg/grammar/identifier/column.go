//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/function"
	"github.com/jaypipes/sqlb/pkg/grammar/sortcolumn"
	"github.com/jaypipes/sqlb/pkg/schema"
	"github.com/jaypipes/sqlb/pkg/types"
)

type Column struct {
	tbl   *Table
	Alias string
	Name  string
}

// Schema returns a pointer to the underlying Schema
func (c *Column) Schema() *schema.Schema {
	return c.tbl.Schema()
}

func (c *Column) From() types.Selection {
	return c.tbl
}

func (c *Column) DisableAliasScan() func() {
	origAlias := c.Alias
	c.Alias = ""
	return func() { c.Alias = origAlias }
}

func (c *Column) ArgCount() int {
	return 0
}

func (c *Column) Size(scanner types.Scanner) int {
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

func (c *Column) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
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

func (c *Column) As(alias string) types.Projection {
	return &Column{
		Alias: alias,
		Name:  c.Name,
		tbl:   c.tbl,
	}
}

func (c *Column) Desc() types.Sortable {
	return sortcolumn.NewDesc(c)
}

func (c *Column) Asc() types.Sortable {
	return sortcolumn.NewAsc(c)
}

func (c *Column) Reverse() types.Projection {
	return function.Reverse(c)
}

func (c *Column) Ascii() types.Projection {
	return function.Ascii(c)
}

func (c *Column) Max() types.Projection {
	return function.Max(c)
}

func (c *Column) Min() types.Projection {
	return function.Min(c)
}

func (c *Column) Sum() types.Projection {
	return function.Sum(c)
}

func (c *Column) Avg() types.Projection {
	return function.Avg(c)
}

func (c *Column) CharLength() types.Projection {
	return function.CharLength(c)
}

func (c *Column) BitLength() types.Projection {
	return function.BitLength(c)
}

func (c *Column) Trim() types.Projection {
	f := function.Trim(c)
	return f
}

func (c *Column) LTrim() types.Projection {
	f := function.LTrim(c)
	return f
}

func (c *Column) RTrim() types.Projection {
	f := function.RTrim(c)
	return f
}

func (c *Column) TrimChars(chars string) types.Projection {
	f := function.TrimChars(c, chars)
	return f
}

func (c *Column) LTrimChars(chars string) types.Projection {
	f := function.LTrimChars(c, chars)
	return f
}

func (c *Column) RTrimChars(chars string) types.Projection {
	f := function.RTrimChars(c, chars)
	return f
}
