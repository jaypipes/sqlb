//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar"
	pkgscanner "github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

// UPDATE <table> SET <column_value_list>[ WHERE <predicates>]
type Update struct {
	table   *ast.TableIdentifier
	columns []*ast.ColumnIdentifier
	values  []interface{}
	where   *ast.WhereClause
}

func (s *Update) ArgCount() int {
	argc := len(s.values)
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	return argc
}

func (s *Update) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_UPDATE]) + len(s.table.Name) + len(grammar.Symbols[grammar.SYM_SET])
	ncols := len(s.columns)
	for _, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <columns> element of the INSERT statement
		size += len(c.Name)
	}
	// NOTE(jaypipes): We do not include the length of interpolation markers,
	// since that differs based on the SQL dialect
	size += len(grammar.Symbols[grammar.SYM_EQUAL]) * ncols
	// Two comma-delimited lists of same number of elements (columns and
	// values)
	size += 2 * (len(grammar.Symbols[grammar.SYM_COMMA_WS]) * (ncols - 1)) // the commas...
	if s.where != nil {
		size += s.where.Size(scanner)
	}
	return size
}

func (s *Update) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_UPDATE])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.Name)
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_SET])

	ncols := len(s.columns)
	for x, c := range s.columns {
		// We don't add the table identifier or use an alias when outputting
		// the column names in the <column_value_lists> element of the UPDATE
		// statement
		bw += copy(b[bw:], c.Name)
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_EQUAL])
		bw += pkgscanner.ScanInterpolationMarker(scanner.Dialect(), b[bw:], *curArg)
		args[*curArg] = s.values[x]
		*curArg++
		if x != (ncols - 1) {
			bw += copy(b[bw:], grammar.Symbols[grammar.SYM_COMMA_WS])
		}
	}

	if s.where != nil {
		bw += s.where.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func (s *Update) AddWhere(e *ast.Expression) *Update {
	if s.where == nil {
		s.where = ast.NewWhereClause(e)
		return s
	}
	s.where.AddExpression(e)
	return s
}

// NewUpdate returns a new UpdateStatement struct that scans into an UPDATE SQL
// statement
func NewUpdate(
	table *ast.TableIdentifier,
	columns []*ast.ColumnIdentifier,
	values []interface{},
	where *ast.WhereClause,
) *Update {
	return &Update{
		table:   table,
		columns: columns,
		values:  values,
		where:   where,
	}
}
