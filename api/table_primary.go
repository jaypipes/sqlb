//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"github.com/jaypipes/sqlb/grammar"
)

// NameAndTablePrimaryFromReferred evaluates the supplied interface argument
// representing a referred table or derived table and returns a *TablePrimary
// and the name of the referrent if the supplied argument can be converted into
// a TablePrimary, or nil if the conversion cannot be done.
func NameAndTablePrimaryFromReferred(
	subject interface{},
) (string, *grammar.TablePrimary) {
	tname := ""
	tp := &grammar.TablePrimary{}
	t, ok := subject.(*Table)
	if ok {
		tname = t.Name()
		tp.TableName = &tname
		if t.alias != "" {
			tp.Correlation = &grammar.Correlation{
				Name: t.Alias(),
			}
		}
	} else {
		// The column is from a derived table
		dt := subject.(*DerivedTable)
		tname = dt.Name()
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
	return tname, tp
}
