//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"errors"
	"fmt"

	"github.com/jaypipes/sqlb/pkg/scanner"
	"github.com/jaypipes/sqlb/pkg/types"
)

var (
	ERR_JOIN_INVALID_NO_SELECT      = errors.New("Unable to join selection. There was no selection to join to.")
	ERR_JOIN_INVALID_UNKNOWN_TARGET = errors.New("Unable to join selection. Target selection was not found.")
)

type SelectQuery struct {
	e       error
	b       []byte
	args    []interface{}
	sel     *SelectStatement
	scanner types.Scanner
}

func (q *SelectQuery) IsValid() bool {
	return q.e == nil && q.sel != nil
}

func (q *SelectQuery) Error() error {
	return q.e
}

func (q *SelectQuery) String() string {
	sizes := q.scanner.Size(q.sel)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.Scan(q.b, q.args, q.sel)
	return string(q.b)
}

func (q *SelectQuery) StringArgs() (string, []interface{}) {
	sizes := q.scanner.Size(q.sel)
	if len(q.args) != sizes.ArgCount {
		q.args = make([]interface{}, sizes.ArgCount)
	}
	if len(q.b) != sizes.BufferSize {
		q.b = make([]byte, sizes.BufferSize)
	}
	q.scanner.Scan(q.b, q.args, q.sel)
	return string(q.b), q.args
}

func (q *SelectQuery) Where(e *Expression) *SelectQuery {
	q.sel.addWhere(e)
	return q
}

func (q *SelectQuery) GroupBy(cols ...types.Projection) *SelectQuery {
	q.sel.addGroupBy(cols...)
	return q
}

func (q *SelectQuery) Having(e *Expression) *SelectQuery {
	q.sel.addHaving(e)
	return q
}

func (q *SelectQuery) OrderBy(scols ...*SortColumn) *SelectQuery {
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
	dt := &DerivedTable{
		alias: alias,
		from:  q.sel,
	}
	derivedSel := &SelectStatement{
		projs:      dt.getAllDerivedColumns(),
		selections: []types.Selection{dt},
	}
	return &SelectQuery{sel: derivedSel, scanner: q.scanner}
}

