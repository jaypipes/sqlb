package sqlb

// DELETE FROM <table> WHERE <predicates>

type deleteStatement struct {
	table *Table
	where *whereClause
}

func (s *deleteStatement) argCount() int {
	argc := 0
	if s.where != nil {
		argc += s.where.argCount()
	}
	return argc
}

func (s *deleteStatement) size(scanner *sqlScanner) int {
	size := len(Symbols[SYM_DELETE]) + len(s.table.name)
	if s.where != nil {
		size += s.where.size(scanner)
	}
	return size
}

func (s *deleteStatement) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], Symbols[SYM_DELETE])
	// We don't add any table alias when outputting the table identifier
	bw += copy(b[bw:], s.table.name)
	if s.where != nil {
		bw += s.where.scan(scanner, b[bw:], args, curArg)
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
