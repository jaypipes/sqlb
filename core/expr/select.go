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
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/meta"
	"github.com/jaypipes/sqlb/core/types"
	"github.com/jaypipes/sqlb/internal/inspect"
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

// Projections returns a slice of Projections referenced by the Selection. The
// slice is sorted by the Projection's name.
func (s *Selection) Projections() []types.Projection {
	cols := slices.Clone(s.cols)
	slices.SortFunc(cols, func(a, b types.Projection) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return cols
}

// Query returns the CursorSpecification if set, otherwise returns the wrapped
// QuerySpecification.
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
						NonJoin: &grammar.NonJoinQueryExpression{
							NonJoin: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									Simple: &grammar.SimpleTable{
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
// invalid SQL construction and we want Select() to be chainable with other
// Select() calls.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `SelectE` function which returns a checkable `error` object.
func Select(
	items ...interface{},
) *Selection {
	s, err := SelectE(items...)
	if err != nil {
		panic(err)
	}
	return s
}

// SelectE returns a QuerySpecification that produces a SELECT SQL statement for
// one or more items. Items can be a Table, a Column, a Function, another
// SELECT query, or even a literal value.
//
// If sqlb cannot compile the supplied arguments into a valid SELECT SQL query,
// SelectE returns an error.
func SelectE(
	items ...interface{},
) (*Selection, error) {
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
								NonJoin: &grammar.NonJoinQueryExpression{
									NonJoin: &grammar.NonJoinQueryTerm{
										Primary: &grammar.NonJoinQueryPrimary{
											Simple: &grammar.SimpleTable{
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
					Value: grammar.ValueExpression{
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
			njqe := body.NonJoin
			if njqe == nil {
				return nil, fmt.Errorf(
					"expected subquery to have non-nil non-join " +
						"query expression",
				)
			}
			njqt := njqe.NonJoin
			if njqt == nil {
				return nil, fmt.Errorf(
					"expected subquery to have non-nil non-join query term",
				)
			}
			njqp := njqt.Primary
			if njqp == nil {
				return nil, fmt.Errorf(
					"expected subquery to have non-nil non-join " +
						"query primary",
				)
			}
			st := njqp.Simple
			if st == nil {
				return nil, fmt.Errorf(
					"expected subquery to have non-nil simple table",
				)
			}
			qs := st.QuerySpecification
			if qs == nil {
				return nil, fmt.Errorf(
					"expected subquery to have non-nil query specification",
				)
			}
			// TODO(jaypipes): Determine if this is a SCALAR subquery or not...
			dc := grammar.DerivedColumn{
				Value: grammar.ValueExpression{
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
				Value: grammar.ValueExpression{
					Row: &grammar.RowValueExpression{
						Primary: &grammar.NonParenthesizedValueExpressionPrimary{
							UnsignedValue: &grammar.UnsignedValueSpecification{
								UnsignedLiteral: &grammar.UnsignedLiteral{
									General: &grammar.GeneralLiteral{
										Value: item,
									},
								},
							},
						},
					},
				},
			}
			fmt.Println("ADDED", fmt.Sprintf("%+v", item))
			sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
		}
	}

	if len(trefByName) == 0 {
		return nil, fmt.Errorf(
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
				From: grammar.FromClause{
					TableReferences: trefs,
				},
			},
		},
		cols: cols,
	}, nil
}

// As returns a Selection as a DerivedTable.
//
// As panics if the Selection has not had a query specification set yet. This
// is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of As() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `AsE` function which returns a checkable `error` object.
func (s *Selection) As(subqueryName string) types.Relation {
	r, err := s.AsE(subqueryName)
	if err != nil {
		panic(err)
	}
	return r
}

// AsE returns a Selection as a DerivedTable. If the Selection has not yet had
// its query specification set, AsE returns an error.
func (s *Selection) AsE(subqueryName string) (types.Relation, error) {
	if s == nil || s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call As before Selection has a query specification",
		)
	}
	return NewDerivedTable(subqueryName, s), nil
}

// Count applies a SELECT COUNT(*) to the Selection
//
// Count panics if the Selection has not had a query specification set yet. This
// is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of Count() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `CountE` function which returns a checkable `error` object.
func (s *Selection) Count() *Selection {
	res, err := s.CountE()
	if err != nil {
		panic(err)
	}
	return res
}

// CountE applies a SELECT COUNT(*) to the Selection and returns the Selection.
// If the Selection has not had its query specification set, CountE returns an
// error.
func (s *Selection) CountE() (*Selection, error) {
	if s == nil || s.qs == nil {
		return nil, fmt.Errorf("called Count() on a nil Selection.")
	}
	dc := grammar.DerivedColumn{
		Value: grammar.ValueExpression{
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
	return s, nil
}

// Limit applies a LIMIT clause to the Selection (or a TOP N clause for T-SQL
// variants)
//
// Limit panics if the Selection has not had a query specification set yet. This
// is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of Limit() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `LimitE` function which returns a checkable `error` object.
func (s *Selection) Limit(count int) *Selection {
	res, err := s.LimitE(count)
	if err != nil {
		panic(err)
	}
	return res
}

// LimitE applies a LIMIT clause to the Selection (or a TOP N clause for T-SQL
// variants). If the Selection has not had its query specification set, LimitE
// returns an error.
func (s *Selection) LimitE(count int) (*Selection, error) {
	if s == nil {
		return nil, fmt.Errorf(
			"cannot call Limit() on a nil Selection",
		)
	}
	if s.cs != nil {
		s.cs.Limit = &grammar.LimitClause{
			Count: count,
		}
		return s, nil
	}
	if s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call Limit() on a nil QuerySpecification",
		)
	}
	cs := &grammar.CursorSpecification{
		Query: grammar.QueryExpression{
			Body: grammar.QueryExpressionBody{
				NonJoin: &grammar.NonJoinQueryExpression{
					NonJoin: &grammar.NonJoinQueryTerm{
						Primary: &grammar.NonJoinQueryPrimary{
							Simple: &grammar.SimpleTable{
								QuerySpecification: s.qs,
							},
						},
					},
				},
			},
		},
		Limit: &grammar.LimitClause{
			Count: count,
		},
	}
	s.cs = cs
	return s, nil
}

// LimitWithOffset applies a LIMIT M OFFSET N clause to the Selection
//
// LimitWithOffset panics if the Selection has not had a query specification
// set yet. This is intentional, as we want compile-time failures for invalid
// SQL construction and we want the result of LimitWithOffset() to be chainable
// with other Selection methods and be usable as an input to the Select()
// function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `LimitWithOffsetE` function which returns a checkable `error`
// object.
func (s *Selection) LimitWithOffset(
	count int,
	offset int,
) *Selection {
	res, err := s.LimitWithOffsetE(count, offset)
	if err != nil {
		panic(err)
	}
	return res
}

// LimitWithOffset applies a LIMIT M OFFSET N clause to the Selection. If the
// Selection has not had its query specification set, LimitWithOffsetE returns
// an error.
func (s *Selection) LimitWithOffsetE(
	count int,
	offset int,
) (*Selection, error) {
	if s == nil {
		return nil, fmt.Errorf(
			"cannot call LimitWithOffset() on a nil Selection",
		)
	}
	if s.cs != nil {
		s.cs.Limit = &grammar.LimitClause{
			Count:  count,
			Offset: &offset,
		}
		return s, nil
	}
	if s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call LimitWithOffset() on a nil QuerySpecification",
		)
	}
	cs := &grammar.CursorSpecification{
		Query: grammar.QueryExpression{
			Body: grammar.QueryExpressionBody{
				NonJoin: &grammar.NonJoinQueryExpression{
					NonJoin: &grammar.NonJoinQueryTerm{
						Primary: &grammar.NonJoinQueryPrimary{
							Simple: &grammar.SimpleTable{
								QuerySpecification: s.qs,
							},
						},
					},
				},
			},
		},
		Limit: &grammar.LimitClause{
			Count:  count,
			Offset: &offset,
		},
	}
	s.cs = cs
	return s, nil
}

// Where adapts the Selection with a filtering expression, returning the
// Selection pointer to support method chaining.
//
// Where panics if the Selection has not had a query specification set yet.
// This is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of Where() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `WhereE` function which returns a checkable `error`
// object.
func (s *Selection) Where(
	exprAny interface{},
) *Selection {
	res, err := s.WhereE(exprAny)
	if err != nil {
		panic(err)
	}
	return res
}

// WhereE adapts the Selection with a filtering expression, returning the
// Selection pointer to support method chaining. If the Selection has not had
// its query specification set or the supplied parameters cannot be converted
// to BooleanValueExpressions, WhereE returns an error.
func (s *Selection) WhereE(
	exprAny interface{},
) (*Selection, error) {
	if s == nil || s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call Where() on a nil QuerySpecification",
		)
	}
	bve := inspect.BooleanValueExpressionFromAny(exprAny)
	if bve == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			exprAny, exprAny,
		)
	}
	te := &s.qs.TableExpression
	if te.Where != nil {
		te.Where.Search = *And(&te.Where.Search, bve)
	} else {
		te.Where = &grammar.WhereClause{
			Search: *bve,
		}
	}
	s.qs.TableExpression = *te
	return s, nil
}