// Returns the projection of the underlying SelectStatement that matches the name
// provided
func (q *SelectQuery) C(name string) types.Projection {
	for _, p := range q.sel.projs {
		switch p.(type) {
		case *DerivedColumn:
			dc := p.(*DerivedColumn)
			if dc.alias != "" && dc.alias == name {
				return dc
			} else if dc.c.name == name {
				return dc
			}
		case *Column:
			c := p.(*Column)
			if c.alias != "" && c.alias == name {
				return c
			} else if c.name == name {
				return c
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
	var rightSel types.Selection
	switch right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.(*SelectQuery).sel.selections[0]
	case types.Selection:
		rightSel = right.(types.Selection)
	}
	return q.doJoin(JOIN_INNER, rightSel, on)
}

func (q *SelectQuery) OuterJoin(right interface{}, on *Expression) *SelectQuery {
	var rightSel types.Selection
	switch right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.(*SelectQuery).sel.selections[0]
	case types.Selection:
		rightSel = right.(types.Selection)
	}
	return q.doJoin(JOIN_OUTER, rightSel, on)
}

// Join to a supplied selection with the supplied ON expression. If the SelectQuery
// does not yet contain a SelectStatement OR if the supplied ON expression does
// not reference any selection that is found in the SelectQuery's SelectStatement, then
// SelectQuery.e will be set to an error.
func (q *SelectQuery) doJoin(
	jt JoinType,
	right types.Selection,
	on *Expression,
) *SelectQuery {
	if q.sel == nil || len(q.sel.selections) == 0 {
		q.e = ERR_JOIN_INVALID_NO_SELECT
		return q
	}

	// Let's first determine which selection is targeted as the LEFT part of
	// the join.
	var left types.Selection
	if on != nil {
		for _, el := range on.elements {
			switch el.(type) {
			case types.Projection:
				p := el.(types.Projection)
				exprSel := p.From()
				if exprSel == right {
					continue
				}
				// Search through the SelectQuery's primary SelectStatement, looking for
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
				// Now search through the SelectQuery's JoinClauses, looking
				// for a selection that is the left side of the ON expression
				for _, j := range q.sel.joins {
					if j.left == exprSel {
						left = j.left
					} else if j.right == exprSel {
						left = j.right
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
					// Now search through the SelectQuery's JoinClauses, looking
					// for a selection that is the left side of the ON expression
					for _, j := range q.sel.joins {
						if j.left == referrent {
							left = j.left
						} else if j.right == referrent {
							left = j.right
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
		// against a DerivedTable constructed from the existing SelectQuery.sel
		// SelectStatement
	}
	if left == nil {
		q.e = ERR_JOIN_INVALID_UNKNOWN_TARGET
		return q
	}
	jc := &JoinClause{
		JoinType: jt,
		left:     left,
		right:    right,
		on:       on,
	}
	q.sel.addJoin(jc)

	// Make sure we remove the right-hand selection from the SelectStatement's
	// selections collection, since it's in a JOIN clause.
	q.sel.removeSelection(right)
	return q
}

func Select(items ...interface{}) *SelectQuery {
	scanner := scanner.New(types.DIALECT_UNKNOWN)
	sq := &SelectQuery{
		scanner: scanner,
	}
	sel := &SelectStatement{
		projs: make([]types.Projection, 0),
	}

	nDerived := 0
	selectionMap := make(map[types.Selection]bool, 0)

	// For each scannable item we've received in the call, check what concrete
	// type they are and, depending on which type they are, either add them to
	// the returned SelectStatement's projections list or query the underlying
	// table metadata to generate a list of all columns in that table.
	for _, item := range items {
		switch item.(type) {
		case *SelectQuery:
			// Project all columns from the subquery to the outer
			// SelectStatement
			isq := item.(*SelectQuery)
			sq.scanner = isq.scanner
			innerSelClause := isq.sel
			if len(innerSelClause.selections) == 1 {
				innerSel := innerSelClause.selections[0]
				switch innerSel.(type) {
				case *DerivedTable:
					// If the inner select clause contains a single
					// selection and that selection is a DerivedTable,
					// that means we were called like so:
					//
					//      Select(Select(...).As("alias"))
					//
					// This means that we do *not* need to generate a
					// derived table but instead simply grab the
					// existing derived table as the single selection
					// for the outer SelectStatement and project all the
					// derived table's projections out into the outer
					// SelectStatement.
					selectionMap[innerSel] = true
					dt := innerSel.(*DerivedTable)
					for _, p := range dt.getAllDerivedColumns() {
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
					dt := &DerivedTable{
						alias: derivedName,
						from:  innerSelClause,
					}
					selectionMap[dt] = true
					for _, p := range dt.getAllDerivedColumns() {
						addToProjections(sel, p)
					}
					nDerived++
				}
			}
		case *Column:
			v := item.(*Column)
			// Set scanner's dialect based on supplied meta's dialect
			sq.scanner.WithDialect(v.tbl.meta.dialect)
			sel.projs = append(sel.projs, v)
			selectionMap[v.tbl] = true
		case *Table:
			v := item.(*Table)
			// Set scanner's dialect based on supplied meta's dialect
			sq.scanner.WithDialect(v.meta.dialect)
			for _, c := range v.Projections() {
				addToProjections(sel, c)
			}
			selectionMap[v] = true
		case *sqlFunc:
			v := item.(*sqlFunc)
			addToProjections(sel, v)
			selectionMap[v.sel] = true
		default:
			// Everything else, make it a literal value projection, so, for
			// instance, a user can do SELECT 1, which is, technically
			// valid SQL.
			p := &value{val: item}
			addToProjections(sel, p)
		}
	}
	selections := make([]types.Selection, len(selectionMap))
	x := 0
	for sel, _ := range selectionMap {
		selections[x] = sel
		x++
	}
	sel.selections = selections
	sq.sel = sel
	return sq
}
