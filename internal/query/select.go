// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package query

import (
	"fmt"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar/clause"
	"github.com/jaypipes/sqlb/internal/grammar/element"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
	"github.com/jaypipes/sqlb/internal/grammar/function"
	"github.com/jaypipes/sqlb/internal/grammar/identifier"
	"github.com/jaypipes/sqlb/internal/grammar/statement"
)

type SelectQuery struct {
	sel *statement.Select
}

func (q *SelectQuery) Element() api.Element {
	return q.sel
}

func (q *SelectQuery) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	return q.sel.String(opts, qargs, curarg)
}

func (q *SelectQuery) ArgCount() int {
	return q.sel.ArgCount()
}

func (q *SelectQuery) Where(e *expression.Expression) *SelectQuery {
	q.sel.AddWhere(e)
	return q
}

func (q *SelectQuery) GroupBy(cols ...api.Projection) *SelectQuery {
	q.sel.AddGroupBy(cols...)
	return q
}

func (q *SelectQuery) Having(e *expression.Expression) *SelectQuery {
	q.sel.AddHaving(e)
	return q
}

func (q *SelectQuery) OrderBy(scols ...api.Orderable) *SelectQuery {
	q.sel.AddOrderBy(scols...)
	return q
}

func (q *SelectQuery) Limit(limit int) *SelectQuery {
	q.sel.SetLimit(limit)
	return q
}

func (q *SelectQuery) LimitWithOffset(limit int, offset int) *SelectQuery {
	q.sel.SetLimitWithOffset(limit, offset)
	return q
}

// Returns a pointer to a new SelectQuery that has aliased its inner selection
// to the supplied name
func (q *SelectQuery) As(alias string) *SelectQuery {
	dt := clause.NewDerivedTable(alias, q.sel)
	derivedSel := statement.NewSelect(
		dt.DerivedColumns(),
		[]api.Selection{dt},
		nil, nil, nil, nil, nil, nil,
	)
	return &SelectQuery{sel: derivedSel}
}

// Returns the projection of the underlying SelectStatement that matches the name
// provided
func (q *SelectQuery) C(name string) api.Projection {
	for _, p := range q.sel.Projections() {
		switch p := p.(type) {
		case *clause.DerivedColumn:
			dc := p
			if dc.Alias != "" && dc.Alias == name {
				return dc
			} else if dc.C().Name == name {
				return dc
			}
		case *identifier.Column:
			c := p
			if c.Alias != "" && c.Alias == name {
				return c
			} else if c.Name == name {
				return c
			}
		case *function.Function:
			f := p
			if f.Alias != "" && f.Alias == name {
				return f
			}
		}
	}
	return nil
}

func (q *SelectQuery) Join(
	right interface{},
	on *expression.Expression,
) *SelectQuery {
	var rightSel api.Selection
	switch right := right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.sel.Selections()[0]
	case api.Selection:
		rightSel = right
	}
	return q.doJoin(api.JoinInner, rightSel, on)
}

func (q *SelectQuery) InnerJoin(
	right interface{},
	on *expression.Expression,
) *SelectQuery {
	return q.Join(right, on)
}

func (q *SelectQuery) OuterJoin(
	right interface{},
	on *expression.Expression,
) *SelectQuery {
	var rightSel api.Selection
	switch right := right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.sel.Selections()[0]
	case api.Selection:
		rightSel = right
	}
	return q.doJoin(api.JoinOuter, rightSel, on)
}

