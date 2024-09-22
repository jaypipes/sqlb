//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <value expression primary>    ::=
//          <parenthesized value expression>
//      |   <nonparenthesized value expression primary>
//
// <parenthesized value expression>    ::=   <left paren> <value expression> <right paren>
//
// <nonparenthesized value expression primary>    ::=
//          <unsigned value specification>
//      |     <column reference>
//      |     <set function specification>
//      |     <window function>
//      |     <scalar subquery>
//      |     <case expression>
//      |     <cast specification>
//      |     <field reference>
//      |     <subtype treatment>
//      |     <method invocation>
//      |     <static method invocation>
//      |     <new specification>
//      |     <attribute or method reference>
//      |     <reference resolution>
//      |     <collection value constructor>
//      |     <array element reference>
//      |     <multiset element reference>
//      |     <routine invocation>
//      |     <next value expression>

type ValueExpressionPrimary struct {
	Parenthesized *ValueExpression
	Primary       *NonParenthesizedValueExpressionPrimary
}

func (p *ValueExpressionPrimary) ArgCount(count *int) {
	if p.Parenthesized != nil {
		p.Parenthesized.ArgCount(count)
	} else if p.Primary != nil {
		p.Primary.ArgCount(count)
	}
}

type NonParenthesizedValueExpressionPrimary struct {
	UnsignedValue   *UnsignedValueSpecification
	ColumnReference *ColumnReference
	SetFunction     *SetFunctionSpecification
	//WindowFunction *WindowFunction
	ScalarSubquery *Subquery
	//CaseExpression *CaseExpression
	//CastSpecification *CastSpecification
	//FieldReference *FieldReference
	//SubtypeTreatment *SubtypeTreatment
	//MethodInvocation *MethodInvocation
	//StaticMethodInvocation *StaticMethodInvocation
	//NewSpecification *NewSpecification
	//AttributeOrMethodReference *AttributeOrMethodReference
	//ReferenceResolution *ReferenceResolution
	//CollectionValueConstructor *CollectionValueConstructor
	//ArrayElementReference *ArrayElementReference
}

func (p *NonParenthesizedValueExpressionPrimary) ArgCount(count *int) {
	if p.UnsignedValue != nil {
		p.UnsignedValue.ArgCount(count)
	} else if p.ColumnReference != nil {
		p.ColumnReference.ArgCount(count)
	} else if p.SetFunction != nil {
		p.SetFunction.ArgCount(count)
	} else if p.ScalarSubquery != nil {
		p.ScalarSubquery.ArgCount(count)
	}
}
