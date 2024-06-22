//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
)

type Column struct {
	tbl   *Table
	Alias string
	Name  string
}

func (c *Column) From() api.Selection {
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

func (c *Column) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
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
	return b.String()
}

func (c *Column) As(alias string) api.Projection {
	return &Column{
		Alias: alias,
		Name:  c.Name,
		tbl:   c.tbl,
	}
}

func (c *Column) Desc() api.Orderable {
	return sortcolumn.NewDesc(c)
}

func (c *Column) Asc() api.Orderable {
	return sortcolumn.NewAsc(c)
}

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