// GroupBy adapts the Selection to group on the supplied columns, returning the
// adapted Selection itself to support method chaining.
//
// GroupBy panics if the Selection has not had a query specification set yet.
// This is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of GroupBy() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `GroupByE` function which returns a checkable `error` object.
func (s *Selection) GroupBy(
	cols ...interface{},
) *Selection {
	res, err := s.GroupByE(cols...)
	if err != nil {
		panic(err)
	}
	return res
}

// GroupByE adapts the Selection to group on the supplied columns, returning
// the adapted Selection itself to support method chaining. If the Selection
// has not had its query specification set or the supplied parameters cannot be
// converted to ColumnReferences, WhereE returns an error.
func (s *Selection) GroupByE(
	cols ...interface{},
) (*Selection, error) {
	if s == nil || s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call GroupBy() on a nil QuerySpecification",
		)
	}
	te := &s.qs.TableExpression
	if te.GroupBy == nil {
		te.GroupBy = &grammar.GroupByClause{}
	}
	ges := te.GroupBy.GroupingElements
	if ges == nil {
		ges = []grammar.GroupingElement{}
	}
	for _, c := range cols {
		cr := inspect.ColumnReferenceFromAny(c)
		if cr == nil {
			return nil, fmt.Errorf(
				"could not convert %s(%T) to expected ColumnReference",
				c, c,
			)
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
	te.GroupBy.GroupingElements = ges
	s.qs.TableExpression = *te
	return s, nil
}

// Having adapts the Selection with the supplied filtering expression as an
// aggregate filter (a HAVING clause expression), returning the adapted
// Selection itself to support method chaining.
//
// Having panics if the Selection has not had a query specification set yet.
// This is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of Having() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `HavingE` function which returns a checkable `error` object.
func (s *Selection) Having(
	exprAny interface{},
) *Selection {
	res, err := s.HavingE(exprAny)
	if err != nil {
		panic(err)
	}
	return res
}

// HavingE adapts the Selection with the supplied filtering expression as an
// aggregate filter (a HAVING clause expression), returning the adapted
// Selection itself to support method chaining. If the Selection has not had
// its query specification set or the supplied parameters cannot be converted
// to BooleanValueExpressions, HavingE returns an error.
func (s *Selection) HavingE(
	exprAny interface{},
) (*Selection, error) {
	if s == nil || s.qs == nil {
		return nil, fmt.Errorf(
			"cannot call Having() on a nil QuerySpecification",
		)
	}
	bve := inspect.BooleanValueExpressionFromAny(exprAny)
	if bve == nil {
		return nil, fmt.Errorf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			exprAny, exprAny,
		)
	}
	te := &s.qs.TableExpression
	if te.Having != nil {
		te.Having.Search = *And(&te.Having.Search, bve)
	} else {
		te.Having = &grammar.HavingClause{
			Search: *bve,
		}
	}
	s.qs.TableExpression = *te
	return s, nil
}

