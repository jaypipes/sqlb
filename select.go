package sqlb

import (
    "errors"
    "fmt"
)

var (
    ERR_JOIN_INVALID = errors.New("Unable to join selection. Either there was no selection to join to or the target selection was not found.")
)

type SelectQuery struct {
    e error
    b []byte
    args []interface{}
    sel selection
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
    q.sel.(*selectClause).addWhere(e)
    return q
}

func (q *SelectQuery) GroupBy(cols ...projection) *SelectQuery {
    q.sel.(*selectClause).addGroupBy(cols...)
    return q
}

func (q *SelectQuery) OrderBy(scols ...*sortColumn) *SelectQuery {
    q.sel.(*selectClause).addOrderBy(scols...)
    return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
    q.sel.(*selectClause).setLimit(limit)
    return q
}

func (q *SelectQuery) LimitWithOffset(limit int, offset int) *SelectQuery {
    q.sel.(*selectClause).setLimitWithOffset(limit, offset)
    return q
}

// Returns a pointer to a new SelectQuery that has aliased its inner selection
// to the supplied name
func (q *SelectQuery) As(alias string) *SelectQuery {
    q.sel.(*selectClause).setAlias(alias)
    aliasProjs := make([]projection, len(q.sel.(*selectClause).projs))
    for x, origProj := range q.sel.(*selectClause).projs {
        aliasProjs[x] = origProj
    }
    derived := &selectClause{
        projs: aliasProjs,
        selections: []selection{q.sel},
    }
    return &SelectQuery{sel: derived}
}

// Join to a supplied selection with the supplied ON expression. If the SelectQuery
// does not yet contain a selectClause OR if the supplied ON expression does
// not reference any selection that is found in the SelectQuery's selectClause, then
// SelectQuery.e will be set to an error.
func (q *SelectQuery) Join(right selection, onExpr *Expression) *SelectQuery {
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
                // Search through the SelectQuery's primary selectClause, looking for
                // the selection that is referred to be the ON expression.
                for _, sel := range q.sel.(*selectClause).selections {
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
    q.sel.(*selectClause).addJoin(jc)

    // Make sure we remove the right-hand selection from the selectClause's
    // selections collection, since it's in a JOIN clause.
    q.sel.(*selectClause).removeSelection(right)
    return q
}

func Select(items ...interface{}) *SelectQuery {
    sel := &selectClause{
        projs: make([]projection, 0),
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
                sel.projs = append(sel.projs, v)
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
    return &SelectQuery{sel: sel}
}