// Join to a supplied selection with the supplied ON expression. If the SelectQuery
// does not yet contain a SelectStatement OR if the supplied ON expression does
// not reference any selection that is found in the SelectQuery's SelectStatement, then
// SelectQuery.e will be set to an error.
func (q *SelectQuery) doJoin(
	jt api.JoinType,
	right api.Selection,
	on *expression.Expression,
) *SelectQuery {
	if q.sel == nil || len(q.sel.Selections()) == 0 {
		panic(api.InvalidJoinNoSelect)
	}

	// Let's first determine which selection is targeted as the LEFT part of
	// the join.
	var left api.Selection
	if on != nil {
		for _, el := range on.Elements() {
			switch el := el.(type) {
			case api.Projection:
				p := el
				exprSel := p.From()
				if exprSel == right {
					continue
				}
				// Search through the SelectQuery's primary SelectStatement, looking for
				// the selection that is referred to in the ON expression.
				for _, sel := range q.sel.Selections() {
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
				for _, j := range q.sel.Joins() {
					if j.Left() == exprSel {
						left = j.Left()
					} else if j.Right() == exprSel {
						left = j.Right()
					}
				}
				if left != nil {
					break
				}
			case *expression.Expression:
				expr := el
				for _, referrent := range expr.Referrents() {
					if referrent == right {
						continue
					}
					for _, sel := range q.sel.Selections() {
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
					for _, j := range q.sel.Joins() {
						if j.Left() == referrent {
							left = j.Left()
						} else if j.Right() == referrent {
							left = j.Right()
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
	}

	// TODO(jaypipes): Handle CROSS JOIN by joining the supplied right
	// against a DerivedTable constructed from the existing SelectQuery.sel
	// SelectStatement

	if left == nil {
		panic(api.InvalidJoinUnknownTarget)
	}
	jc := clause.NewJoin(jt, left, right, on)
	q.sel.AddJoin(jc)

	// Make sure we remove the right-hand selection from the SelectStatement's
	// selections collection, since it's in a JOIN clause.
	q.sel.RemoveSelection(right)
	return q
}

func Select(
	items ...interface{},
) *SelectQuery {
	sel := statement.NewSelect(make([]api.Projection, 0), nil, nil, nil, nil, nil, nil, nil)

	nDerived := 0
	selectionMap := make(map[api.Selection]bool, 0)

	// For each scannable item we've received in the call, check what concrete
	// type they are and, depending on which type they are, either add them to
	// the returned SelectStatement's projections list or query the underlying
	// table metadata to generate a list of all columns in that table.
	for _, item := range items {
		switch item := item.(type) {
		case *SelectQuery:
			// Project all columns from the subquery to the outer
			// SelectStatement
			isq := item
			innerSelClause := isq.sel
			if len(innerSelClause.Selections()) == 1 {
				innerSel := innerSelClause.Selections()[0]
				switch innerSel := innerSel.(type) {
				case *clause.DerivedTable:
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
					dt := innerSel
					for _, p := range dt.DerivedColumns() {
						sel.AddProjection(p)
					}
				default:
					// This means we were called like so:
					//
					//     Select(Select(...))
					//
					// So we need to construct a derived table manually
					// and name it derivedN.
					derivedName := fmt.Sprintf("derived%d", nDerived)
					dt := clause.NewDerivedTable(derivedName, innerSelClause)
					selectionMap[dt] = true
					for _, p := range dt.DerivedColumns() {
						sel.AddProjection(p)
					}
					nDerived++
				}
			}
		case *identifier.Column:
			v := item
			if v == nil {
				panic("specified a non-existent column")
			}
			sel.AddProjection(v)
			selectionMap[v.From()] = true
		case *identifier.Table:
			v := item
			for _, c := range v.Projections() {
				sel.AddProjection(c)
			}
			selectionMap[v] = true
		case *function.Function:
			v := item
			sel.AddProjection(v)
			selectionMap[v.From()] = true
		default:
			// Everything else, make it a literal value projection, so, for
			// instance, a user can do SELECT 1, which is, technically
			// valid SQL.
			p := element.NewValue(nil, item)
			sel.AddProjection(p)
		}
	}
	selections := make([]api.Selection, len(selectionMap))
	x := 0
	for sel := range selectionMap {
		selections[x] = sel
		x++
	}
	sel.ReplaceSelections(selections)
	return &SelectQuery{sel: sel}
}
