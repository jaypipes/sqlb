//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <table reference>    ::=   <table primary or joined table> [ <sample clause> ]
//
// <table primary or joined table>    ::=   <table primary> | <joined table>
//
// <sample clause>    ::=
//          TABLESAMPLE <sample method> <left paren> <sample percentage> <right paren> [ <repeatable clause> ]
//
// <sample method>    ::=   BERNOULLI | SYSTEM
//
// <repeatable clause>    ::=   REPEATABLE <left paren> <repeat argument> <right paren>
//
// <sample percentage>    ::=   <numeric value expression>
//
// <repeat argument>    ::=   <numeric value expression>

// TableReference represents the <table reference> SQL grammar element
type TableReference struct {
	Primary *TablePrimary
	Joined  *JoinedTable
}

// TableReferenceFromAny evaluates the supplied
// interface argument and returns a *TableReference if
// the supplied argument can be converted into a
// TableReference, or nil if the conversion cannot be
// done.
func TableReferenceFromAny(
	subject interface{},
) *TableReference {
	switch v := subject.(type) {
	case *TableReference:
		return v
	case TableReference:
		return &v
	case *TablePrimary:
		return &TableReference{Primary: v}
	case TablePrimary:
		return &TableReference{Primary: &v}
	case *JoinedTable:
		return &TableReference{Joined: v}
	case JoinedTable:
		return &TableReference{Joined: &v}
	}
	return nil
}
