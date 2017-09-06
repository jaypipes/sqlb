package sqlb

type selectClause struct {
    projs []projection
    selections []selection
    joins []*joinClause
    where *whereClause
    groupBy *groupByClause
    orderBy *orderByClause
    limit *limitClause
}

func (s *selectClause) argCount() int {
    argc := 0
    for _, p := range s.projs {
        argc += p.argCount()
    }
    for _, sel := range s.selections {
        argc += sel.argCount()
    }
    for _, join := range s.joins {
        argc += join.argCount()
    }
    if s.where != nil {
        argc += s.where.argCount()
    }
    if s.groupBy != nil {
        argc += s.groupBy.argCount()
    }
    if s.orderBy != nil {
        argc += s.orderBy.argCount()
    }
    if s.limit != nil {
        argc += s.limit.argCount()
    }
    return argc
}

func (s *selectClause) size() int {
    size := len(Symbols[SYM_SELECT])
    nprojs := len(s.projs)
    for _, p := range s.projs {
        size += p.size()
    }
    size += (len(Symbols[SYM_COMMA_WS]) * (nprojs - 1))  // the commas...
    nsels := len(s.selections)
    if nsels > 0 {
        size += len(Symbols[SYM_FROM])
        for _, sel := range s.selections {
            size += sel.size()
        }
        size += (len(Symbols[SYM_COMMA_WS]) * (nsels - 1))  // the commas...
        for _, join := range s.joins {
            size += join.size()
        }
    }
    if s.where != nil {
        size += s.where.size()
    }
    if s.groupBy != nil {
        size += s.groupBy.size()
    }
    if s.orderBy != nil {
        size += s.orderBy.size()
    }
    if s.limit != nil {
        size += s.limit.size()
    }
    return size
}

func (s *selectClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_SELECT])
    nprojs := len(s.projs)
    for x, p := range s.projs {
        pbw, pac := p.scan(b[bw:], args[ac:])
        bw += pbw
        ac += pac
        if x != (nprojs - 1) {
            bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
        }
    }
    nsels := len(s.selections)
    if nsels > 0 {
        bw += copy(b[bw:], Symbols[SYM_FROM])
        for x, sel := range s.selections {
            sbw, sac := sel.scan(b[bw:], args)
            bw += sbw
            ac += sac
            if x != (nsels - 1) {
                bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
            }
        }
        for _, join := range s.joins {
            jbw, jac := join.scan(b[bw:], args)
            bw += jbw
            ac += jac
        }
    }
    if s.where != nil {
        wbw, wac := s.where.scan(b[bw:], args[ac:])
        bw += wbw
        ac += wac
    }
    if s.groupBy != nil {
        gbbw, gbac := s.groupBy.scan(b[bw:], args[ac:])
        bw += gbbw
        ac += gbac
    }
    if s.orderBy != nil {
        obbw, obac := s.orderBy.scan(b[bw:], args[ac:])
        bw += obbw
        ac += obac
    }
    if s.limit != nil {
        lbw, lac := s.limit.scan(b[bw:], args[ac:])
        bw += lbw
        ac += lac
    }
    return bw, ac
}

func (s *selectClause) addJoin(jc *joinClause) *selectClause {
    s.joins = append(s.joins, jc)
    return s
}

func (s *selectClause) addWhere(e *Expression) *selectClause {
    if s.where == nil {
        s.where = &whereClause{filters: make([]*Expression, 0)}
    }
    s.where.filters = append(s.where.filters, e)
    return s
}

// Given one or more columns, either set or add to the GROUP BY clause for
// the selectClause
func (s *selectClause) addGroupBy(cols ...projection) *selectClause {
    if len(cols) == 0 {
        return s
    }
    gb := s.groupBy
    if gb == nil {
        gb = &groupByClause{
            cols: make([]projection, len(cols)),
        }
        for x, c := range cols {
            gb.cols[x] = c
        }
    } else {
        for _, c := range cols {
            gb.cols = append(gb.cols, c)
        }
    }
    s.groupBy = gb
    return s
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the selectClause
func (s *selectClause) addOrderBy(sortCols ...*sortColumn) *selectClause {
    if len(sortCols) == 0 {
        return s
    }
    ob := s.orderBy
    if ob == nil {
        ob = &orderByClause{
            scols: make([]*sortColumn, len(sortCols)),
        }
        for x, sc := range sortCols {
            ob.scols[x] = sc
        }
    } else {
        for _, sc := range sortCols {
            ob.scols = append(ob.scols, sc)
        }
    }
    s.orderBy = ob
    return s
}

func (s *selectClause) setLimitWithOffset(limit int, offset int) *selectClause {
    lc := &limitClause{limit: limit}
    lc.offset = &offset
    s.limit = lc
    return s
}

func (s *selectClause) setLimit(limit int) *selectClause {
    lc := &limitClause{limit: limit}
    s.limit = lc
    return s
}

func containsJoin(s *selectClause, j *joinClause) bool {
    for _, sj := range s.joins {
        if j == sj {
            return true
        }
    }
    return false
}

func addToProjections(s *selectClause, p projection) {
    s.projs = append(s.projs, p)
}

func (s *selectClause) removeSelection(toRemove selection) {
    idx := -1
    for x, sel := range s.selections {
        if sel == toRemove {
            idx = x
            break
        }
    }
    if idx == -1 {
        return
    }
    s.selections = append(s.selections[:idx], s.selections[idx + 1:]...)
}
