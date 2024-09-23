//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doQueryExpression(
	el *grammar.QueryExpression,
	qargs []interface{},
	curarg *int,
) {
	body := el.Body
	if body.NonJoin != nil {
		b.doNonJoinQueryExpression(body.NonJoin, qargs, curarg)
	} else if body.Joined != nil {
		b.doJoinedTable(body.Joined, qargs, curarg)
	}
}

func (b *Builder) doNonJoinQueryExpression(
	el *grammar.NonJoinQueryExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.NonJoin != nil {
		b.doNonJoinQueryTerm(el.NonJoin, qargs, curarg)
	}
}

func (b *Builder) doNonJoinQueryTerm(
	el *grammar.NonJoinQueryTerm,
	qargs []interface{},
	curarg *int,
) {
	if el.Primary != nil {
		b.doNonJoinQueryPrimary(el.Primary, qargs, curarg)
	}
}

func (b *Builder) doNonJoinQueryPrimary(
	el *grammar.NonJoinQueryPrimary,
	qargs []interface{},
	curarg *int,
) {
	if el.Simple != nil {
		b.doSimpleTable(el.Simple, qargs, curarg)
	} else if el.Parenthesized != nil {
		b.Write(grammar.Symbols[grammar.SYM_LPAREN])
		b.doNonJoinQueryExpression(el.Parenthesized, qargs, curarg)
		b.Write(grammar.Symbols[grammar.SYM_RPAREN])
	}
}

func (b *Builder) doSimpleTable(
	el *grammar.SimpleTable,
	qargs []interface{},
	curarg *int,
) {
	if el.QuerySpecification != nil {
		b.doQuerySpecification(el.QuerySpecification, qargs, curarg)
	}
}
