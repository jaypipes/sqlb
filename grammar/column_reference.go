//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <column reference>    ::=
//          <basic identifier chain>
//      |   MODULE <period> <qualified identifier> <period> <column name>

type ColumnReference struct {
	BasicIdentifierChain *IdentifierChain
	//ModuleIdentifier *ModuleIdentifier
	Correlation *Correlation
}
