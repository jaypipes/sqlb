//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <value specification>    ::=   <literal> | <general value specification>
//
// <unsigned value specification>    ::=   <unsigned literal> | <general value specification>
//
// <general value specification>    ::=
//          <host parameter specification>
//      |     <SQL parameter reference>
//      |     <dynamic parameter specification>
//      |     <embedded variable specification>
//      |     <current collation specification>
//      |     CURRENT_DEFAULT_TRANSFORM_GROUP
//      |     CURRENT_PATH
//      |     CURRENT_ROLE
//      |     CURRENT_TRANSFORM_GROUP_FOR_TYPE <path-resolved user-defined type name>
//      |     CURRENT_USER
//      |     SESSION_USER
//      |     SYSTEM_USER
//      |     USER
//      |     VALUE
//
// <simple value specification>    ::=
//          <literal>
//      |     <host parameter name>
//      |     <SQL parameter reference>
//      |     <embedded variable name>
//
// <target specification>    ::=
//          <host parameter specification>
//      |     <SQL parameter reference>
//      |     <column reference>
//      |     <target array element specification>
//      |     <dynamic parameter specification>
//      |     <embedded variable specification>
//
// <simple target specification>    ::=
//          <host parameter specification>
//      |     <SQL parameter reference>
//      |     <column reference>
//      |     <embedded variable name>
//
// <host parameter specification>    ::=   <host parameter name> [ <indicator parameter> ]
//
// <dynamic parameter specification>    ::=   <question mark>
//
// <embedded variable specification>    ::=   <embedded variable name> [ <indicator variable> ]
//
// <indicator variable>    ::=   [ INDICATOR ] <embedded variable name>
//
// <indicator parameter>    ::=   [ INDICATOR ] <host parameter name>
//
// <target array element specification>    ::=
//          <target array reference> <left bracket or trigraph> <simple value specification> <right bracket or trigraph>
//
// <target array reference>    ::=   <SQL parameter reference> | <column reference>
//
// <current collation specification>    ::=   CURRENT_COLLATION <left paren> <string value expression> <right paren>

type ValueSpecification struct {
	Literal                    *Literal
	UnsignedValueSpecification *UnsignedValueSpecification
}

type UnsignedValueSpecification struct {
	UnsignedLiteral           *UnsignedLiteral
	GeneralValueSpecification *GeneralValueSpecification
}

type GeneralValueSpecification struct {
	//HostParameterSpecification *HostParameterSpecification
	//SQLParameterReference *SQLParameterReference
	//DynamicParameterSpecification *DynamicParameterSpecification
	//EmbeddedVariableSpecification *EmbeddedVariableSpecification
	//CurrentCollationSpecification *CurrentCollationSpecification
}
