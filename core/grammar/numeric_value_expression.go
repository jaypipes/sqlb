//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <numeric value expression>    ::=
//          <term>
//      |     <numeric value expression> <plus sign> <term>
//      |     <numeric value expression> <minus sign> <term>
//
// <term>    ::=
//          <factor>
//      |     <term> <asterisk> <factor>
//      |     <term> <solidus> <factor>
//
// <factor>    ::=   [ <sign> ] <numeric primary>
//
// <numeric primary>    ::=
//          <value expression primary>
//      |     <numeric value function>

type Sign int

const (
	SignPlus Sign = iota
	SignMinus
)

var SignSymbol = map[Sign]string{
	SignPlus:  "+",
	SignMinus: "-",
}

type NumericOperation int

const (
	NumericOperationAdd NumericOperation = iota
	NumericOperationSubtract
	NumericOperationMultiply
	NumericOperationDivide
)

var NumericOperationSymbol = map[NumericOperation]string{
	NumericOperationAdd:      " + ",
	NumericOperationSubtract: " - ",
	NumericOperationMultiply: " * ",
	NumericOperationDivide:   " / ",
}

type NumericValueExpression struct {
	Unary       *Term
	AddSubtract *AddSubtractExpression
}

func (e *NumericValueExpression) ArgCount(count *int) {
	if e.Unary != nil {
		e.Unary.ArgCount(count)
	} else if e.AddSubtract != nil {
		e.AddSubtract.Left.ArgCount(count)
		e.AddSubtract.Right.ArgCount(count)
	}
}

type AddSubtractExpression struct {
	Left     NumericValueExpression
	Right    Term
	Subtract bool
}

func (e *AddSubtractExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type Term struct {
	Unary          *Factor
	MultiplyDivide *MultiplyDivideExpression
}

func (t *Term) ArgCount(count *int) {
	if t.Unary != nil {
		t.Unary.ArgCount(count)
	} else if t.MultiplyDivide != nil {
		t.MultiplyDivide.Left.ArgCount(count)
		t.MultiplyDivide.Right.ArgCount(count)
	}
}

type MultiplyDivideExpression struct {
	Left   Term
	Right  Factor
	Divide bool
}

func (e *MultiplyDivideExpression) ArgCount(count *int) {
	e.Left.ArgCount(count)
	e.Right.ArgCount(count)
}

type Factor struct {
	Sign    Sign
	Primary NumericPrimary
}

func (f *Factor) ArgCount(count *int) {
	f.Primary.ArgCount(count)
}

type NumericPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *NumericValueFunction
}

func (p *NumericPrimary) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	} else if p.Function != nil {
		p.Function.ArgCount(count)
	}
}
