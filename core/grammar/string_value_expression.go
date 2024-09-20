//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <string value expression>    ::=   <character value expression> | <blob value expression>
//
// <character value expression>    ::=   <concatenation> | <character factor>
//
// <concatenation>    ::=   <character value expression> <concatenation operator> <character factor>
//
// <character factor>    ::=   <character primary> [ <collate clause> ]
//
// <character primary>    ::=   <value expression primary> | <string value function>
//
// <blob value expression>    ::=   <blob concatenation> | <blob factor>
//
// <blob factor>    ::=   <blob primary>
//
// <blob primary>    ::=   <value expression primary> | <string value function>
//
// <blob concatenation>    ::=   <blob value expression> <concatenation operator> <blob factor>

type StringValueExpression struct {
	Character *CharacterValueExpression
	Blob      *BlobValueExpression
}

type CharacterValueExpression struct {
	// Concatenation *Concatenation
	Factor *CharacterFactor
}

type CharacterFactor struct {
	Primary   CharacterPrimary
	Collation *string
}

type CharacterPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *StringValueFunction
}

type BlobValueExpression struct {
	// BlobConcatenation
	Factor *BlobFactor
}

type BlobFactor struct {
	Primary BlobPrimary
}

type BlobPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *StringValueFunction
}