//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <identifier chain>    ::=   <identifier> [ { <period> <identifier> }... ]
//
// <basic identifier chain>    ::=   <identifier chain>

type IdentifierChain struct {
	Identifiers []string
}

// <identifier>    ::=   <actual identifier>
//
// <actual identifier>    ::=   <regular identifier> | <delimited identifier>
//
// <SQL language identifier>    ::=
//          <SQL language identifier start> [ { <underscore> | <SQL language identifier part> }... ]
//
// <SQL language identifier start>    ::=   <simple Latin letter>
//
// <SQL language identifier part>    ::=   <simple Latin letter> | <digit>
//
// <authorization identifier>    ::=   <role name> | <user identifier>
//
// <table name>    ::=   <local or schema qualified name>
//
// <domain name>    ::=   <schema qualified name>
//
// <schema name>    ::=   [ <catalog name> <period> ] <unqualified schema name>
//
// <catalog name>    ::=   <identifier>
//
// <schema qualified name>    ::=   [ <schema name> <period> ] <qualified identifier>
//
// <local or schema qualified name>    ::=   [ <local or schema qualifier> <period> ] <qualified identifier>
//
// <local or schema qualifier>    ::=   <schema name> | MODULE
//
// <qualified identifier>    ::=   <identifier>
//
// <column name>    ::=   <identifier>
//
// <correlation name>    ::=   <identifier>
//
// <query name>    ::=   <identifier>
//
// <SQL-client module name>    ::=   <identifier>
//
// <procedure name>    ::=   <identifier>
//
// <schema qualified routine name>    ::=   <schema qualified name>
//
// <method name>    ::=   <identifier>
//
// <specific name>    ::=   <schema qualified name>
//
// <cursor name>    ::=   <local qualified name>
//
// <local qualified name>    ::=   [ <local qualifier> <period> ] <qualified identifier>
//
// <local qualifier>    ::=   MODULE
//
// <host parameter name>    ::=   <colon> <identifier>
//
// <SQL parameter name>    ::=   <identifier>
//
// <constraint name>    ::=   <schema qualified name>
//
// <external routine name>    ::=   <identifier> | <character string literal>
//
// <trigger name>    ::=   <schema qualified name>
//
// <collation name>    ::=   <schema qualified name>
//
// <character set name>    ::=   [ <schema name> <period> ] <SQL language identifier>
//
// <transliteration name>    ::=   <schema qualified name>
//
// <transcoding name>    ::=   <schema qualified name>
//
// <user-defined type name>    ::=   <schema qualified type name>
//
// <schema-resolved user-defined type name>    ::=   <user-defined type name>
//
// <schema qualified type name>    ::=   [ <schema name> <period> ] <qualified identifier>
//
// <attribute name>    ::=   <identifier>
//
// <field name>    ::=   <identifier>
//
// <savepoint name>    ::=   <identifier>
//
// <sequence generator name>    ::=   <schema qualified name>
//
// <role name>    ::=   <identifier>
//
// <user identifier>    ::=   <identifier>
//
// <connection name>    ::=   <simple value specification>
//
// <SQL-server name>    ::=   <simple value specification>
//
// <connection user name>    ::=   <simple value specification>
//
// <SQL statement name>    ::=   <statement name> | <extended statement name>
//
// <statement name>    ::=   <identifier>
//
// <extended statement name>    ::=   [ <scope option> ] <simple value specification>
//
// <dynamic cursor name>    ::=   <cursor name> | <extended cursor name>
//
// <extended cursor name>    ::=   [ <scope option> ] <simple value specification>
//
// <descriptor name>    ::=   [ <scope option> ] <simple value specification>
//
// <scope option>    ::=   GLOBAL | LOCAL
//
// <window name>    ::=   <identifier>

type SchemaQualifiedName struct {
	SchemaName  *string
	Identifiers IdentifierChain
}