// OrderBy adds an ORDER BY to the Selection.
//
// OrderBy panics if the Selection has not had a query specification set yet.
// This is intentional, as we want compile-time failures for invalid SQL
// construction and we want the result of OrderBy() to be chainable with other
// Selection methods and be usable as an input to the Select() function.
//
// If you are constructing SQL expressions dynamically with user-supplied
// input, use the `OrderByE` function which returns a checkable `error` object.
func (s *Selection) OrderBy(
	specAnys ...interface{},
) *Selection {
	res, err := s.OrderByE(specAnys...)
	if err != nil {
		panic(err)
	}
	return res
}

// OrderBy adds an ORDER BY to the Selection. If the Selection has not had its
// query specification set or the supplied parameters cannot be converted to
// ValueExpressions, OrderByE returns an error.
func (s *Selection) OrderByE(
	specAnys ...interface{},
) (*Selection, error) {
	if s == nil {
		return nil, fmt.Errorf(
			"cannot call OrderBy() on a nil Selection",
		)
	}
	if s.cs == nil {
		s.cs = &grammar.CursorSpecification{
			Query: grammar.QueryExpression{
				Body: grammar.QueryExpressionBody{
					NonJoin: &grammar.NonJoinQueryExpression{
						NonJoin: &grammar.NonJoinQueryTerm{
							Primary: &grammar.NonJoinQueryPrimary{
								Simple: &grammar.SimpleTable{
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
				return nil, fmt.Errorf(
					"could not convert %s(%T) to expected ValueExpression",
					specAny, specAny,
				)
			}
			specs = append(specs, grammar.SortSpecification{Key: *ve})
		}
	}
	if s.cs.OrderBy == nil {
		s.cs.OrderBy = &grammar.OrderByClause{
			SortSpecifications: []grammar.SortSpecification{},
		}
	}
	s.cs.OrderBy.SortSpecifications = append(
		s.cs.OrderBy.SortSpecifications, specs...,
	)
	return s, nil
}
