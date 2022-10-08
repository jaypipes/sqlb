//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package clause

import (
	"strings"

	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/grammar/expression"
	"github.com/jaypipes/sqlb/pkg/types"
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

func (c *Having) Size(scanner types.Scanner) int {
	size := 0
	nconditions := len(c.conditions)
	if nconditions > 0 {
		size += len(scanner.FormatOptions().SeparateClauseWith)
		size += len(grammar.Symbols[grammar.SYM_HAVING])
		size += len(grammar.Symbols[grammar.SYM_AND]) * (nconditions - 1)
		for _, condition := range c.conditions {
			size += condition.Size(scanner)
		}
	}
	return size
}

func (c *Having) Scan(scanner types.Scanner, b *strings.Builder, args []interface{}, curArg *int) {
	if len(c.conditions) > 0 {
		b.WriteString(scanner.FormatOptions().SeparateClauseWith)
		b.Write(grammar.Symbols[grammar.SYM_HAVING])
		for x, condition := range c.conditions {
			if x > 0 {
				b.Write(grammar.Symbols[grammar.SYM_AND])
			}
			condition.Scan(scanner, b, args, curArg)
		}
	}
}

// NewHaving returns a new Having populated with zero or more
// Expression conditions
func NewHaving(conds ...*expression.Expression) *Having {
	return &Having{
		conditions: conds,
	}
}
