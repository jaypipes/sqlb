package sqlb

import (
    "errors"
    "fmt"
)

var (
    ERR_JOIN_INVALID_NO_SELECT = errors.New("Unable to join selection. There was no selection to join to.")
    ERR_JOIN_INVALID_UNKNOWN_TARGET = errors.New("Unable to join selection. Target selection was not found.")
)

type SelectQuery struct {
    e error
    b []byte
    args []interface{}
    sel *selectClause
}

func (q *SelectQuery) IsValid() bool {
    return q.e == nil &&  q.sel != nil
}

func (q *SelectQuery) Error() error {
    return q.e
}

func (q *SelectQuery) String() string {
    size := q.sel.size()
    argc := q.sel.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.sel.scan(q.b, q.args)
    return string(q.b)
}

func (q *SelectQuery) StringArgs() (string, []interface{}) {
    size := q.sel.size()
    argc := q.sel.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.sel.scan(q.b, q.args)
    return string(q.b), q.args
}

func (q *SelectQuery) Where(e *Expression) *SelectQuery {
    q.sel.addWhere(e)
    return q
}

func (q *SelectQuery) GroupBy(cols ...projection) *SelectQuery {
    q.sel.addGroupBy(cols...)
    return q
}

func (q *SelectQuery) OrderBy(scols ...*sortColumn) *SelectQuery {
    q.sel.addOrderBy(scols...)
    return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
    q.sel.setLimit(limit)
    return q
}

func (q *SelectQuery) LimitWithOffset(limit int, offset int) *SelectQuery {
    q.sel.setLimitWithOffset(limit, offset)
    return q
}

// Returns a pointer to a new SelectQuery that has aliased its inner selection
// to the supplied name
func (q *SelectQuery) As(alias string) *SelectQuery {
    dt := &derivedTable{
        alias: alias,
        from: q.sel,
    }
    derivedSel := &selectClause{
        projs: dt.getAllDerivedColumns(),
        selections: []selection{dt},
    }
    return &SelectQuery{sel: derivedSel}
}

// Returns the projection of the underlying selectClause that matches the name
// provided
func (q *SelectQuery) Column(name string) projection {
    for _, p := range q.sel.projs {
        switch p.(type) {
        case *derivedColumn:
            dc := p.(*derivedColumn)
            if dc.alias != "" && dc.alias == name {
                return dc
            } else if dc.c.cdef.name == name {
                return dc
            }
        case *Column:
            c := p.(*Column)
            if c.alias != "" && c.alias == name {
                return c
            } else if c.cdef.name == name {
                return c
            }
        case *ColumnDef:
            cd := p.(*ColumnDef)
            if cd.name == name {
                return cd
            }
        case *sqlFunc:
            f := p.(*sqlFunc)
            if f.alias != "" && f.alias == name {
                return f
            }
        }
    }
    return nil
}

func (q *SelectQuery) Join(right interface{}, on *Expression) *SelectQuery {
    var rightSel selection
    switch right.(type) {
    case selection:
        rightSel = right.(selection)
    case *SelectQuery:
        // Joining to a derived table
        rightSel = right.(*SelectQuery).sel.selections[0]
    }
    return q.doJoin(JOIN_INNER, rightSel, on)
}

func (q *SelectQuery) OuterJoin(right selection, on *Expression) *SelectQuery {
    return q.doJoin(JOIN_OUTER, right, on)
}

