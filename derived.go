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

// DerivedTable is a SELECT in the FROM clause. It is always aliased and the
// projections for a derived table take this alias as their selection alias.
//
// The projections of a derived table are not the same as the projections for
// the SELECT that is being wrapped. For example, given the following SQL
// statement:
//
// SELECT u.id, u.name FROM (
//   SELECT users.id, users.name FROM users
// ) AS u
//
// The inner SELECT's projections are columns from the users Table or TableDef.
// However, the derived table's projections are separate and include the alias
// of the derived table as the selection alias (u instead of users).
type DerivedTable struct {
	alias string
	from  *SelectStatement
}

// DerivedColumns returns a collection of DerivedColumn projections that have
// been constructed to refer to this derived table and not have any outer alias
func (dt *DerivedTable) DerivedColumns() []types.Projection {
	nprojs := len(dt.from.projs)
	projs := make([]types.Projection, nprojs)
	for x := 0; x < nprojs; x++ {
		p := dt.from.projs[x]
		switch p.(type) {
		case *ast.ColumnIdentifier:
			projs[x] = &DerivedColumn{dt: dt, c: p.(*ast.ColumnIdentifier)}
		}
	}
	return projs
}

func (dt *DerivedTable) Projections() []types.Projection {
	nprojs := len(dt.from.projs)
	projs := make([]types.Projection, nprojs)
	for x := 0; x < nprojs; x++ {
		p := dt.from.projs[x]
		switch p.(type) {
		case *ast.ColumnIdentifier:
		}
	}
	return projs
}

func (dt *DerivedTable) ArgCount() int {
	return dt.from.ArgCount()
}

func (dt *DerivedTable) Size(scanner types.Scanner) int {
	size := dt.from.Size(scanner)
	size += (len(grammar.Symbols[grammar.SYM_LPAREN]) + len(grammar.Symbols[grammar.SYM_RPAREN]) +
		len(grammar.Symbols[grammar.SYM_AS]) + len(dt.alias))
	return size
}

func (dt *DerivedTable) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_LPAREN])
	bw += dt.from.Scan(scanner, b[bw:], args, curArg)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_RPAREN])
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
	bw += copy(b[bw:], dt.alias)
	return bw
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
//   SELECT users.id, users.name FROM users
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
//   SELECT
//     users.id AS user_id,
//     users.name AS user_name
//   FROM users
// ) AS u
//
// <outer> should instead contain:
//
// &DerivedColumn{dt: dt, c: &Column{alias: "user_id". name: "id", tbl: users}},
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
//   SELECT users.name AS user_name
//   FROM users
// ) AS u
type DerivedColumn struct {
	alias string // This is the outermost alias
	c     *ast.ColumnIdentifier
	dt    *DerivedTable
}

func (dc *DerivedColumn) From() types.Selection {
	return dc.dt
}

func (dc *DerivedColumn) DisableAliasScan() func() {
	origAlias := dc.alias
	dc.alias = ""
	return func() { dc.alias = origAlias }
}

func (dc *DerivedColumn) ArgCount() int {
	return 0
}

func (dc *DerivedColumn) Size(scanner types.Scanner) int {
	size := len(dc.dt.alias)
	size += len(grammar.Symbols[grammar.SYM_PERIOD])
	if dc.c.Alias != "" {
		size += len(dc.c.Alias)
	} else {
		size += len(dc.c.Name)
	}
	if dc.alias != "" {
		size += len(grammar.Symbols[grammar.SYM_AS]) + len(dc.alias)
	}
	return size
}

func (dc *DerivedColumn) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], dc.dt.alias)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_PERIOD])
	if dc.c.Alias != "" {
		bw += copy(b[bw:], dc.c.Alias)
	} else {
		bw += copy(b[bw:], dc.c.Name)
	}
	if dc.alias != "" {
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AS])
		bw += copy(b[bw:], dc.alias)
	}
	return bw
}
