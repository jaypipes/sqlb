//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <string value function>    ::=   <character value function> | <blob value function>
//
// <character value function>    ::=
//          <character substring function>
//      |     <regular expression substring function>
//      |     <fold>
//      |     <transcoding>
//      |     <character transliteration>
//      |     <trim function>
//      |     <character overlay function>
//      |     <normalize function>
//      |     <specific type method>
//
// <blob substring function>    ::=
//          SUBSTRING <left paren> <blob value expression> FROM <start position> [ FOR <string length> ] <right paren>
//
// <blob trim function>    ::=   TRIM <left paren> <blob trim operands> <right paren>
//
// <blob trim operands>    ::=   [ [ <trim specification> ] [ <trim octet> ] FROM ] <blob trim source>
//
// <blob trim source>    ::=   <blob value expression>
//
// <trim octet>    ::=   <blob value expression>
//
// <blob overlay function>    ::=
//          OVERLAY <left paren> <blob value expression> PLACING <blob value expression>
//          FROM <start position> [ FOR <string length> ] <right paren>
//
// <start position>    ::=   <numeric value expression>
//
// <string length>    ::=   <numeric value expression>

type StringValueFunction struct {
	Character *CharacterValueFunction
	Blob      *BlobValueFunction
}

func (f *StringValueFunction) ArgCount(count *int) {
	if f.Character != nil {
		f.Character.ArgCount(count)
	} else if f.Blob != nil {
		f.Blob.ArgCount(count)
	}
}

type CharacterValueFunction struct {
	Substring       *CharacterSubstringFunction
	RegexSubstring  *RegexSubstringFunction
	Fold            *FoldFunction
	Transcoding     *TranscodingFunction
	Transliteration *CharacterTransliterationFunction
	Trim            *TrimFunction
	Overlay         *CharacterOverlayFunction
	Normalize       *NormalizeFunction
	SpecificType    *SpecificTypeFunction
}

func (f *CharacterValueFunction) ArgCount(count *int) {
	if f.Substring != nil {
		f.Substring.ArgCount(count)
	} else if f.RegexSubstring != nil {
		f.RegexSubstring.ArgCount(count)
	} else if f.Fold != nil {
		f.Fold.ArgCount(count)
	} else if f.Transcoding != nil {
		f.Transcoding.ArgCount(count)
	} else if f.Transliteration != nil {
		f.Transliteration.ArgCount(count)
	} else if f.Trim != nil {
		f.Trim.ArgCount(count)
	} else if f.Overlay != nil {
		f.Overlay.ArgCount(count)
	} else if f.Normalize != nil {
		f.Normalize.ArgCount(count)
	}
}

// <character substring function>    ::=
//          SUBSTRING <left paren> <character value expression> FROM <start position>
//          [ FOR <string length> ] [ USING <char length units> ] <right paren>

type CharacterSubstringFunction struct {
	Subject CharacterValueExpression
	From    NumericValueExpression
	For     *NumericValueExpression
	Using   CharacterLengthUnits
}

func (f *CharacterSubstringFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
	f.From.ArgCount(count)
	if f.For != nil {
		f.For.ArgCount(count)
	}
}

// <regular expression substring function>    ::=
//          SUBSTRING <left paren> <character value expression>
//          SIMILAR <character value expression> ESCAPE <escape character> <right paren>

type RegexSubstringFunction struct {
	Subject CharacterValueExpression
	Similar CharacterValueExpression
	Escape  CharacterValueExpression
}

func (f *RegexSubstringFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
	f.Similar.ArgCount(count)
	f.Escape.ArgCount(count)
}

// <fold>    ::=   { UPPER | LOWER } <left paren> <character value expression> <right paren>

type FoldCase int

const (
	FoldCaseUpper FoldCase = iota
	FoldCaseLower
)

var FoldCaseSymbols = map[FoldCase]string{
	FoldCaseUpper: "UPPER",
	FoldCaseLower: "LOWER",
}

type FoldFunction struct {
	Case    FoldCase
	Subject CharacterValueExpression
}

func (f *FoldFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
}

// <transcoding>    ::=   CONVERT <left paren> <character value expression> USING <transcoding name> <right paren>

type TranscodingFunction struct {
	Subject CharacterValueExpression
	Using   SchemaQualifiedName
}

func (f *TranscodingFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
}

// <character transliteration>    ::=   TRANSLATE <left paren> <character value expression> USING <transliteration name> <right paren>

type CharacterTransliterationFunction struct {
	Subject CharacterValueExpression
	Using   SchemaQualifiedName
}

func (f *CharacterTransliterationFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
}

// <trim function>    ::=   TRIM <left paren> <trim operands> <right paren>
//
// <trim operands>    ::=   [ [ <trim specification> ] [ <trim character> ] FROM ] <trim source>
//
// <trim source>    ::=   <character value expression>
//
// <trim specification>    ::=   LEADING | TRAILING | BOTH
//
// <trim character>    ::=   <character value expression>

type TrimSpecification int

const (
	TrimSpecificationBoth TrimSpecification = iota
	TrimSpecificationLeading
	TrimSpecificationTrailing
)

var TrimSpecificationSymbols = map[TrimSpecification]string{
	TrimSpecificationBoth:     "BOTH",
	TrimSpecificationLeading:  "LEADING",
	TrimSpecificationTrailing: "TRAILING",
}

type TrimFunction struct {
	Specification TrimSpecification
	Character     *CharacterValueExpression
	Subject       CharacterValueExpression
}

func (f *TrimFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
	if f.Character != nil {
		f.Character.ArgCount(count)
	}
}

// <character overlay function>    ::=
//          OVERLAY <left paren> <character value expression> PLACING <character value expression>
//          FROM <start position> [ FOR <string length> ] [ USING <char length units> ] <right paren>

type CharacterOverlayFunction struct {
	Subject CharacterValueExpression
	Placing CharacterValueExpression
	From    NumericValueExpression
	For     *NumericValueExpression
	Using   CharacterLengthUnits
}

func (f *CharacterOverlayFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
	f.Placing.ArgCount(count)
	f.From.ArgCount(count)
	if f.For != nil {
		f.For.ArgCount(count)
	}
}

// <normalize function>    ::=   NORMALIZE <left paren> <character value expression> <right paren>

type NormalizeFunction struct {
	Subject CharacterValueExpression
}

func (f *NormalizeFunction) ArgCount(count *int) {
	f.Subject.ArgCount(count)
}

// <specific type method>    ::=   <user-defined type value expression> <period> SPECIFICTYPE

type SpecificTypeFunction struct{}

func (f *SpecificTypeFunction) ArgCount(count *int) {
}

// <blob value function>    ::=
//          <blob substring function>
//      |     <blob trim function>
//      |     <blob overlay function>

type BlobValueFunction struct {
}

func (f *BlobValueFunction) ArgCount(count *int) {
}
