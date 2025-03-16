//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package builder

import (
	"github.com/jaypipes/sqlb/core/grammar"
	"github.com/jaypipes/sqlb/core/grammar/symbol"
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
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
	case grammar.JoinTypeLeftOuter:
		b.WriteString(symbol.Left)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
	case grammar.JoinTypeFullOuter:
		b.WriteString(symbol.Cross)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
	}
	b.doTableReference(&el.Right, qargs, curarg)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.On)
	b.WriteString(symbol.Space)
	b.doBooleanValueExpression(&el.On, qargs, curarg)
}

func (b *Builder) doNaturalJoin(
	el *grammar.NaturalJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Natural)
	b.WriteString(symbol.Space)
	switch el.Type {
	case grammar.JoinTypeInner:
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
	case grammar.JoinTypeLeftOuter:
		b.WriteString(symbol.Left)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
	case grammar.JoinTypeFullOuter:
		b.WriteString(symbol.Cross)
		b.WriteString(symbol.Space)
		b.WriteString(symbol.Join)
		b.WriteString(symbol.Space)
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
	b.WriteString(symbol.Union)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Join)
	b.WriteString(symbol.Space)
	b.doTablePrimary(&el.Right, qargs, curarg)
}

func (b *Builder) doCrossJoin(
	el *grammar.CrossJoin,
	qargs []interface{},
	curarg *int,
) {
	b.doTableReference(&el.Left, qargs, curarg)
	b.WriteString(b.opts.FormatSeparateClauseWith())
	b.WriteString(symbol.Cross)
	b.WriteString(symbol.Space)
	b.WriteString(symbol.Join)
	b.WriteString(symbol.Space)
	b.doTablePrimary(&el.Right, qargs, curarg)
}
