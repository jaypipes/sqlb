//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
	"github.com/jaypipes/sqlb/meta"
)

type Column struct {
	tbl   *Table
	Alias string
	Name  string
}

// Meta returns a pointer to the underlying Meta
func (c *Column) Meta() *meta.Meta {
	return c.tbl.Meta()
}

func (c *Column) From() builder.Selection {
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

func (c *Column) Size(b *builder.Builder) int {
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

func (c *Column) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	if c.tbl.Alias != "" {
		b.WriteString(c.tbl.Alias)
	} else {
		b.WriteString(c.tbl.Name)
	}
	b.Write(grammar.Symbols[grammar.SYM_PERIOD])
	b.WriteString(c.Name)
	if c.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(c.Alias)
	}
}

func (c *Column) As(alias string) builder.Projection {
	return &Column{
		Alias: alias,
		Name:  c.Name,
		tbl:   c.tbl,
	}
}

func (c *Column) Desc() builder.Sortable {
	return sortcolumn.NewDesc(c)
}

func (c *Column) Asc() builder.Sortable {
	return sortcolumn.NewAsc(c)
}

func (c *Column) Reverse() builder.Projection {
	return function.Reverse(c)
}

func (c *Column) Ascii() builder.Projection {
	return function.Ascii(c)
}

func (c *Column) Max() builder.Projection {
	return function.Max(c)
}

func (c *Column) Min() builder.Projection {
	return function.Min(c)
}

func (c *Column) Sum() builder.Projection {
	return function.Sum(c)
}

func (c *Column) Avg() builder.Projection {
	return function.Avg(c)
}

func (c *Column) CharLength() builder.Projection {
	return function.CharLength(c)
}

func (c *Column) BitLength() builder.Projection {
	return function.BitLength(c)
}

func (c *Column) Trim() builder.Projection {
	f := function.Trim(c)
	return f
}

func (c *Column) LTrim() builder.Projection {
	f := function.LTrim(c)
	return f
}

func (c *Column) RTrim() builder.Projection {
	f := function.RTrim(c)
	return f
}

func (c *Column) TrimChars(chars string) builder.Projection {
	f := function.TrimChars(c, chars)
	return f
}

func (c *Column) LTrimChars(chars string) builder.Projection {
	f := function.LTrimChars(c, chars)
	return f
}

func (c *Column) RTrimChars(chars string) builder.Projection {
	f := function.RTrimChars(c, chars)
	return f
}
