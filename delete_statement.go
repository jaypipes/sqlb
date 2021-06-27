//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

// DELETE FROM <table> WHERE <predicates>

type deleteStatement struct {
	table *Table
	where *whereClause
}

func (s *deleteStatement) ArgCount() int {
	argc := 0
	if s.where != nil {
		argc += s.where.ArgCount()
	}
	return argc
}

func (s *deleteStatement) Size(scanner types.Scanner) int {
	size := len(grammar.Symbols[grammar.SYM_DELETE]) + len(s.table.name)
	if s.where != nil {
		size += s.where.Size(scanner)
	}
	return size
}

func (s *deleteStatement) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], grammar.Symbols[grammar.SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.name)
	if s.where != nil {
		bw += s.where.Scan(scanner, b[bw:], args, curArg)
	}
	return bw
}

func (s *deleteStatement) addWhere(e *Expression) *deleteStatement {
	if s.where == nil {
		s.where = &whereClause{filters: make([]*Expression, 0)}
	}
	s.where.filters = append(s.where.filters, e)
	return s
}
