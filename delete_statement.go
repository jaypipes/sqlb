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

func (s *deleteStatement) size() int {
    size := len(Symbols[SYM_DELETE]) + len(s.table.name)
    if s.where != nil {
        size += s.where.size()
    }
    return size
}

func (s *deleteStatement) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_DELETE])
    // We don't add any table alias when outputting the table identifier
    bw += copy(b[bw:], s.table.name)
    if s.where != nil {
        wbw, wac := s.where.scan(b[bw:], args[ac:])
        bw += wbw
        ac += wac
    }
    return bw, ac
}

func (s *deleteStatement) addWhere(e *Expression) *deleteStatement {
    if s.where == nil {
        s.where = &whereClause{filters: make([]*Expression, 0)}
    }
    s.where.filters = append(s.where.filters, e)
    return s
}
