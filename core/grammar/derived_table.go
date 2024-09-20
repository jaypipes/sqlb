//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <derived table>    ::=   <table subquery>
//
// <table subquery>    ::=   <subquery>

type DerivedTable struct {
	Subquery
}
