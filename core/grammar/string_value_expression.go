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

func (e *StringValueExpression) ArgCount(count *int) {
	if e.Character != nil {
		e.Character.ArgCount(count)
	} else if e.Blob != nil {
		e.Blob.ArgCount(count)
	}
}

type CharacterValueExpression struct {
	// Concatenation *Concatenation
	Factor *CharacterFactor
}

func (e *CharacterValueExpression) ArgCount(count *int) {
	e.Factor.ArgCount(count)
}

type CharacterFactor struct {
	Primary   CharacterPrimary
	Collation *string
}

func (f *CharacterFactor) ArgCount(count *int) {
	f.Primary.ArgCount(count)
}

type CharacterPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *StringValueFunction
}

func (p *CharacterPrimary) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	} else if p.Function != nil {
		p.Function.ArgCount(count)
	}
}

type BlobValueExpression struct {
	// BlobConcatenation
	Factor *BlobFactor
}

func (e *BlobValueExpression) ArgCount(count *int) {
	e.Factor.ArgCount(count)
}

type BlobFactor struct {
	Primary BlobPrimary
}

func (f *BlobFactor) ArgCount(count *int) {
	f.Primary.ArgCount(count)
}

type BlobPrimary struct {
	Primary  *ValueExpressionPrimary
	Function *StringValueFunction
}

func (p *BlobPrimary) ArgCount(count *int) {
	if p.Primary != nil {
		p.Primary.ArgCount(count)
	} else if p.Function != nil {
		p.Function.ArgCount(count)
	}
}
