//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package api

import (
	"fmt"

	"github.com/jaypipes/sqlb/grammar"
	"github.com/samber/lo"
)

// Join adapts the Selection after joining the FromClause's last TableReference
// to the first parameter which must be convertible to a TableReference.
func (s *Selection) Join(
	rightAny interface{},
	onAny interface{},
) *Selection {
	if s.qs == nil {
		panic("attempt to join against nil query specification")
	}
	if len(s.qs.TableExpression.FromClause.TableReferences) == 0 {
		msg := "attempt to join against nothing. before calling Join() " +
			"first call Select()"
		panic(msg)
	}
	right := TableReferenceFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"attempted join on invalid type %s(%T)",
			rightAny, rightAny,
		)
		panic(msg)
	}
	on := BooleanValueExpressionFromAny(onAny)
	if on == nil {
		msg := fmt.Sprintf(
			"invalid join condition %s(%T)",
			onAny, onAny,
		)
		panic(msg)
	}
	return s.doJoin(grammar.JoinTypeInner, right, on)
}

// OuterJoin adapts the Selection after left-joining the FromClause's last
// TableReference to the first parameter which must be convertible to a
// TableReference.
func (s *Selection) OuterJoin(
	rightAny interface{},
	onAny interface{},
) *Selection {
	if s.qs == nil {
		panic("attempt to join against nil query specification")
	}
	if len(s.qs.TableExpression.FromClause.TableReferences) == 0 {
		msg := "attempt to join against nothing. before calling OuterJoin() " +
			"first call Select()"
		panic(msg)
	}
	right := TableReferenceFromAny(rightAny)
	if right == nil {
		msg := fmt.Sprintf(
			"attempted join on invalid type %s(%T)",
			rightAny, rightAny,
		)
		panic(msg)
	}
	on := BooleanValueExpressionFromAny(onAny)
	if on == nil {
		msg := fmt.Sprintf(
			"invalid join condition %s(%T)",
			onAny, onAny,
		)
		panic(msg)
	}
	return s.doJoin(grammar.JoinTypeLeftOuter, right, on)
}

func (s *Selection) doJoin(
	joinType grammar.JoinType,
	right *grammar.TableReference,
	on *grammar.BooleanValueExpression,
) *Selection {
	// We have to remove all referenced TablePrimary TableReferences from the
	// existing QuerySpecification's list of TableReferences because these
	// named TablePrimaries will be output when the builder returns the
	// JoinedTable SQL string.
	//
	// Let's walk through what we might occur with a simple single table join.
	// If the user does this:
	//
	// sel := sqlb.Select(users.C("name"), articles.C("title"))
	//
	// the `sel` variable would contain a QuerySpecification that looks like
	// this:
	//
	// SelectList:
	//   SelectSublist:
	//   - ColumnReference:
	//     - users
	//     - name
	//   - ColumnReference:
	//     - articles
	//     - title
	// TableExpression:
	//   FromClause:
	//     TableReferences:
	//     - Primary:
	//         TableName: users
	//     - Primary:
	//         TableName: articles
	//
	// which, if printed by the sqlb Builder, would produce the following SQL:
	//
	// SELECT users.id, articles.title FROM users, articles
	//
	// which, of course, is a cartesian product of the users and articles
	// tables.
	//
	// afterwards, the user does this:
	//
	// sel = sel.Join(articles, sqlb.Equal(users.C("id"), articles.C("author"))
	//
	// and we need to end up with a new QuerySpecification that looks like
	// this:
	//
	// SelectList:
	//   SelectSublist:
	//   - ColumnReference:
	//     - users
	//     - name
	//   - ColumnReference:
	//     - articles
	//     - title
	// TableExpression:
	//   FromClause:
	//     TableReferences:
	//     - JoinedTable:
	//         QualifiedJoin:
	//			 JoinType: INNER
	//           Left:
	//             Primary:
	//               TableName: users
	//           Right:
	//             Primary:
	//               TableName: articles
	//           On:
	//             Unary:
	//               Unary:
	//                 Test:
	//                   Primary:
	//                     Predicate:
	//                       Comparison:
	//                         Operation: EQUALS
	//                         A:
	//                           NonParenthesizedValueExpressionPrimary:
	//                             ColumnReference:
	//                               - users
	//                               - id
	//                         B:
	//                           NonParenthesizedValueExpressionPrimary:
	//                             ColumnReference:
	//                               - articles
	//                               - author
	var rightTR *grammar.TableReference
	if right.Primary != nil {
		rp := right.Primary
		var search string
		if rp.Correlation != nil {
			// ON condition might reference a table or derived table by its
			// alias...
			search = rp.Correlation.Name
		} else if rp.TableName != nil {
			search = *rp.TableName
		} else if rp.QueryName != nil {
			search = *rp.QueryName
		}
		rightTR = TableReferenceByName(
			s.qs.TableExpression.FromClause.TableReferences,
			search,
		)
		if rightTR == nil {
			// We are joining a table that has not yet been referenced in the
			// query specification
			trefs := s.qs.TableExpression.FromClause.TableReferences
			trefs = append(trefs, *right)
			s.qs.TableExpression.FromClause.TableReferences = trefs
			rightTR = right
		}
	} else if right.Joined != nil {
		/*
			j := right.Joined
			if j.Qualified != nil {
				trefs := []grammar.TableReference{
					j.Qualified.Left,
					j.Qualified.Right,
				}
				rightTR = TableReferenceByName(search)
			} else if j.Natural != nil {

			} else if j.Union != nil {

			} else if j.Cross != nil {

			}
		*/
	}
	rightname := ""
	if rightTR.Primary != nil {
		rp := rightTR.Primary
		if rp.Correlation != nil {
			// ON condition might reference a table or derived table by its
			// alias...
			rightname = rp.Correlation.Name
		} else if rp.TableName != nil {
			rightname = *rp.TableName
		} else if rp.QueryName != nil {
			rightname = *rp.QueryName
		}
	}
	referreds := lo.Uniq(ReferredFromBooleanValueExpression(on))
	// remove the referrant so we can just look for the reference on the left
	// side of the join
	referreds = lo.Without(referreds, rightname)
	// Find the table references that the ON condition refers to and replace
	// any table primaries with a JoinedTable containing the joined right table
	// reference.
	updatedTRefs := []grammar.TableReference{}
	for _, referred := range referreds {
		leftTR := TableReferenceByName(s.qs.TableExpression.FromClause.TableReferences, referred)
		if leftTR != nil {
			jt := &grammar.JoinedTable{
				Qualified: &grammar.QualifiedJoin{
					Type:  joinType,
					Right: *rightTR,
					Left:  *leftTR,
					On:    *on,
				},
			}
			updatedTRefs = append(updatedTRefs, grammar.TableReference{
				Joined: jt,
			})
		} else {
			msg := fmt.Sprintf(
				"ON condition refers to a table '%s' that is not "+
					"in a JOIN or FROM clause",
				referred,
			)
			panic(msg)
		}
	}

	// Now we replace the right table reference with a new one representing the
	// joined table
	s.qs.TableExpression.FromClause.TableReferences = updatedTRefs
	return s
}
