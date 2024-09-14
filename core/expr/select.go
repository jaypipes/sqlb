//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package expr

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jaypipes/sqlb/core/fn"
	"github.com/jaypipes/sqlb/core/inspect"
	"github.com/jaypipes/sqlb/core/meta"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/grammar"
)

// Selection wraps a grammar.QuerySpecification, adding methods to inspect the
// wrapped query specifications projections/columns.
type Selection struct {
	qs    *grammar.QuerySpecification
	cs    *grammar.CursorSpecification
	cols  []types.Projection
	alias string
}

func (s *Selection) Name() string {
	return s.alias
}

// Alias returns the alias of the Selection
func (s *Selection) Alias() string {
	return s.alias
}

// AliasOrName returns the alias of the Selection
func (s *Selection) AliasOrName() string {
	return s.alias
}

// Projections returns a slice of Projection things referenced by the
// Selectable. The slice should be sorted by the Projection's name.
func (s *Selection) Projections() []types.Projection {
	cols := slices.Clone(s.cols)
	slices.SortFunc(cols, func(a, b types.Projection) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

func (s *Selection) Query() interface{} {
	if s.cs != nil {
		return s.cs
	}
	return s.qs
}

// QuerySpecification returns the object as a `*grammar.QuerySpecification`
func (s *Selection) QuerySpecification() *grammar.QuerySpecification {
	return s.qs
}

// TablePrimary returns the object as a `*grammar.TablePrimary`
func (s *Selection) TablePrimary() *grammar.TablePrimary {
	tp := &grammar.TablePrimary{
		DerivedTable: &grammar.DerivedTable{
			Subquery: grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: grammar.QueryExpressionBody{
						NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
							NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									SimpleTable: &grammar.SimpleTable{
										QuerySpecification: s.qs,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	if s.alias != "" {
		tp.Correlation = &grammar.Correlation{
			Name: s.alias,
		}
	}
	return tp
}

// TableReference returns the object as a `*grammar.TableReference`
func (s *Selection) TableReference() *grammar.TableReference {
	// A Selection that only has a QuerySpecification set on it is can be
	// considered a derived table (subquery in the FROM clause)
	return &grammar.TableReference{
		Primary: s.TablePrimary(),
	}
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (s *Selection) C(name string) types.Projection {
	for _, c := range s.cols {
		if strings.EqualFold(c.Name(), name) {
			return c
		}
	}
	return nil
}

// Select returns a QuerySpecification that produces a SELECT SQL statement for
// one or more items. Items can be a Table, a Column, a Function, another
// SELECT query, or even a literal value.
//
// Select panics if sqlb cannot compile the supplied arguments into a valid
// SELECT SQL query. This is intentional, as we want compile-time failures for
// invalid SQL construction.
func Select(
	items ...interface{},
) *Selection {
	cols := []types.Projection{}
	sels := []grammar.SelectSublist{}
	trefByName := map[string]grammar.TableReference{}
	nDerived := 0
	// For each scannable item we've received in the call, check what concrete
	// type they are and, depending on which type they are, either add them to
	// the returned SelectStatement's projections list or query the underlying
	// table metadata to generate a list of all columns in that table.
	for _, item := range items {
		switch item := item.(type) {
		case *Selection:
			// a derived table. The user has done something like:
			//
			// sqlb.Select(sqlb.Select(users).As("u"))
			//
			// and we need to produce the following SQL:
			//
			// SELECT u.id, u.name FROM (SELECT users.id, users.name FROM users) AS u
			derivedName := item.alias
			if derivedName == "" {
				derivedName = fmt.Sprintf("derived%d", nDerived)
				nDerived++
			}
			tp := grammar.TablePrimary{
				DerivedTable: &grammar.DerivedTable{
					Subquery: grammar.Subquery{
						QueryExpression: grammar.QueryExpression{
							Body: grammar.QueryExpressionBody{
								NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
									NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
										Primary: &grammar.NonJoinQueryPrimary{
											SimpleTable: &grammar.SimpleTable{
												QuerySpecification: item.qs,
											},
										},
									},
								},
							},
						},
					},
				},
				Correlation: &grammar.Correlation{
					Name: derivedName,
				},
			}
			selAsTableCols := map[string]types.Projection{}
			selAsTable := &DerivedTable{
				name: derivedName,
			}
			tref := grammar.TableReference{Primary: &tp}
			trefByName[derivedName] = tref
			// We need to project all columns from the supplied Selection's
			// QuerySpecification to the outer QuerySpecification.
			for _, c := range item.Projections() {
				outerCol := meta.NewColumn(selAsTable, c.Name())
				dc := grammar.DerivedColumn{
					ValueExpression: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: &grammar.ColumnReference{
									BasicIdentifierChain: &grammar.IdentifierChain{
										Identifiers: []string{derivedName, c.Name()},
									},
								},
							},
						},
					},
				}
				selAsTableCols[c.Name()] = outerCol
				sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
				cols = append(cols, outerCol)
			}
			selAsTable.columns = selAsTableCols
		case *grammar.Subquery:
			// Project all columns from the subquery to the outer
			// QuerySpecification
			body := item.QueryExpression.Body
			njqe := body.NonJoinQueryExpression
			if njqe == nil {
				panic("expected subquery to have non-nil non-join query expression")
			}
			njqt := njqe.NonJoinQueryTerm
			if njqt == nil {
				panic("expected subquery to have non-nil non-join query term")
			}
			njqp := njqt.Primary
			if njqp == nil {
				panic("expected subquery to have non-nil non-join query primary")
			}
			st := njqp.SimpleTable
			if st == nil {
				panic("expected subquery to have non-nil simple table")
			}
			qs := st.QuerySpecification
			if qs == nil {
				panic("expected subquery to have non-nil query specification")
			}
			// TODO(jaypipes): Determine if this is a SCALAR subquery or not...
			dc := grammar.DerivedColumn{
				ValueExpression: grammar.ValueExpression{
					Row: &grammar.RowValueExpression{
						Primary: &grammar.NonParenthesizedValueExpressionPrimary{
							ScalarSubquery: item,
						},
					},
				},
			}
			sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
		case types.Relation:
			tname := item.AliasOrName()
			tr := item.TableReference()
			trefByName[tname] = *tr
			for _, p := range item.Projections() {
				dc := p.DerivedColumn()
				sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
				cols = append(cols, p)
			}
		case types.Projection:
			dc := item.DerivedColumn()
			cols = append(cols, item)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			ref := item.References()
			if ref != nil {
				tname := ref.AliasOrName()
				tr := ref.TableReference()
				trefByName[tname] = *tr
			}
		default:
			// Everything else, make it a general literal value projection, so, for
			// instance, a user can do SELECT 1, which is, technically
			// valid SQL.
			dc := grammar.DerivedColumn{
				ValueExpression: grammar.ValueExpression{
					Row: &grammar.RowValueExpression{
						Primary: &grammar.NonParenthesizedValueExpressionPrimary{
							UnsignedValue: &grammar.UnsignedValueSpecification{
								UnsignedLiteral: &grammar.UnsignedLiteral{
									GeneralLiteral: &grammar.GeneralLiteral{
										Value: item,
									},
								},
							},
						},
					},
				},
			}
			sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
		}
	}

	if len(trefByName) == 0 {
		panic(
			"no entries in FROM clause. you must pass Select() at " +
				"least one element that references a table or subquery",
		)
	}

	trefs := make([]grammar.TableReference, 0, len(trefByName))
	for _, tref := range trefByName {
		trefs = append(trefs, tref)
	}
	return &Selection{
		qs: &grammar.QuerySpecification{
			SelectList: grammar.SelectList{
				Sublists: sels,
			},
			TableExpression: grammar.TableExpression{
				FromClause: grammar.FromClause{
					TableReferences: trefs,
				},
			},
		},
		cols: cols,
	}
}

// As returns a Selection as a DerivedTable
func (s *Selection) As(subqueryName string) types.Relation {
	if s.qs == nil {
		panic("cannot call As before Selection has a query specification")
	}
	return NewDerivedTable(subqueryName, s)
}

// Count applies a SELECT COUNT(*) to the Selection
func (s *Selection) Count() *Selection {
	if s.qs == nil {
		panic("called Count() on a nil Selection.")
	}
	dc := grammar.DerivedColumn{
		ValueExpression: grammar.ValueExpression{
			Row: &grammar.RowValueExpression{
				Primary: &grammar.NonParenthesizedValueExpressionPrimary{
					SetFunction: &grammar.SetFunctionSpecification{
						Aggregate: fn.Count().AggregateFunction,
					},
				},
			},
		},
	}
	s.qs.SelectList.Sublists = []grammar.SelectSublist{
		{
			DerivedColumn: &dc,
		},
	}
	return s
}

// Limit applies a LIMIT clause to the Selection (or a TOP N clause for T-SQL
// variants)
func (s *Selection) Limit(count int) *Selection {
	if s.cs != nil {
		s.cs.LimitClause = &grammar.LimitClause{
			Count: count,
		}
		return s
	}
	if s.qs == nil {
		panic("cannot call Limit() on a nil QuerySpecification")
	}
	cs := &grammar.CursorSpecification{
		QueryExpression: grammar.QueryExpression{
			Body: grammar.QueryExpressionBody{
				NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
					NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
						Primary: &grammar.NonJoinQueryPrimary{
							SimpleTable: &grammar.SimpleTable{
								QuerySpecification: s.qs,
							},
						},
					},
				},
			},
		},
		LimitClause: &grammar.LimitClause{
			Count: count,
		},
	}
	s.cs = cs
	return s
}

