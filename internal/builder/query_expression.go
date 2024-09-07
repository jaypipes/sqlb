//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/grammar"
)

func (b *Builder) doQueryExpression(
	el *grammar.QueryExpression,
	qargs []interface{},
	curarg *int,
) {
	body := el.Body
	if body.NonJoinQueryExpression != nil {
		b.doNonJoinQueryExpression(body.NonJoinQueryExpression, qargs, curarg)
	} else if body.JoinedTable != nil {
		b.doJoinedTable(body.JoinedTable, qargs, curarg)
	}
}

func (b *Builder) doNonJoinQueryExpression(
	el *grammar.NonJoinQueryExpression,
	qargs []interface{},
	curarg *int,
) {
	if el.NonJoinQueryTerm != nil {
		b.doNonJoinQueryTerm(el.NonJoinQueryTerm, qargs, curarg)
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
	if el.SimpleTable != nil {
		b.doSimpleTable(el.SimpleTable, qargs, curarg)
	} else if el.ParenthesizedNonJoinQueryExpression != nil {
		b.Write(grammar.Symbols[grammar.SYM_LPAREN])
		b.doNonJoinQueryExpression(el.ParenthesizedNonJoinQueryExpression, qargs, curarg)
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
