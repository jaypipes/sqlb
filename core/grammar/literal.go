//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <literal>    ::=   <signed numeric literal> | <general literal>
//
// <unsigned literal>    ::=   <unsigned numeric literal> | <general literal>
//
// <general literal>    ::=
//          <character string literal>
//      |     <national character string literal>
//      |     <Unicode character string literal>
//      |     <binary string literal>
//      |     <datetime literal>
//      |     <interval literal>
//      |     <boolean literal>
//
// <character string literal>    ::=
//          [ <introducer> <character set specification> ]
//          <quote> [ <character representation> ... ] <quote>
//          [ { <separator> <quote> [ <character representation> ... ] <quote> }... ]
//
// <introducer>    ::=   <underscore>
//
// <character representation>    ::=   <nonquote character> | <quote symbol>
//
// <nonquote character>    ::=   !! See the Syntax Rules.
//
// The <quote symbol> rule consists of two immediately adjacent <quote> marks with no spaces. As usual, this would be best handled in the lexical analyzer, not in the grammar.
//
// <quote symbol>    ::=   <quote> <quote>
//
// <national character string literal>    ::=
//          N <quote> [ <character representation> ... ] <quote>
//          [ { <separator> <quote> [ <character representation> ... ] <quote> }... ]
//
// <Unicode character string literal>    ::=
//          [ <introducer> <character set specification> ]
//          U <ampersand> <quote> [ <Unicode representation> ... ] <quote>
//          [ { <separator> <quote> [ <Unicode representation> ... ] <quote> }... ]
//          [ ESCAPE <escape character> ]
//
// <Unicode representation>    ::=   <character representation> | <Unicode escape value>
//
// <binary string literal>    ::=
//          X <quote> [ { <hexit> <hexit> }... ] <quote>
//          [ { <separator> <quote> [ { <hexit> <hexit> }... ] <quote> }... ]
//          [ ESCAPE <escape character> ]
//
// <hexit>    ::=   <digit> | A | B | C | D | E | F | a | b | c | d | e | f
//
// <signed numeric literal>    ::=   [ <sign> ] <unsigned numeric literal>
//
// <unsigned numeric literal>    ::=   <exact numeric literal> | <approximate numeric literal>
//
// <exact numeric literal>    ::=
//          <unsigned integer> [ <period> [ <unsigned integer> ] ]
//      |     <period> <unsigned integer>
//
// <sign>    ::=   <plus sign> | <minus sign>
//
// <approximate numeric literal>    ::=   <mantissa> E <exponent>
//
// <mantissa>    ::=   <exact numeric literal>
//
// <exponent>    ::=   <signed integer>
//
// <signed integer>    ::=   [ <sign> ] <unsigned integer>
//
// <unsigned integer>    ::=   <digit> ...
//
// <datetime literal>    ::=   <date literal> | <time literal> | <timestamp literal>
//
// <date literal>    ::=   DATE <date string>
//
// <time literal>    ::=   TIME <time string>
//
// <timestamp literal>    ::=   TIMESTAMP <timestamp string>
//
// <date string>    ::=   <quote> <unquoted date string> <quote>
//
// <time string>    ::=   <quote> <unquoted time string> <quote>
//
// <timestamp string>    ::=   <quote> <unquoted timestamp string> <quote>
//
// <time zone interval>    ::=   <sign> <hours value> <colon> <minutes value>
//
// <date value>    ::=   <years value> <minus sign> <months value> <minus sign> <days value>
//
// <time value>    ::=   <hours value> <colon> <minutes value> <colon> <seconds value>
//
// <interval literal>    ::=   INTERVAL [ <sign> ] <interval string> <interval qualifier>
//
// <interval string>    ::=   <quote> <unquoted interval string> <quote>
//
// <unquoted date string>    ::=   <date value>
//
// <unquoted time string>    ::=   <time value> [ <time zone interval> ]
//
// <unquoted timestamp string>    ::=   <unquoted date string> <space> <unquoted time string>
//
// <unquoted interval string>    ::=   [ <sign> ] { <year-month literal> | <day-time literal> }
//
// <year-month literal>    ::=   <years value> | [ <years value> <minus sign> ] <months value>
//
// <day-time literal>    ::=   <day-time interval> | <time interval>
//
// <day-time interval>    ::=
//          <days value> [ <space> <hours value> [ <colon> <minutes value> [ <colon> <seconds value> ] ] ]
//
// <time interval>    ::=
//          <hours value> [ <colon> <minutes value> [ <colon> <seconds value> ] ]
//      |     <minutes value> [ <colon> <seconds value> ]
//      |     <seconds value>
//
// <years value>    ::=   <datetime value>
//
// <months value>    ::=   <datetime value>
//
// <days value>    ::=   <datetime value>
//
// <hours value>    ::=   <datetime value>
//
// <minutes value>    ::=   <datetime value>
//
// <seconds value>    ::=   <seconds integer value> [ <period> [ <seconds fraction> ] ]
//
// <seconds integer value>    ::=   <unsigned integer>
//
// <seconds fraction>    ::=   <unsigned integer>
//
// <datetime value>    ::=   <unsigned integer>
//
// <boolean literal>    ::=   TRUE | FALSE | UNKNOWN

type Literal struct {
	SignedNumeric *SignedNumericLiteral
	General       *GeneralLiteral
}

func (l *Literal) ArgCount(count *int) {
	if l.SignedNumeric != nil {
		l.SignedNumeric.ArgCount(count)
	} else if l.General != nil {
		l.General.ArgCount(count)
	}
}

type UnsignedLiteral struct {
	UnsignedNumeric *UnsignedNumericLiteral
	General         *GeneralLiteral
}

func (l *UnsignedLiteral) ArgCount(count *int) {
	if l.UnsignedNumeric != nil {
		l.UnsignedNumeric.ArgCount(count)
	} else if l.General != nil {
		l.General.ArgCount(count)
	}
}

type SignedNumericLiteral struct {
	Value interface{}
}

func (l *SignedNumericLiteral) ArgCount(count *int) {
	*count++
}

type UnsignedNumericLiteral struct {
	Value interface{}
}

func (l *UnsignedNumericLiteral) ArgCount(count *int) {
	*count++
}

type GeneralLiteral struct {
	Value interface{}
}

func (l *GeneralLiteral) ArgCount(count *int) {
	*count++
}
