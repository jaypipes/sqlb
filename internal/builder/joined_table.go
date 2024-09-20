//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
)

func (b *Builder) doJoinedTable(
	el *grammar.JoinedTable,
	qargs []interface{},
	curarg *int,
) {
	if el.Qualified != nil {
		b.doQualifiedJoin(el.Qualified, qargs, curarg)
	} else if el.Natural != nil {
		b.doNaturalJoin(el.Natural, qargs, curarg)
	} else if el.Union != nil {
		b.doUnionJoin(el.Union, qargs, curarg)
	} else if el.Cross != nil {
		b.doCrossJoin(el.Cross, qargs, curarg)
	}
}

func (b *Builder) doQualifiedJoin(
	el *grammar.QualifiedJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	switch el.Type {
	case grammar.JoinTypeInner:
		b.Write(grammar.Symbols[grammar.SYM_JOIN])
	case grammar.JoinTypeLeftOuter:
		b.Write(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case grammar.JoinTypeFullOuter:
		b.Write(grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	b.doTableReference(&el.Right, qargs, curarg)
	b.Write(grammar.Symbols[grammar.SYM_ON])
	b.doBooleanValueExpression(&el.On, qargs, curarg)
}

func (b *Builder) doNaturalJoin(
	el *grammar.NaturalJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_NATURAL])
	switch el.Type {
	case grammar.JoinTypeInner:
		b.Write(grammar.Symbols[grammar.SYM_JOIN])
	case grammar.JoinTypeLeftOuter:
		b.Write(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case grammar.JoinTypeFullOuter:
		b.Write(grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	b.doTablePrimary(&el.Right, qargs, curarg)
}

func (b *Builder) doUnionJoin(
	el *grammar.UnionJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_UNION])
	b.Write(grammar.Symbols[grammar.SYM_JOIN])
	b.doTablePrimary(&el.Right, qargs, curarg)
}

func (b *Builder) doCrossJoin(
	el *grammar.CrossJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.Write(grammar.Symbols[grammar.SYM_CROSS])
	b.Write(grammar.Symbols[grammar.SYM_JOIN])
	b.doTablePrimary(&el.Right, qargs, curarg)
}
