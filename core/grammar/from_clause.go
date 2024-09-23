//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <from clause>    ::=   FROM <table reference list>
//
// <table reference list>    ::=   <table reference> [ { <comma> <table reference> }... ]

type FromClause struct {
	TableReferences []TableReference
}

func (c *FromClause) ArgCount(count *int) {
	for _, tref := range c.TableReferences {
		tref.ArgCount(count)
	}
}
