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
// <fold>    ::=   { UPPER | LOWER } <left paren> <character value expression> <right paren>
//
// <transcoding>    ::=   CONVERT <left paren> <character value expression> USING <transcoding name> <right paren>
//
// <character transliteration>    ::=   TRANSLATE <left paren> <character value expression> USING <transliteration name> <right paren>
//
// <trim function>    ::=   TRIM <left paren> <trim operands> <right paren>
//
// <trim operands>    ::=   [ [ <trim specification> ] [ <trim character> ] FROM ] <trim source>
//
// <trim source>    ::=   <character value expression>
//
// <trim specification>    ::=   LEADING | TRAILING | BOTH
//
// <trim character>    ::=   <character value expression>
//
// <character overlay function>    ::=
//          OVERLAY <left paren> <character value expression> PLACING <character value expression>
//          FROM <start position> [ FOR <string length> ] [ USING <char length units> ] <right paren>
//
// <normalize function>    ::=   NORMALIZE <left paren> <character value expression> <right paren>
//
// <specific type method>    ::=   <user-defined type value expression> <period> SPECIFICTYPE
//
// <blob value function>    ::=
//          <blob substring function>
//      |     <blob trim function>
//      |     <blob overlay function>
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

type CharacterValueFunction struct {
	Substring       *SubstringFunction
	RegexSubstring  *RegexSubstringFunction
	Fold            *FoldFunction
	Transcoding     *TranscodingFunction
	Transliteration *TransliterationFunction
	Trim            *TrimFunction
	Overlay         *OverlayFunction
	Normalize       *NormalizeFunction
	SpecificType    *SpecificTypeFunction
}

// <character substring function>    ::=
//          SUBSTRING <left paren> <character value expression> FROM <start position>
//          [ FOR <string length> ] [ USING <char length units> ] <right paren>

type SubstringFunction struct {
	Subject CharacterValueExpression
	From    NumericValueExpression
	For     *NumericValueExpression
	Using   CharacterLengthUnits
}

// <regular expression substring function>    ::=
//          SUBSTRING <left paren> <character value expression>
//          SIMILAR <character value expression> ESCAPE <escape character> <right paren>

type RegexSubstringFunction struct {
	Subject CharacterValueExpression
	Similar CharacterValueExpression
	Escape  CharacterValueExpression
}

type FoldFunction struct{}

type TranscodingFunction struct{}

type TransliterationFunction struct{}

type TrimFunction struct{}

type OverlayFunction struct{}

type NormalizeFunction struct{}

type SpecificTypeFunction struct{}

type BlobValueFunction struct{}
