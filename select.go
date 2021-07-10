//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package sqlb

import (
	"errors"
	"fmt"

	"github.com/jaypipes/sqlb/pkg/ast"
	"github.com/jaypipes/sqlb/pkg/grammar/statement"
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
	sel     *statement.Select
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

func (q *SelectQuery) Where(e *ast.Expression) *SelectQuery {
	q.sel.AddWhere(e)
	return q
}

func (q *SelectQuery) GroupBy(cols ...types.Projection) *SelectQuery {
	q.sel.AddGroupBy(cols...)
	return q
}

func (q *SelectQuery) Having(e *ast.Expression) *SelectQuery {
	q.sel.AddHaving(e)
	return q
}

func (q *SelectQuery) OrderBy(scols ...*ast.SortColumn) *SelectQuery {
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
	dt := ast.NewDerivedTable(alias, q.sel)
	derivedSel := statement.NewSelect(
		dt.DerivedColumns(),
		[]types.Selection{dt},
		nil, nil, nil, nil, nil, nil,
	)
	return &SelectQuery{sel: derivedSel, scanner: q.scanner}
}

// Returns the projection of the underlying SelectStatement that matches the name
// provided
func (q *SelectQuery) C(name string) types.Projection {
	for _, p := range q.sel.Projections() {
		switch p.(type) {
		case *ast.DerivedColumn:
			dc := p.(*ast.DerivedColumn)
			if dc.Alias != "" && dc.Alias == name {
				return dc
			} else if dc.C().Name == name {
				return dc
			}
		case *ast.ColumnIdentifier:
			c := p.(*ast.ColumnIdentifier)
			if c.Alias != "" && c.Alias == name {
				return c
			} else if c.Name == name {
				return c
			}
		case *ast.Function:
			f := p.(*ast.Function)
			if f.Alias != "" && f.Alias == name {
				return f
			}
		}
	}
	return nil
}

func (q *SelectQuery) Join(right interface{}, on *ast.Expression) *SelectQuery {
	var rightSel types.Selection
	switch right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.(*SelectQuery).sel.Selections()[0]
	case types.Selection:
		rightSel = right.(types.Selection)
	}
	return q.doJoin(types.JOIN_INNER, rightSel, on)
}

func (q *SelectQuery) OuterJoin(right interface{}, on *ast.Expression) *SelectQuery {
	var rightSel types.Selection
	switch right.(type) {
	case *SelectQuery:
		// Joining to a derived table
		rightSel = right.(*SelectQuery).sel.Selections()[0]
	case types.Selection:
		rightSel = right.(types.Selection)
	}
	return q.doJoin(types.JOIN_OUTER, rightSel, on)
}

// Join to a supplied selection with the supplied ON expression. If the SelectQuery
// does not yet contain a SelectStatement OR if the supplied ON expression does
// not reference any selection that is found in the SelectQuery's SelectStatement, then
// SelectQuery.e will be set to an error.
func (q *SelectQuery) doJoin(
	jt types.JoinType,
	right types.Selection,
	on *ast.Expression,
) *SelectQuery {
	if q.sel == nil || len(q.sel.Selections()) == 0 {
		q.e = ERR_JOIN_INVALID_NO_SELECT
		return q
	}

	// Let's first determine which selection is targeted as the LEFT part of
	// the join.
	var left types.Selection
	if on != nil {
		for _, el := range on.Elements() {
			switch el.(type) {
			case types.Projection:
				p := el.(types.Projection)
				exprSel := p.From()
				if exprSel == right {
					continue
				}
				// Search through the SelectQuery's primary SelectStatement, looking for
				// the selection that is referred to be the ON expression.
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
			case *ast.Expression:
				expr := el.(*ast.Expression)
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
	} else {
		// TODO(jaypipes): Handle CROSS JOIN by joining the supplied right
		// against a DerivedTable constructed from the existing SelectQuery.sel
		// SelectStatement
	}
	if left == nil {
		q.e = ERR_JOIN_INVALID_UNKNOWN_TARGET
		return q
	}
	jc := ast.NewJoinClause(jt, left, right, on)
	q.sel.AddJoin(jc)

	// Make sure we remove the right-hand selection from the SelectStatement's
	// selections collection, since it's in a JOIN clause.
	q.sel.RemoveSelection(right)
	return q
}

func Select(items ...interface{}) *SelectQuery {
	scanner := scanner.New(types.DIALECT_UNKNOWN)
	sq := &SelectQuery{
		scanner: scanner,
	}
	sel := statement.NewSelect(make([]types.Projection, 0), nil, nil, nil, nil, nil, nil, nil)

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
			if len(innerSelClause.Selections()) == 1 {
				innerSel := innerSelClause.Selections()[0]
				switch innerSel.(type) {
				case *ast.DerivedTable:
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
					dt := innerSel.(*ast.DerivedTable)
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
					dt := ast.NewDerivedTable(derivedName, innerSelClause)
					selectionMap[dt] = true
					for _, p := range dt.DerivedColumns() {
						sel.AddProjection(p)
					}
					nDerived++
				}
			}
		case *ast.ColumnIdentifier:
			v := item.(*ast.ColumnIdentifier)
			// Set scanner's dialect based on supplied meta's dialect
			if v == nil {
				panic("specified a non-existent column")
			}
			sq.scanner.WithDialect(v.Schema().Dialect)
			sel.AddProjection(v)
			selectionMap[v.From()] = true
		case *ast.TableIdentifier:
			v := item.(*ast.TableIdentifier)
			// Set scanner's dialect based on supplied meta's dialect
			sq.scanner.WithDialect(v.Schema().Dialect)
			for _, c := range v.Projections() {
				sel.AddProjection(c)
			}
			selectionMap[v] = true
		case *ast.Function:
			v := item.(*ast.Function)
			sel.AddProjection(v)
			selectionMap[v.From()] = true
		default:
			// Everything else, make it a literal value projection, so, for
			// instance, a user can do SELECT 1, which is, technically
			// valid SQL.
			p := ast.NewValue(nil, item)
			sel.AddProjection(p)
		}
	}
	selections := make([]types.Selection, len(selectionMap))
	x := 0
	for sel, _ := range selectionMap {
		selections[x] = sel
		x++
	}
	sel.ReplaceSelections(selections)
	sq.sel = sel
	return sq
}
