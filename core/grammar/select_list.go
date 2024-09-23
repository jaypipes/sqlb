//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

// <select list>    ::=   <asterisk> | <select sublist> [ { <comma> <select sublist> }... ]
//
// <select sublist>    ::=   <derived column> | <qualified asterisk>
//
// <qualified asterisk>    ::=
//          <asterisked identifier chain> <period> <asterisk>
//      |     <all fields reference>
//
// <asterisked identifier chain>    ::=   <asterisked identifier> [ { <period> <asterisked identifier> }... ]
//
// <asterisked identifier>    ::=   <identifier>
//
// <all fields reference>    ::=   <value expression primary> <period> <asterisk> [ AS <left paren> <all fields column name list> <right paren> ]
//
// <all fields column name list>    ::=   <column name list>

type SelectList struct {
	Asterisk bool
	Sublists []SelectSublist
}

func (l *SelectList) ArgCount(count *int) {
	if !l.Asterisk {
		for _, sub := range l.Sublists {
			sub.ArgCount(count)
		}
	}
}

type SelectSublist struct {
	Asterisk      bool
	DerivedColumn *DerivedColumn
}

func (l *SelectSublist) ArgCount(count *int) {
	if l.DerivedColumn != nil {
		l.DerivedColumn.ArgCount(count)
	}
}
