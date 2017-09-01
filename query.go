package sqlb

import (
    "errors"
    "fmt"
)

var (
    ERR_JOIN_INVALID = errors.New("Unable to join selection. Either there was no selection to join to or the target selection was not found.")
)

type Query struct {
    e error
    b []byte
    args []interface{}
    sel *selectClause
}

func (q *Query) IsValid() bool {
    return q.e == nil &&  q.sel != nil
}

func (q *Query) Error() error {
    return q.e
}

func (q *Query) String() string {
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

func (q *Query) StringArgs() (string, []interface{}) {
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

func (q *Query) Where(e *Expression) *Query {
    q.sel.addWhere(e)
    return q
}

func (q *Query) GroupBy(cols ...projection) *Query {
    q.sel.addGroupBy(cols...)
    return q
}

func (q *Query) OrderBy(scols ...*sortColumn) *Query {
    q.sel.addOrderBy(scols...)
    return q
}

func (q *Query) Limit(limit int) *Query {
    q.sel.setLimit(limit)
    return q
}

func (q *Query) LimitWithOffset(limit int, offset int) *Query {
    q.sel.setLimitWithOffset(limit, offset)
    return q
}

func (q *Query) As(alias string) *Query {
    q.sel.setAlias(alias)
    return q
}

// Join to a supplied selection with the supplied ON expression. If the Query
// does not yet contain a selectClause OR if the supplied ON expression does
// not reference any selection that is found in the Query's selectClause, then
// Query.e will be set to an error.
func (q *Query) Join(right selection, onExpr *Expression) *Query {
    if q.sel == nil {
        q.e = ERR_JOIN_INVALID
        fmt.Println("No select clause.")
        return q
    }

    // Let's first determine which selection is targeted as the LEFT part of
    // the join.
    var left selection
    rightSelId := right.selectionId()
    for _, el := range onExpr.elements {
        switch el.(type) {
            case projection:
                p := el.(projection)
                exprSelId := p.from().selectionId()
                if exprSelId == rightSelId {
                    continue
                }
                // Search through the Query's primary selectClause, looking for
                // the selection that is referred to be the ON expression.
                for _, sel := range q.sel.selections {
                    if sel.selectionId() == exprSelId {
                        left = sel
                        break
                    }
                }
                if left != nil {
                    break
                }
        }
    }
    if left == nil {
        q.e = ERR_JOIN_INVALID
        return q
    }
    jc := Join(left, right, onExpr)
    q.sel.addJoin(jc)

    // Make sure we remove the right-hand selection from the selectClause's
    // selections collection, since it's in a JOIN clause.
    q.sel.removeSelection(right)
    return q
}

func Select(items ...interface{}) *Query {
    sel := &selectClause{
        projections: make([]projection, 0),
    }

    selectionMap := make(map[uint64]selection, 0)
    projectionMap := make(map[uint64]projection, 0)

    // For each scannable item we've received in the call, check what concrete
    // type they are and, depending on which type they are, either add them to
    // the returned selectClause's projections list or query the underlying
    // table metadata to generate a list of all columns in that table.
    for _, item := range items {
        switch item.(type) {
            case *joinClause:
                j := item.(*joinClause)
                if ! containsJoin(sel, j) {
                    sel.joins = append(sel.joins, j)
                    if _, ok := selectionMap[j.left.selectionId()]; ! ok {
                        selectionMap[j.left.selectionId()] = j.left
                        for _, proj := range j.left.projections() {
                            projId := proj.projectionId()
                            _, projExists := projectionMap[projId]
                            if ! projExists {
                                addToProjections(sel, proj)
                                projectionMap[projId] = proj
                            }
                        }
                    }
                    if _, ok := selectionMap[j.right.selectionId()]; ! ok {
                        for _, proj := range j.right.projections() {
                            projId := proj.projectionId()
                            _, projExists := projectionMap[projId]
                            if ! projExists {
                                addToProjections(sel, proj)
                                projectionMap[projId] = proj
                            }
                        }
                    }
                }
            case *Column:
                v := item.(*Column)
                sel.projections = append(sel.projections, v)
                selectionMap[v.tbl.selectionId()] = v.tbl
            case *Table:
                v := item.(*Table)
                for _, cd := range v.tdef.projections() {
                    addToProjections(sel, cd)
                }
                selectionMap[v.selectionId()] = v
            case *TableDef:
                v := item.(*TableDef)
                for _, cd := range v.projections() {
                    addToProjections(sel, cd)
                }
                selectionMap[v.selectionId()] = v
            case *ColumnDef:
                v := item.(*ColumnDef)
                addToProjections(sel, v)
                selectionMap[v.tdef.selectionId()] = v.tdef
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
    for _, sel := range selectionMap {
        selections[x] = sel
        x++
    }
    sel.selections = selections
    return &Query{sel: sel}
}
