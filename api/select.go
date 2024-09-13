//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"
	"slices"
	"strings"

	"github.com/jaypipes/sqlb/grammar"
)

// Selection wraps a grammar.QuerySpecification, adding methods to inspect the
// wrapped query specifications projections/columns.
type Selection struct {
	qs    *grammar.QuerySpecification
	cs    *grammar.CursorSpecification
	cols  []*Column
	alias string
}

// ColumnsSorted returns a slice of the Selection's Columns sorted by Column
// name.
func (s *Selection) ColumnsSorted() []*Column {
	cols := slices.Clone(s.cols)
	slices.SortFunc(cols, func(a, b *Column) int {
		return strings.Compare(a.name, b.name)
	})
	return cols
}

func (s *Selection) Query() interface{} {
	if s.cs != nil {
		return s.cs
	}
	return s.qs
}

// C returns a pointer to a Column with a name matching the supplied string, or
// nil if no such column is known
//
// The name matching is done using case-insensitive matching, since this is how
// the SQL standard works for identifiers and symbols (even though Microsoft
// SQL Server uses case-sensitive identifier names).
func (s *Selection) C(name string) *Column {
	for _, c := range s.cols {
		if strings.EqualFold(c.name, name) {
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
	cols := []*Column{}
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
			selAsTableCols := map[string]*Column{}
			selAsTable := &DerivedTable{
				name: derivedName,
			}
			tref := grammar.TableReference{Primary: &tp}
			trefByName[derivedName] = tref
			// We need to project all columns from the supplied Selection's
			// QuerySpecification to the outer QuerySpecification.
			for _, c := range item.ColumnsSorted() {
				outerCol := &Column{
					t:    selAsTable,
					name: c.name,
				}
				dc := grammar.DerivedColumn{
					ValueExpression: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: &grammar.ColumnReference{
									BasicIdentifierChain: &grammar.IdentifierChain{
										Identifiers: []string{derivedName, c.name},
									},
								},
							},
						},
					},
				}
				selAsTableCols[c.name] = outerCol
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
		case *Column:
			if item == nil {
				panic("specified a non-existent column")
			}
			tname := item.TableName()
			dc := grammar.DerivedColumn{
				ValueExpression: grammar.ValueExpression{
					Row: &grammar.RowValueExpression{
						Primary: &grammar.NonParenthesizedValueExpressionPrimary{
							ColumnReference: &grammar.ColumnReference{
								BasicIdentifierChain: &grammar.IdentifierChain{
									Identifiers: []string{tname, item.name},
								},
							},
						},
					},
				},
			}
			if item.alias != "" {
				dc.As = &item.alias
			}
			sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
			cols = append(cols, item)
			tnameNoAlias := item.TableNameNoAlias()
			tp := &grammar.TablePrimary{}
			_, ok := item.t.(*Table)
			if ok {
				tp.TableName = &tnameNoAlias
				if item.TableAlias() != "" {
					tp.Correlation = &grammar.Correlation{
						Name: item.TableAlias(),
					}
				}
			} else {
				// The column is from a derived table
				dt := item.t.(*DerivedTable)
				tp.DerivedTable = &grammar.DerivedTable{
					Subquery: grammar.Subquery{
						QueryExpression: grammar.QueryExpression{
							Body: grammar.QueryExpressionBody{
								NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
									NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
										Primary: &grammar.NonJoinQueryPrimary{
											SimpleTable: &grammar.SimpleTable{
												QuerySpecification: dt.Query(),
											},
										},
									},
								},
							},
						},
					},
				}
				// Derived tables are always named/aliased
				tp.Correlation = &grammar.Correlation{
					Name: dt.Name(),
				}
			}
			tr := grammar.TableReference{Primary: tp}
			trefByName[tname] = tr
		case *Table:
			if item == nil {
				panic("specified a non-existent table")
			}
			tname := item.name
			if item.alias != "" {
				tname = item.alias
			}
			for _, c := range item.ColumnsSorted() {
				dc := grammar.DerivedColumn{
					ValueExpression: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: &grammar.ColumnReference{
									BasicIdentifierChain: &grammar.IdentifierChain{
										Identifiers: []string{tname, c.name},
									},
								},
							},
						},
					},
				}
				sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
				cols = append(cols, c)
				// The table reference should point to the original table but
				// keep any alias, which is why we set the Correlation on the
				// TablePrimary here and use the tname (which is set to the
				// table's Alias, if any, above) as the table references map
				// key.
				tr := grammar.TableReference{
					Primary: &grammar.TablePrimary{
						TableName: &item.name,
					},
				}
				if item.alias != "" {
					tr.Primary.Correlation = &grammar.Correlation{
						Name: item.alias,
					}
				}
				trefByName[tname] = tr
			}
		case *DerivedTable:
			tname := item.name
			for _, c := range item.ColumnsSorted() {
				dc := grammar.DerivedColumn{
					ValueExpression: grammar.ValueExpression{
						Row: &grammar.RowValueExpression{
							Primary: &grammar.NonParenthesizedValueExpressionPrimary{
								ColumnReference: &grammar.ColumnReference{
									BasicIdentifierChain: &grammar.IdentifierChain{
										Identifiers: []string{tname, c.name},
									},
								},
							},
						},
					},
				}
				sels = append(sels, grammar.SelectSublist{DerivedColumn: &dc})
				cols = append(cols, c)
				// The table reference should point to the original table but
				// keep any alias, which is why we set the Correlation on the
				// TablePrimary here and use the tname (which is set to the
				// table's Alias, if any, above) as the table references map
				// key.
				tp := &grammar.TablePrimary{}
				tp.DerivedTable = &grammar.DerivedTable{
					Subquery: grammar.Subquery{
						QueryExpression: grammar.QueryExpression{
							Body: grammar.QueryExpressionBody{
								NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
									NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
										Primary: &grammar.NonJoinQueryPrimary{
											SimpleTable: &grammar.SimpleTable{
												QuerySpecification: item.Query(),
											},
										},
									},
								},
							},
						},
					},
				}
				// Derived tables are always named/aliased
				tp.Correlation = &grammar.Correlation{
					Name: item.Name(),
				}
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *AggregateFunction:
			if item == nil {
				panic("specified a non-existent aggregate function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *SubstringFunction:
			if item == nil {
				panic("specified a non-existent substring function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *RegexSubstringFunction:
			if item == nil {
				panic("specified a non-existent regex substring function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *FoldFunction:
			if item == nil {
				panic("specified a non-existent fold function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *TranscodingFunction:
			if item == nil {
				panic("specified a non-existent transcoding function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
			}
		case *TransliterationFunction:
			if item == nil {
				panic("specified a non-existent tranliteration function")
			}
			dc := DerivedColumnFromAnyAndAlias(
				item, item.alias,
			)
			sels = append(sels, grammar.SelectSublist{DerivedColumn: dc})
			if item.Referred != nil {
				tname, tp := NameAndTablePrimaryFromReferred(item.Referred)
				tr := grammar.TableReference{Primary: tp}
				trefByName[tname] = tr
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
func (s *Selection) As(subqueryName string) *DerivedTable {
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
						Aggregate: Count().AggregateFunction,
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
	expr interface{},
) *Selection {
	if s.qs == nil {
		panic("cannot call Where() on a nil QuerySpecification")
	}
	bve := BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	te := &s.qs.TableExpression
	if te.WhereClause != nil {
		te.WhereClause.SearchCondition = And(&te.WhereClause.SearchCondition, bve)
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
		cr := ColumnReferenceFromAny(c)
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
	expr interface{},
) *Selection {
	if s.qs == nil {
		panic("cannot call Having() on a nil QuerySpecification")
	}
	bve := BooleanValueExpressionFromAny(expr)
	if bve == nil {
		msg := fmt.Sprintf(
			"could not convert %s(%T) to expected BooleanValueExpression",
			expr, expr,
		)
		panic(msg)
	}
	te := &s.qs.TableExpression
	if te.HavingClause != nil {
		te.HavingClause.SearchCondition = And(&te.HavingClause.SearchCondition, bve)
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
			ve := ValueExpressionFromAny(specAny)
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
