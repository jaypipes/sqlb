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

type Having struct {
	conditions []*expression.Expression
}

func (c *Having) Conditions() []*expression.Expression {
	return c.conditions
}

func (c *Having) AddCondition(e *expression.Expression) {
	c.conditions = append(c.conditions, e)
}

func (c *Having) ArgCount() int {
	argc := 0
	for _, condition := range c.conditions {
		argc += condition.ArgCount()
	}
	return argc
}

func (c *Having) String(
	opts api.Options,
	qargs []interface{},
	curarg *int,
) string {
	b := &strings.Builder{}
	if len(c.conditions) > 0 {
		b.WriteString(opts.FormatSeparateClauseWith())
		b.Write(grammar.Symbols[grammar.SYM_HAVING])
		for x, condition := range c.conditions {
			if x > 0 {
				b.Write(grammar.Symbols[grammar.SYM_AND])
			}
			b.WriteString(condition.String(opts, qargs, curarg))
		}
	}
	return b.String()
}

// NewHaving returns a new Having populated with zero or more
// Expression conditions
func NewHaving(conds ...*expression.Expression) *Having {
	return &Having{
		conditions: conds,
	}
}
