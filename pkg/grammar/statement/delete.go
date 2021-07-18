//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package statement

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/clause"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/grammar/identifier"
	"github.com/jaypipes/sqlb/pkg/types"
)

// DELETE FROM <table> WHERE <predicates>
type Delete struct {
	table *identifier.Table
	where *clause.Where
}

func (s *Delete) ArgCount() int {
	argc := 0
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	return argc
}

func (s *Delete) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_DELETE]) + len(s.table.Name)
	if s.where != nil {
		size += s.where.Size(scanner)
	}
	return size
}

func (s *Delete) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.Name)
	if s.where != nil {
		bw += s.where.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func (s *Delete) AddWhere(e *expression.Expression) *Delete {
	if s.where == nil {
		s.where = clause.NewWhere(e)
		return s
	}
	s.where.AddExpression(e)
	return s
}

// NewDelete returns a new DeleteStatement struct that scans into a DELETE SQL
// statement
func NewDelete(
	table *identifier.Table,
	where *clause.Where,
) *Delete {
	return &Delete{
		table: table,
		where: where,
	}
}
