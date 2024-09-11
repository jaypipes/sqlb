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

type AddSubtractExpression struct {
	Left     NumericValueExpression
	Right    Term
	Subtract bool
}

type Term struct {
	Unary          *Factor
	MultiplyDivide *MultiplyDivideExpression
}

type MultiplyDivideExpression struct {
	Left   Term
	Right  Factor
	Divide bool
}

type Factor struct {
	Sign    Sign
	Primary NumericPrimary
}

type NumericPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *NumericValueFunction
}
