//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package identifier

import (
	"strings"

	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
	"github.com/jaypipes/sqlb/internal/scanner"
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

func (c *Column) From() scanner.Selection {
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

func (c *Column) Size(s *scanner.Scanner) int {
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

func (c *Column) Scan(s *scanner.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
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

func (c *Column) As(alias string) scanner.Projection {
	return &Column{
		Alias: alias,
		Name:  c.Name,
		tbl:   c.tbl,
	}
}

func (c *Column) Desc() scanner.Sortable {
	return sortcolumn.NewDesc(c)
}

func (c *Column) Asc() scanner.Sortable {
	return sortcolumn.NewAsc(c)
}

func (c *Column) Reverse() scanner.Projection {
	return function.Reverse(c)
}

func (c *Column) Ascii() scanner.Projection {
	return function.Ascii(c)
}

func (c *Column) Max() scanner.Projection {
	return function.Max(c)
}

func (c *Column) Min() scanner.Projection {
	return function.Min(c)
}

func (c *Column) Sum() scanner.Projection {
	return function.Sum(c)
}

func (c *Column) Avg() scanner.Projection {
	return function.Avg(c)
}

func (c *Column) CharLength() scanner.Projection {
	return function.CharLength(c)
}

func (c *Column) BitLength() scanner.Projection {
	return function.BitLength(c)
}

func (c *Column) Trim() scanner.Projection {
	f := function.Trim(c)
	return f
}

func (c *Column) LTrim() scanner.Projection {
	f := function.LTrim(c)
	return f
}

func (c *Column) RTrim() scanner.Projection {
	f := function.RTrim(c)
	return f
}

func (c *Column) TrimChars(chars string) scanner.Projection {
	f := function.TrimChars(c, chars)
	return f
}

func (c *Column) LTrimChars(chars string) scanner.Projection {
	f := function.LTrimChars(c, chars)
	return f
}

func (c *Column) RTrimChars(chars string) scanner.Projection {
	f := function.RTrimChars(c, chars)
	return f
}