// Limit applies a LIMIT M OFFSET N clause to the Selection
func (s *Selection) LimitWithOffset(
	count int,
	offset int,
) *Selection {
	if s.cs != nil {
		s.cs.LimitClause = &grammar.LimitClause{
			Count:  count,
			Offset: &offset,
		}
		return s
	}
	if s.qs == nil {
		panic("cannot call Limit() on a nil QuerySpecification")
	}
	cs := &grammar.CursorSpecification{
		QueryExpression: grammar.QueryExpression{
			Body: grammar.QueryExpressionBody{
				NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
					NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
						Primary: &grammar.NonJoinQueryPrimary{
							SimpleTable: &grammar.SimpleTable{
								QuerySpecification: s.qs,
							},
						},
					},
				},
			},
		},
		LimitClause: &grammar.LimitClause{
			Count:  count,
			Offset: &offset,
		},
	}
	s.cs = cs
	return s
}

// Where adapts the Selection with a filtering expression, returning the
// Selection pointer to support method chaining.
func (s *Selection) Where(
	exprAny interface{},
) *Selection {
	if s.qs == nil {
		panic("cannot call Where() on a nil QuerySpecification")
	}
	bve := inspect.BooleanValueExpressionFromAny(exprAny)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			exprAny, exprAny,
		)
		panic(msg)
	}
	te := &s.qs.TableExpression
	if te.WhereClause != nil {
		te.WhereClause.SearchCondition = *And(&te.WhereClause.SearchCondition, bve)
	} else {
		te.WhereClause = &grammar.WhereClause{
			SearchCondition: *bve,
		}
	}
	s.qs.TableExpression = *te
	return s
}