// Join to a supplied selection with the supplied ON expression. If the SelectQuery
// does not yet contain a selectClause OR if the supplied ON expression does
// not reference any selection that is found in the SelectQuery's selectClause, then
// SelectQuery.e will be set to an error.
func (q *SelectQuery) doJoin(
    jt joinType,
    right selection,
    on *Expression,
) *SelectQuery {
    if q.sel == nil || len(q.sel.selections) == 0 {
        q.e = ERR_JOIN_INVALID_NO_SELECT
        return q
    }

    // Let's first determine which selection is targeted as the LEFT part of
    // the join.
    var left selection
    if on != nil {
        for _, el := range on.elements {
            switch el.(type) {
            case projection:
                p := el.(projection)
                exprSel := p.from()
                if exprSel == right {
                    continue
                }
                // Search through the SelectQuery's primary selectClause, looking for
                // the selection that is referred to be the ON expression.
                for _, sel := range q.sel.selections {
                    if sel == exprSel {
                        left = sel
                        break
                    }
                }
                if left != nil {
                    break
                }
            case *Expression:
                expr := el.(*Expression)
                for _, referrent := range expr.referrents() {
                    if referrent == right {
                        continue
                    }
                    for _, sel := range q.sel.selections {
                        if sel == referrent {
                            left = sel
                            break
                        }
                    }
                    if left != nil {
                        break
                    }
                }
                if left != nil {
                    break
                }
            }
        }
    } else {
        // TODO(jaypipes): Handle CROSS JOIN by joining the supplied right
        // against a derivedTable constructed from the existing SelectQuery.sel
        // selectClause
    }
    if left == nil {
        q.e = ERR_JOIN_INVALID_UNKNOWN_TARGET
        return q
    }
    jc := &joinClause{
        joinType: jt,
        left: left,
        right: right,
        on: on,
    }
    q.sel.addJoin(jc)

    // Make sure we remove the right-hand selection from the selectClause's
    // selections collection, since it's in a JOIN clause.
    q.sel.removeSelection(right)
    return q
}

func Select(items ...interface{}) *SelectQuery {
    sel := &selectClause{
        projs: make([]projection, 0),
    }

    nDerived := 0
    selectionMap := make(map[selection]bool, 0)
    projectionMap := make(map[uint64]projection, 0)

    // For each scannable item we've received in the call, check what concrete
    // type they are and, depending on which type they are, either add them to
    // the returned selectClause's projections list or query the underlying
    // table metadata to generate a list of all columns in that table.
    for _, item := range items {
        switch item.(type) {
            case *SelectQuery:
                // Project all columns from the subquery to the outer
                // selectClause
                sq := item.(*SelectQuery)
                innerSelClause := sq.sel
                if len(innerSelClause.selections) == 1 {
                    innerSel := innerSelClause.selections[0]
                    switch innerSel.(type) {
                        case *derivedTable:
                            // If the inner select clause contains a single
                            // selection and that selection is a derivedTable,
                            // that means we were called like so:
                            //
                            //      Select(Select(...).As("alias"))
                            //
                            // This means that we do *not* need to generate a
                            // derived table but instead simply grab the
                            // existing derived table as the single selection
                            // for the outer selectClause and project all the
                            // derived table's projections out into the outer
                            // selectClause.
                            selectionMap[innerSel] = true
                            dt := innerSel.(*derivedTable)
                            for _, p := range dt.getAllDerivedColumns() {
                                pid := p.projectionId()
                                projectionMap[pid] = p
                                addToProjections(sel, p)
                            }
                        default:
                            // This means we were called like so:
                            //
                            //     Select(Select(...))
                            //
                            // So we need to construct a derived table manually
                            // and name it derivedN.
                            derivedName := fmt.Sprintf("derived%d", nDerived)
                            dt := &derivedTable{
                                alias: derivedName,
                                from: innerSelClause,
                            }
                            selectionMap[dt] = true
                            for _, p := range dt.getAllDerivedColumns() {
                                pid := p.projectionId()
                                projectionMap[pid] = p
                                addToProjections(sel, p)
                            }
                            nDerived++
                    }
                }
            case *Column:
                v := item.(*Column)
                sel.projs = append(sel.projs, v)
                selectionMap[v.tbl] = true
            case *Table:
                v := item.(*Table)
                for _, cd := range v.tdef.projections() {
                    addToProjections(sel, cd)
                }
                selectionMap[v] = true
            case *TableDef:
                v := item.(*TableDef)
                for _, cd := range v.projections() {
                    addToProjections(sel, cd)
                }
                selectionMap[v] = true
            case *ColumnDef:
                v := item.(*ColumnDef)
                addToProjections(sel, v)
                selectionMap[v.tdef] = true
            default:
                // Everything else, make it a literal value projection, so, for
                // instance, a user can do SELECT 1, which is, technically
                // valid SQL.
                p := &value{val: item}
                addToProjections(sel, p)
        }
    }
    selections := make([]selection, len(selectionMap))
    x := 0
    for sel, _ := range selectionMap {
        selections[x] = sel
        x++
    }
    sel.selections = selections
    return &SelectQuery{sel: sel}
}
