//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

// ArgCount examines the supplied element and recursively determines the number
// of query arguments represented in the element. The value of the `count`
// pointer to int is incremented recursively.
func ArgCount(target interface{}, count *int) {
	switch el := target.(type) {
	case *grammar.SelectList:
		if !el.Asterisk {
			for _, s := range el.Sublists {
				ArgCount(s, count)
			}
		}
	case grammar.SelectSublist:
		if !el.Asterisk {
			ArgCount(el.DerivedColumn, count)
		}
	case *grammar.ValueExpression:
		if el.Common != nil {
			ArgCount(el.Common, count)
		} else if el.Boolean != nil {
			ArgCount(el.Boolean, count)
		} else if el.Row != nil {
			ArgCount(el.Row.Primary, count)
		}
	case *grammar.ValueExpressionPrimary:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Parenthesized != nil {
			ArgCount(el.Parenthesized, count)
		}
	case *grammar.Factor:
		ArgCount(&el.Primary, count)
	case *grammar.NumericPrimary:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Function != nil {
			ArgCount(el.Function, count)
		}
	case *grammar.Term:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		} else if el.MultiplyDivide != nil {
			ArgCount(&el.MultiplyDivide.Left, count)
			ArgCount(&el.MultiplyDivide.Right, count)
		}
	case *grammar.NumericValueExpression:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		} else if el.AddSubtract != nil {
			ArgCount(&el.AddSubtract.Left, count)
			ArgCount(&el.AddSubtract.Right, count)
		}
	case *grammar.DatetimeValueExpression:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		} else if el.AddInterval != nil {
			ArgCount(&el.AddInterval.Left, count)
			ArgCount(&el.AddInterval.Right, count)
		} else if el.AddSubtract != nil {
			ArgCount(&el.AddSubtract.Left, count)
			ArgCount(&el.AddSubtract.Right, count)
		}
	case *grammar.DatetimeFactor:
		ArgCount(&el.Primary, count)
	case *grammar.DatetimePrimary:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Function != nil {
			ArgCount(el.Function, count)
		}
	case *grammar.DatetimeTerm:
		ArgCount(&el.Factor, count)
	case *grammar.BooleanValueExpression:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		}
		if el.OrLeft != nil {
			ArgCount(el.OrLeft, count)
		}
		if el.OrRight != nil {
			ArgCount(el.OrRight, count)
		}
	case *grammar.BooleanTerm:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		}
		if el.AndLeft != nil {
			ArgCount(el.AndLeft, count)
		}
		if el.AndRight != nil {
			ArgCount(el.AndRight, count)
		}
	case grammar.BooleanTerm:
		if el.Unary != nil {
			ArgCount(el.Unary, count)
		}
		if el.AndLeft != nil {
			ArgCount(el.AndLeft, count)
		}
		if el.AndRight != nil {
			ArgCount(el.AndRight, count)
		}
	case *grammar.BooleanFactor:
		ArgCount(el.Test, count)
	case grammar.BooleanFactor:
		ArgCount(el.Test, count)
	case *grammar.BooleanTest:
		ArgCount(el.Primary, count)
	case grammar.BooleanTest:
		ArgCount(el.Primary, count)
	case grammar.BooleanPrimary:
		if el.Predicate != nil {
			ArgCount(el.Predicate, count)
		} else if el.BooleanPredicand != nil {
			ArgCount(el.BooleanPredicand, count)
		}
	case *grammar.BooleanPrimary:
		if el.Predicate != nil {
			ArgCount(el.Predicate, count)
		} else if el.BooleanPredicand != nil {
			ArgCount(el.BooleanPredicand, count)
		}
	case *grammar.CommonValueExpression:
		if el.Numeric != nil {
			ArgCount(el.Numeric, count)
		} else if el.String != nil {
			ArgCount(el.String, count)
		} else if el.Datetime != nil {
			ArgCount(el.Datetime, count)
		} else if el.Interval != nil {
			ArgCount(el.Interval, count)
		}
	case *grammar.NonParenthesizedValueExpressionPrimary:
		if el.UnsignedValue != nil {
			ArgCount(el.UnsignedValue, count)
		} else if el.ColumnReference != nil {
			ArgCount(el.ColumnReference, count)
		} else if el.SetFunction != nil {
			ArgCount(el.SetFunction, count)
		}
	case *grammar.UnsignedValueSpecification:
		if el.UnsignedLiteral != nil {
			ArgCount(el.UnsignedLiteral, count)
		} else if el.GeneralValue != nil {
			ArgCount(el.GeneralValue, count)
		}
	case *grammar.ValueSpecification:
		if el.Literal != nil {
			ArgCount(el.Literal, count)
		} else if el.UnsignedValue != nil {
			ArgCount(el.UnsignedValue, count)
		}
	case *grammar.StringValueExpression:
		if el.Character != nil {
			ArgCount(el.Character, count)
		} else if el.Blob != nil {
			ArgCount(el.Blob, count)
		}
	case *grammar.CharacterValueExpression:
		if el.Factor != nil {
			ArgCount(&el.Factor.Primary, count)
		}
	case *grammar.CharacterPrimary:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Function != nil {
			ArgCount(el.Function, count)
		}
	case *grammar.StringValueFunction:
		if el.Character != nil {
			ArgCount(el.Character, count)
		} else if el.Blob != nil {
			ArgCount(el.Blob, count)
		}
	case *grammar.CharacterValueFunction:
		if el.Substring != nil {
			ss := el.Substring
			ArgCount(&ss.Subject, count)
			ArgCount(&ss.From, count)
			if ss.For != nil {
				ArgCount(ss.For, count)
			}
		} else if el.RegexSubstring != nil {
			ss := el.RegexSubstring
			ArgCount(&ss.Subject, count)
			ArgCount(&ss.Similar, count)
			ArgCount(&ss.Escape, count)
		} else if el.Fold != nil {
			ArgCount(&el.Fold.Subject, count)
		} else if el.Transcoding != nil {
			ArgCount(&el.Transcoding.Subject, count)
		} else if el.Transliteration != nil {
			ArgCount(&el.Transliteration.Subject, count)
		} else if el.Trim != nil {
			ArgCount(&el.Trim.Subject, count)
			if el.Trim.Character != nil {
				ArgCount(el.Trim.Character, count)
			}
		}
	case *grammar.QueryExpression:
		ArgCount(&el.Body, count)
	case grammar.QueryExpression:
		ArgCount(&el.Body, count)
	case *grammar.QueryExpressionBody:
		if el.NonJoinQueryExpression != nil {
			ArgCount(el.NonJoinQueryExpression, count)
		} else if el.JoinedTable != nil {
			ArgCount(el.JoinedTable, count)
		}
	case *grammar.NonJoinQueryExpression:
		if el.NonJoinQueryTerm != nil {
			ArgCount(el.NonJoinQueryTerm, count)
		}
	case *grammar.NonJoinQueryTerm:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Intersect != nil {
			ArgCount(el.Intersect, count)
		}
	case *grammar.NonJoinQueryPrimary:
		if el.SimpleTable != nil {
			ArgCount(el.SimpleTable, count)
		} else if el.ParenthesizedNonJoinQueryExpression != nil {
			ArgCount(el.ParenthesizedNonJoinQueryExpression, count)
		}
	case *grammar.SimpleTable:
		if el.QuerySpecification != nil {
			ArgCount(el.QuerySpecification, count)
		}
	case *grammar.TableReference:
		if el.Primary != nil {
			ArgCount(el.Primary, count)
		} else if el.Joined != nil {
			ArgCount(el.Joined, count)
		}
	case *grammar.TablePrimary:
		if el.DerivedTable != nil {
			ArgCount(&el.DerivedTable.Subquery.QueryExpression, count)
		}
	case *grammar.JoinedTable:
		if el.Qualified != nil {
			j := el.Qualified
			ArgCount(&j.Left, count)
			ArgCount(&j.Right, count)
			ArgCount(&j.On, count)
		} else if el.Cross != nil {
			j := el.Cross
			ArgCount(&j.Left, count)
			ArgCount(&j.Right, count)
		} else if el.Union != nil {
			j := el.Union
			ArgCount(&j.Left, count)
			ArgCount(&j.Right, count)
		} else if el.Natural != nil {
			j := el.Natural
			ArgCount(&j.Left, count)
			ArgCount(&j.Right, count)
		}
	case *grammar.DerivedColumn:
		ArgCount(&el.ValueExpression, count)
	case *grammar.TableExpression:
		ArgCount(&el.FromClause, count)
		if el.WhereClause != nil {
			ArgCount(el.WhereClause, count)
		}
		if el.GroupByClause != nil {
			ArgCount(el.GroupByClause, count)
		}
		if el.HavingClause != nil {
			ArgCount(el.HavingClause, count)
		}
	case *grammar.RowValuePredicand:
		if el.CommonValueExpression != nil {
			ArgCount(el.CommonValueExpression, count)
		} else if el.NonParenthesizedValueExpressionPrimary != nil {
			ArgCount(el.NonParenthesizedValueExpressionPrimary, count)
		} else if el.BooleanPredicand != nil {
			ArgCount(el.BooleanPredicand, count)
		}
	case grammar.RowValuePredicand:
		if el.CommonValueExpression != nil {
			ArgCount(el.CommonValueExpression, count)
		} else if el.NonParenthesizedValueExpressionPrimary != nil {
			ArgCount(el.NonParenthesizedValueExpressionPrimary, count)
		} else if el.BooleanPredicand != nil {
			ArgCount(el.BooleanPredicand, count)
		}
	case *grammar.CursorSpecification:
		ArgCount(&el.QueryExpression, count)
		if el.OrderByClause != nil {
			ArgCount(el.OrderByClause, count)
		}
		if el.LimitClause != nil {
			ArgCount(el.LimitClause, count)
		}
	case *grammar.QuerySpecification:
		ArgCount(&el.SelectList, count)
		ArgCount(&el.TableExpression, count)
	case *grammar.UpdateStatementSearched:
		*count += len(el.Values)
		if el.WhereClause != nil {
			ArgCount(el.WhereClause, count)
		}
	case *grammar.DeleteStatementSearched:
		if el.WhereClause != nil {
			ArgCount(el.WhereClause, count)
		}
	case *grammar.InsertStatement:
		*count += len(el.Values)
	case *grammar.FromClause:
		for _, tref := range el.TableReferences {
			ArgCount(&tref, count)
		}
	case *grammar.LimitClause:
		if el.Offset != nil {
			*count = *count + 2
		} else {
			*count++
		}
	case *grammar.WhereClause:
		ArgCount(&el.SearchCondition, count)
	case *grammar.HavingClause:
		ArgCount(&el.SearchCondition, count)
	case *grammar.Predicate:
		if el.Comparison != nil {
			ArgCount(el.Comparison, count)
		} else if el.In != nil {
			ArgCount(el.In, count)
		} else if el.Between != nil {
			ArgCount(el.Between, count)
		} else if el.Null != nil {
			ArgCount(el.Null, count)
		}
	case *grammar.ComparisonPredicate:
		ArgCount(el.A, count)
		ArgCount(el.B, count)
	case *grammar.BetweenPredicate:
		ArgCount(el.Target, count)
		ArgCount(el.Start, count)
		ArgCount(el.End, count)
	case *grammar.InPredicate:
		ArgCount(el.Target, count)
		for _, rve := range el.Values {
			ArgCount(rve, count)
		}
	case *grammar.NullPredicate:
		ArgCount(el.Target, count)
	case *grammar.Literal, grammar.Literal, *grammar.UnsignedLiteral, grammar.UnsignedLiteral, string, []byte, rune, bool, float64, float32, int8, int16, int64, uint, uint8, uint16, uint32, uint64:
		// a literal common value expression, which can be a string, a number,
		// a null value, etc. The value of the common value expression is
		// contained in a query argument.
		*count++
	default:
		//fmt.Printf("ArgCount on %T: %v\n", el, el)
		// the default is no argument
	}
}
