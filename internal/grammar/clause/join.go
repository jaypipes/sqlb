//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/api"
	"github.com/jaypipes/sqlb/internal/grammar"
	"github.com/jaypipes/sqlb/internal/grammar/expression"
)

type Join struct {
	joinType api.JoinType
	left     api.Selection
	right    api.Selection
	on       *expression.Expression
}

func (j *Join) Left() api.Selection {
	return j.left
}

func (j *Join) Right() api.Selection {
	return j.right
}

func (j *Join) ArgCount() int {
	ac := 0
	if j.on != nil {
		ac = j.on.ArgCount()
	}
	return ac + j.left.ArgCount() + j.right.ArgCount()
}

func (j *Join) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	b.WriteString(opts.FormatSeparateClauseWith())
	switch j.joinType {
	case api.JoinInner:
		b.Write(grammar.Symbols[grammar.SYM_JOIN])
	case api.JoinOuter:
		b.Write(grammar.Symbols[grammar.SYM_LEFT_JOIN])
	case api.JoinCross:
		b.Write(grammar.Symbols[grammar.SYM_CROSS_JOIN])
	}
	b.WriteString(j.right.String(opts, qargs, curarg))
	if j.on != nil {
		b.Write(grammar.Symbols[grammar.SYM_ON])
		b.WriteString(j.on.String(opts, qargs, curarg))
	}
	return b.String()
}

func InnerJoin(left api.Selection, right api.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: api.JoinInner,
		left:     left,
		right:    right,
		on:       on,
	}
}

func OuterJoin(left api.Selection, right api.Selection, on *expression.Expression) *Join {
	return &Join{
		joinType: api.JoinOuter,
		left:     left,
		right:    right,
		on:       on,
	}
}

func CrossJoin(left api.Selection, right api.Selection) *Join {
	return &Join{
		joinType: api.JoinCross,
		left:     left,
		right:    right,
	}
}

func NewJoin(
	jt api.JoinType,
	left api.Selection,
	right api.Selection,
	on *expression.Expression,
) *Join {
	switch jt {
	case api.JoinInner:
		return InnerJoin(left, right, on)
	case api.JoinOuter:
		return OuterJoin(left, right, on)
	case api.JoinCross:
		return CrossJoin(left, right)
	}
	return nil
}
