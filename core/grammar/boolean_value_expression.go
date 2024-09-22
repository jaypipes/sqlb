//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <boolean value expression>    ::=
//          <boolean term>
//      |     <boolean value expression> OR <boolean term>
//
// <boolean term>    ::=
//          <boolean factor>
//      |     <boolean term> AND <boolean factor>
//
// <boolean factor>    ::=   [ NOT ] <boolean test>
//
// <boolean test>    ::=   <boolean primary> [ IS [ NOT ] <truth value> ]
//
// <truth value>    ::=   TRUE | FALSE | UNKNOWN
//
// <boolean primary>    ::=   <predicate> | <boolean predicand>
//
// <boolean predicand>    ::=
//          <parenthesized boolean value expression>
//      |     <nonparenthesized value expression primary>
//
// <parenthesized boolean value expression>    ::=   <left paren> <boolean value expression> <right paren>

// BooleanValueExpression represents a boolean comparison expression in the SQL
// statement, e.g. "a = b"
type BooleanValueExpression struct {
	Unary   *BooleanTerm
	OrLeft  *BooleanValueExpression
	OrRight *BooleanTerm
}

func (e *BooleanValueExpression) ArgCount(count *int) {
	if e.Unary != nil {
		e.Unary.ArgCount(count)
	}
	if e.OrLeft != nil {
		e.OrLeft.ArgCount(count)
	}
	if e.OrRight != nil {
		e.OrRight.ArgCount(count)
	}
}

type BooleanTerm struct {
	Unary    *BooleanFactor
	AndLeft  *BooleanTerm
	AndRight *BooleanFactor
}

func (t *BooleanTerm) ArgCount(count *int) {
	if t.Unary != nil {
		t.Unary.ArgCount(count)
	}
	if t.AndLeft != nil {
		t.AndLeft.ArgCount(count)
	}
	if t.AndRight != nil {
		t.AndRight.ArgCount(count)
	}
}

type BooleanFactor struct {
	Test BooleanTest
	Not  bool
}

func (f *BooleanFactor) ArgCount(count *int) {
	f.Test.ArgCount(count)
}

type BooleanTest struct {
	Primary BooleanPrimary
}

func (t *BooleanTest) ArgCount(count *int) {
	t.Primary.ArgCount(count)
}

type BooleanPrimary struct {
	Predicate *Predicate
	Predicand *BooleanPredicand
}

func (p *BooleanPrimary) ArgCount(count *int) {
	if p.Predicate != nil {
		p.Predicate.ArgCount(count)
	} else if p.Predicand != nil {
		p.Predicand.ArgCount(count)
	}
}

type BooleanPredicand struct {
	Parenthesized *BooleanValueExpression
	Primary       *NonParenthesizedValueExpressionPrimary
}

func (p *BooleanPredicand) ArgCount(count *int) {
	if p.Parenthesized != nil {
		p.Parenthesized.ArgCount(count)
	} else if p.Primary != nil {
		p.Primary.ArgCount(count)
	}
}
