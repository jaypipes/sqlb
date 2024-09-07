//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"github.com/jaypipes/sqlb/grammar"
)

// TableReferenceFromAny evaluates the supplied interface argument and returns
// a *TableReference if the supplied argument can be converted into a
// TableReference, or nil if the conversion cannot be done.
func TableReferenceFromAny(
	subject interface{},
) *grammar.TableReference {
	switch v := subject.(type) {
	case *grammar.TableReference:
		return v
	case grammar.TableReference:
		return &v
	case *Selection:
		// A Selection that only has a QuerySpecification set on it is can be
		// considered a derived table (subquery in the FROM clause)
		tr := &grammar.TableReference{
			Primary: &grammar.TablePrimary{
				DerivedTable: &grammar.DerivedTable{
					Subquery: grammar.Subquery{
						QueryExpression: grammar.QueryExpression{
							Body: grammar.QueryExpressionBody{
								NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
									NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
										Primary: &grammar.NonJoinQueryPrimary{
											SimpleTable: &grammar.SimpleTable{
												QuerySpecification: v.qs,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		if v.alias != "" {
			tr.Primary.Correlation = &grammar.Correlation{
				Name: v.alias,
			}
		}
		return tr
	case *Table:
		tr := &grammar.TableReference{
			Primary: &grammar.TablePrimary{
				TableName: &v.name,
			},
		}
		if v.alias != "" {
			tr.Primary.Correlation = &grammar.Correlation{
				Name: v.alias,
			}
		}
		return tr
	case *DerivedTable:
		tp := &grammar.TablePrimary{}
		tp.DerivedTable = &grammar.DerivedTable{
			Subquery: grammar.Subquery{
				QueryExpression: grammar.QueryExpression{
					Body: grammar.QueryExpressionBody{
						NonJoinQueryExpression: &grammar.NonJoinQueryExpression{
							NonJoinQueryTerm: &grammar.NonJoinQueryTerm{
								Primary: &grammar.NonJoinQueryPrimary{
									SimpleTable: &grammar.SimpleTable{
										QuerySpecification: v.Query(),
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
			Name: v.Name(),
		}
		return &grammar.TableReference{Primary: tp}
	case *grammar.TablePrimary:
		return &grammar.TableReference{Primary: v}
	case grammar.TablePrimary:
		return &grammar.TableReference{Primary: &v}
	case *grammar.JoinedTable:
		return &grammar.TableReference{Joined: v}
	case grammar.JoinedTable:
		return &grammar.TableReference{Joined: &v}
	}
	return nil
}

// TableReferenceByName returns a pointer to the TableReference that has a name
// or alias/correlation ID matching the supplied string.
func TableReferenceByName(
	refs []grammar.TableReference,
	search string,
) *grammar.TableReference {
	for _, ref := range refs {
		if ref.Primary != nil {
			p := ref.Primary
			if p.Correlation != nil && p.Correlation.Name == search {
				return &ref
			} else if p.TableName != nil && *p.TableName == search {
				return &ref
			} else if p.QueryName != nil && *p.QueryName == search {
				return &ref
			}
		} else if ref.Joined != nil {
			jt := ref.Joined
			if jt.Qualified != nil {
				found := TableReferenceByName(
					[]grammar.TableReference{
						jt.Qualified.Left,
						jt.Qualified.Right,
					},
					search,
				)
				if found != nil {
					return &ref
				}
			}
		}
	}
	return nil
}
