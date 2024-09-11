//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <data type>    ::=
//          <predefined type>
//      |     <row type>
//      |     <path-resolved user-defined type name>
//      |     <reference type>
//      |     <collection type>
//
// <predefined type>    ::=
//          <character string type> [ CHARACTER SET <character set specification> ] [ <collate clause> ]
//      |     <national character string type> [ <collate clause> ]
//      |     <binary large object string type>
//      |     <numeric type>
//      |     <boolean type>
//      |     <datetime type>
//      |     <interval type>
//
// <character string type>    ::=
//          CHARACTER [ <left paren> <length> <right paren> ]
//      |     CHAR [ <left paren> <length> <right paren> ]
//      |     CHARACTER VARYING <left paren> <length> <right paren>
//      |     CHAR VARYING <left paren> <length> <right paren>
//      |     VARCHAR <left paren> <length> <right paren>
//      |     CHARACTER LARGE OBJECT [ <left paren> <large object length> <right paren> ]
//      |     CHAR LARGE OBJECT [ <left paren> <large object length> <right paren> ]
//      |     CLOB [ <left paren> <large object length> <right paren> ]
//
// <national character string type>    ::=
//          NATIONAL CHARACTER [ <left paren> <length> <right paren> ]
//      |     NATIONAL CHAR [ <left paren> <length> <right paren> ]
//      |     NCHAR [ <left paren> <length> <right paren> ]
//      |     NATIONAL CHARACTER VARYING <left paren> <length> <right paren>
//      |     NATIONAL CHAR VARYING <left paren> <length> <right paren>
//      |     NCHAR VARYING <left paren> <length> <right paren>
//      |     NATIONAL CHARACTER LARGE OBJECT [ <left paren> <large object length> <right paren> ]
//      |     NCHAR LARGE OBJECT [ <left paren> <large object length> <right paren> ]
//      |     NCLOB [ <left paren> <large object length> <right paren> ]
//
// <binary large object string type>    ::=
//          BINARY LARGE OBJECT [ <left paren> <large object length> <right paren> ]
//      |     BLOB [ <left paren> <large object length> <right paren> ]
//
// <numeric type>    ::=   <exact numeric type> | <approximate numeric type>
//
// <exact numeric type>    ::=
//          NUMERIC [ <left paren> <precision> [ <comma> <scale> ] <right paren> ]
//      |     DECIMAL [ <left paren> <precision> [ <comma> <scale> ] <right paren> ]
//      |     DEC [ <left paren> <precision> [ <comma> <scale> ] <right paren> ]
//      |     SMALLINT
//      |     INTEGER
//      |     INT
//      |     BIGINT
//
// <approximate numeric type>    ::=
//          FLOAT [ <left paren> <precision> <right paren> ]
//      |     REAL
//      |     DOUBLE PRECISION
//
// <length>    ::=   <unsigned integer>
//
// <large object length>    ::=
//          <unsigned integer> [ <multiplier> ] [ <char length units> ]
//      |     <large object length token> [ <char length units> ]
//
// <char length units>    ::=   CHARACTERS | CODE_UNITS | OCTETS
//
// <precision>    ::=   <unsigned integer>
//
// <scale>    ::=   <unsigned integer>
//
// <boolean type>    ::=   BOOLEAN
//
// <datetime type>    ::=
//          DATE
//      |     TIME [ <left paren> <time precision> <right paren> ] [ <with or without time zone> ]
//      |     TIMESTAMP [ <left paren> <timestamp precision> <right paren> ] [ <with or without time zone> ]
//
// <with or without time zone>    ::=   WITH TIME ZONE | WITHOUT TIME ZONE
//
// <time precision>    ::=   <time fractional seconds precision>
//
// <timestamp precision>    ::=   <time fractional seconds precision>
//
// <time fractional seconds precision>    ::=   <unsigned integer>
//
// <interval type>    ::=   INTERVAL <interval qualifier>
//
// <row type>    ::=   ROW <row type body>
//
// <row type body>    ::=   <left paren> <field definition> [ { <comma> <field definition> }... ] <right paren>
//
// <reference type>    ::=   REF <left paren> <referenced type> <right paren> [ <scope clause> ]
//
// <scope clause>    ::=   SCOPE <table name>
//
// <referenced type>    ::=   <path-resolved user-defined type name>
//
// <path-resolved user-defined type name>    ::=   <user-defined type name>
//
// <collection type>    ::=   <array type> | <multiset type>
//
// <array type>    ::=   <data type> ARRAY [ <left bracket or trigraph> <unsigned integer> <right bracket or trigraph> ]
//
// <multiset type>    ::=   <data type> MULTISET

type CharacterLengthUnits int

const (
	CharacterLengthUnitsCharacters CharacterLengthUnits = iota
	CharacterLengthUnitsCodeUnits
	CharacterLengthUnitsOctets
)

var CharacterLengthUnitsSymbol = map[CharacterLengthUnits]string{
	CharacterLengthUnitsCharacters: "CHARACTERS",
	CharacterLengthUnitsCodeUnits:  "CODE_UNITS",
	CharacterLengthUnitsOctets:     "OCTETS",
}
