//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"github.com/jaypipes/sqlb/internal/builder"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/sortcolumn"
)

// DerivedTable is a SELECT in the FROM clause. It is always aliased and the
// projections for a derived table take this alias as their selection alias.
//
// The projections of a derived table are not the same as the projections for
// the SELECT that is being wrapped. For example, given the following SQL
// statement:
//
// SELECT u.id, u.name FROM (
//
//	SELECT users.id, users.name FROM users
//
// ) AS u
//
// The inner SELECT's projections are columns from the users Table or TableDef.
// However, the derived table's projections are separate and include the alias
// of the derived table as the selection alias (u instead of users).
type DerivedTable struct {
	Alias string
	from  builder.Selection
}

// DerivedColumns returns a collection of DerivedColumn projections that have
// been constructed to refer to this derived table and not have any outer alias
func (dt *DerivedTable) DerivedColumns() []builder.Projection {
	projs := []builder.Projection{}
	for _, p := range dt.from.Projections() {
		switch p.(type) {
		case *identifier.Column:
			projs = append(projs, &DerivedColumn{dt: dt, c: p.(*identifier.Column)})
		}
	}
	return projs
}

func (dt *DerivedTable) Projections() []builder.Projection {
	return dt.from.Projections()
}

func (dt *DerivedTable) ArgCount() int {
	return dt.from.ArgCount()
}

func (dt *DerivedTable) Size(b *builder.Builder) int {
	size := dt.from.Size(b)
	size += (len(grammar.Symbols[grammar.SYM_LPAREN]) + len(grammar.Symbols[grammar.SYM_RPAREN]) +
		len(grammar.Symbols[grammar.SYM_AS]) + len(dt.Alias))
	return size
}

func (dt *DerivedTable) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.Write(grammar.Symbols[grammar.SYM_LPAREN])
	dt.from.Scan(b, args, curArg)
	b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	b.Write(grammar.Symbols[grammar.SYM_AS])
	b.WriteString(dt.Alias)
}

// NewDerivedTable returns a new DerivedTable struct representing a SELECT in
// the FROM clause.
func NewDerivedTable(
	alias string,
	sel builder.Selection,
) *DerivedTable {
	return &DerivedTable{
		Alias: alias,
		from:  sel,
	}
}

// DerivedColumn is a type of projection that is produced from a derived table
// (SELECT in the FROM clause). What makes a derived column unique is that it
// uses the alias of the underlying column as its name in the outer projection.
//
// The inner projection is a column against an underlying table or table def.
// The outer projection will have the selection alias of the derived table and
// the name of the projection will be the alias or name of the underlying
// column. For example, given the following SQL:
//
// SELECT <outer> FROM (
//
//	SELECT users.id, users.name FROM users
//
// ) AS u
//
// <outer> should contain:
//
// &DerivedColumn{dt: dt, c: &Column{name: "id", tbl: users}},
// &DerivedColumn{dt: dt, c: &Column{name: "name", tbl: users}}
//
// when scanned into <outer>, that should produce:
//
// []byte("u.id, u.name")
//
// However, let's say that the inner projections have been
// aliased, like so:
//
// SELECT <outer> FROM (
//
//	SELECT
//	  users.id AS user_id,
//	  users.name AS user_name
//	FROM users
//
// ) AS u
//
// <outer> should instead contain:
//
// &DerivedColumn{dt: dt, c: &Column{alias: "user_id", name: "id", tbl: users}},
// &DerivedColumn{dt: dt, c: &Column{alias: "user_name", name: "name", tbl: users}}
//
// which, when scanned into <outer>, should produce:
//
// []byte("u.user_id, u.user_name")
//
// Finally, the DerivedColumn can itself have an alias, which can result in the
// outermost projection looking like so:
//
// SELECT u.user_name AS uname FROM (
//
//	SELECT users.name AS user_name
//	FROM users
//
// ) AS u
type DerivedColumn struct {
	Alias string // This is the outermost alias
	c     *identifier.Column
	dt    *DerivedTable
}

func (dc *DerivedColumn) C() *identifier.Column {
	return dc.c
}

func (dc *DerivedColumn) From() builder.Selection {
	return dc.dt
}

func (dc *DerivedColumn) DisableAliasScan() func() {
	origAlias := dc.Alias
	dc.Alias = ""
	return func() { dc.Alias = origAlias }
}

func (dc *DerivedColumn) ArgCount() int {
	return 0
}

func (dc *DerivedColumn) Size(b *builder.Builder) int {
	size := len(dc.dt.Alias)
	size += len(grammar.Symbols[grammar.SYM_PERIOD])
	if dc.c.Alias != "" {
		size += len(dc.c.Alias)
	} else {
		size += len(dc.c.Name)
	}
	if dc.Alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(dc.Alias)
	}
	return size
}

func (dc *DerivedColumn) Scan(b *builder.Builder, args []interface{}, curArg *int) {
	b.WriteString(dc.dt.Alias)
	b.Write(grammar.Symbols[grammar.SYM_PERIOD])
	if dc.c.Alias != "" {
		b.WriteString(dc.c.Alias)
	} else {
		b.WriteString(dc.c.Name)
	}
	if dc.Alias != "" {
		b.Write(grammar.Symbols[grammar.SYM_AS])
		b.WriteString(dc.Alias)
	}
}

func (dc *DerivedColumn) As(alias string) builder.Projection {
	return &DerivedColumn{
		Alias: alias,
		c:     dc.c,
		dt:    dc.dt,
	}
}

func (dc *DerivedColumn) Desc() builder.Sortable {
	return sortcolumn.NewDesc(dc)
}

func (dc *DerivedColumn) Asc() builder.Sortable {
	return sortcolumn.NewAsc(dc)
}