// GroupBy adapts the Selection to group on the supplied columns, returning the
// adapted Selection itself to support method chaining.
func (s *Selection) GroupBy(
	cols ...interface{},
) *Selection {
	if s.qs == nil {
		panic("cannot call Where() on a nil QuerySpecification")
	}
	te := &s.qs.TableExpression
	if te.GroupByClause == nil {
		te.GroupByClause = &grammar.GroupByClause{}
	}
	ges := te.GroupByClause.GroupingElements
	if ges == nil {
		ges = []grammar.GroupingElement{}
	}
	for _, c := range cols {
		cr := inspect.ColumnReferenceFromAny(c)
		if cr == nil {
			msg := fmt.Sprintf(
				"could not convert %s(%T) to expected ColumnReference",
				c, c,
			)
			panic(msg)
		}
		ge := grammar.GroupingElement{
			OrdinaryGroupingSet: &grammar.OrdinaryGroupingSet{
				GroupingColumnReference: &grammar.GroupingColumnReference{
					ColumnReference: cr,
				},
			},
		}
		ges = append(ges, ge)
	}
	te.GroupByClause.GroupingElements = ges
	s.qs.TableExpression = *te
	return s
}

// Having adapts the Selection with the supplied filtering expression as an
// aggregate filter (a HAVING clause expression), returning the adapted
// Selection itself to support method chaining.
func (s *Selection) Having(
	exprAny interface{},
) *Selection {
	if s.qs == nil {
		panic("cannot call Having() on a nil QuerySpecification")
	}
	bve := inspect.BooleanValueExpressionFromAny(exprAny)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			exprAny, exprAny,
		)
		panic(msg)
	}
	te := &s.qs.TableExpression
	if te.HavingClause != nil {
		te.HavingClause.SearchCondition = *And(&te.HavingClause.SearchCondition, bve)
	} else {
		te.HavingClause = &grammar.HavingClause{
			SearchCondition: *bve,
		}
	}
	s.qs.TableExpression = *te
	return s
}

// OrderBy adds an ORDER BY to the Selection.
func (s *Selection) OrderBy(
	specAnys ...interface{},
) *Selection {
	if s.cs == nil {
		s.cs = &grammar.CursorSpecification{
			QueryExpression: grammar.QueryExpression{
				Body: grammar.QueryExpressionBody{
					NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
						NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
							Primary: &grammar.NonJoinQueryPrimary{
								SimpleTable: &grammar.SimpleTable{
									QuerySpecification: s.qs,
								},
							},
						},
					},
				},
			},
		}
	}
	specs := []grammar.SortSpecification{}
	for _, specAny := range specAnys {
		switch v := specAny.(type) {
		case *grammar.SortSpecification:
			specs = append(specs, *v)
		case grammar.SortSpecification:
			specs = append(specs, v)
		case *grammar.ValueExpression:
			specs = append(specs, grammar.SortSpecification{Key: *v})
		case grammar.ValueExpression:
			specs = append(specs, grammar.SortSpecification{Key: v})
		default:
			ve := inspect.ValueExpressionFromAny(specAny)
			if ve == nil {
				msg := fmt.Sprintf(
					"could not convert %s(%T) to expected ValueExpression",
					specAny, specAny,
				)
				panic(msg)
			}
			specs = append(specs, grammar.SortSpecification{Key: *ve})
		}
	}
	if s.cs.OrderByClause == nil {
		s.cs.OrderByClause = &grammar.OrderByClause{
			SortSpecifications: []grammar.SortSpecification{},
		}
	}
	s.cs.OrderByClause.SortSpecifications = append(
		s.cs.OrderByClause.SortSpecifications, specs...,
	)
	return s
}
