package sqlb

type selectClause struct {
    alias string
    projected *List
    selections []selection
    joins []*joinClause
    where *whereClause
    groupBy *groupByClause
    orderBy *orderByClause
    limit *limitClause
}

func (s *selectClause) isSubselect() bool {
    return s.alias != ""
}

func (s *selectClause) argCount() int {
    argc := s.projected.argCount()
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

func (s *selectClause) setAlias(alias string) {
    s.alias = alias
}

func (s *selectClause) As(alias string) *selectClause {
    s.setAlias(alias)
    return s
}

func (s *selectClause) size() int {
    size := len(Symbols[SYM_SELECT]) + len(Symbols[SYM_FROM])
    size += s.projected.size()
    for _, sel := range s.selections {
        size += sel.size()
    }
    for _, join := range s.joins {
        size += join.size()
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
    if s.isSubselect() {
        size += (len(Symbols[SYM_LPAREN]) + len(Symbols[SYM_RPAREN]) +
                 len(Symbols[SYM_AS]) + len(s.alias))
    }
    return size
}

func (s *selectClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    if s.isSubselect() {
        bw += copy(b[bw:], Symbols[SYM_LPAREN])
    }
    bw += copy(b[bw:], Symbols[SYM_SELECT])
    pbw, pac := s.projected.scan(b[bw:], args)
    bw += pbw
    ac += pac
    bw += copy(b[bw:], Symbols[SYM_FROM])
    for _, sel := range s.selections {
        sbw, sac := sel.scan(b[bw:], args)
        bw += sbw
        ac += sac
    }
    for _, join := range s.joins {
        jbw, jac := join.scan(b[bw:], args)
        bw += jbw
        ac += jac
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
    if s.isSubselect() {
        bw += copy(b[bw:], Symbols[SYM_RPAREN])
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], s.alias)
    }
    return bw, ac
}

func (s *selectClause) Where(e *Expression) *selectClause {
    if s.where == nil {
        s.where = &whereClause{filters: make([]*Expression, 0)}
    }
    s.where.filters = append(s.where.filters, e)
    return s
}

// Given one or more columns, either set or add to the GROUP BY clause for
// the selectClause
func (s *selectClause) GroupBy(cols ...Columnar) *selectClause {
    if len(cols) == 0 {
        return s
    }
    gb := s.groupBy
    if gb == nil {
        gb = &groupByClause{
            cols: &List{
                elements: make([]element, len(cols)),
            },
        }
        for x, c := range cols {
            gb.cols.elements[x] = c.Column()
        }
    } else {
        for _, c := range cols {
            gb.cols.elements = append(gb.cols.elements, c.Column())
        }
    }
    s.groupBy = gb
    return s
}

// Given one or more sort columns, either set or add to the ORDER BY clause for
// the selectClause
func (s *selectClause) OrderBy(sortCols ...*sortColumn) *selectClause {
    if len(sortCols) == 0 {
        return s
    }
    ob := s.orderBy
    if ob == nil {
        ob = &orderByClause{
            cols: &List{
                elements: make([]element, len(sortCols)),
            },
        }
        for x, sc := range sortCols {
            ob.cols.elements[x] = sc
        }
    } else {
        for _, sc := range sortCols {
            ob.cols.elements = append(ob.cols.elements, sc)
        }
    }
    s.orderBy = ob
    return s
}

func (s *selectClause) LimitWithOffset(limit int, offset int) *selectClause {
    lc := &limitClause{limit: limit}
    lc.offset = &offset
    s.limit = lc
    return s
}

func (s *selectClause) Limit(limit int) *selectClause {
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
    s.projected.elements = append(s.projected.elements, p)
}

func Select(items ...element) *selectClause {
    // TODO(jaypipes): Make the memory allocation more efficient below by
    // looping through the elements and determining the number of element struct
    // pointers to allocate instead of just making an empty array of element
    // pointers.
    res := &selectClause{
        projected: &List{},
    }

    selectionMap := make(map[uint64]selection, 0)
    projectionMap := make(map[uint64]projection, 0)

    // For each scannable item we've received in the call, check what concrete
    // type they are and, depending on which type they are, either add them to
    // the returned selectClause's projected List or query the underlying
    // table metadata to generate a list of all columns in that table.
    for _, item := range items {
        switch item.(type) {
            case *joinClause:
                j := item.(*joinClause)
                if ! containsJoin(res, j) {
                    res.joins = append(res.joins, j)
                    if _, ok := selectionMap[j.left.selectionId()]; ! ok {
                        selectionMap[j.left.selectionId()] = j.left
                        for _, proj := range j.left.projections() {
                            projId := proj.projectionId()
                            _, projExists := projectionMap[projId]
                            if ! projExists {
                                addToProjections(res, proj)
                                projectionMap[projId] = proj
                            }
                        }
                    }
                    if _, ok := selectionMap[j.right.selectionId()]; ! ok {
                        for _, proj := range j.right.projections() {
                            projId := proj.projectionId()
                            _, projExists := projectionMap[projId]
                            if ! projExists {
                                addToProjections(res, proj)
                                projectionMap[projId] = proj
                            }
                        }
                    }
                }
            case *Column:
                v := item.(*Column)
                res.projected.elements = append(res.projected.elements, v)
                selectionMap[v.tbl.selectionId()] = v.tbl
            case *List:
                v := item.(*List)
                for _, el := range v.elements {
                    res.projected.elements = append(res.projected.elements, el)
                    if isColumn(el) {
                        c := el.(*Column)
                        selectionMap[c.tbl.selectionId()] = c.tbl
                    }
                }
            case *Table:
                v := item.(*Table)
                for _, cd := range v.tdef.projections() {
                    addToProjections(res, cd)
                }
                selectionMap[v.selectionId()] = v
            case *TableDef:
                v := item.(*TableDef)
                for _, cd := range v.projections() {
                    addToProjections(res, cd)
                }
                selectionMap[v.selectionId()] = v
            case *ColumnDef:
                v := item.(*ColumnDef)
                addToProjections(res, v)
                selectionMap[v.tdef.selectionId()] = v.tdef
        }
    }
    selections := make([]selection, len(selectionMap))
    x := 0
    for _, sel := range selectionMap {
        selections[x] = sel
        x++
    }
    res.selections = selections
    return res
}
