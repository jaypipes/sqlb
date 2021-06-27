//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ast

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type HavingClause struct {
	conditions []*Expression
}

func (c *HavingClause) Conditions() []*Expression {
	return c.conditions
}

func (c *HavingClause) AddCondition(e *Expression) {
	c.conditions = append(c.conditions, e)
}

func (c *HavingClause) ArgCount() int {
	argc := 0
	for _, condition := range c.conditions {
		argc += condition.ArgCount()
	}
	return argc
}

func (c *HavingClause) Size(scanner types.Scanner) int {
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

func (c *HavingClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(c.conditions) > 0 {
		bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
		bw += copy(b[bw:], grammar.Symbols[grammar.SYM_HAVING])
		for x, condition := range c.conditions {
			if x > 0 {
				bw += copy(b[bw:], grammar.Symbols[grammar.SYM_AND])
			}
			bw += condition.Scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}

// NewHavingClause returns a new HavingClause populated with zero or more
// Expression conditions
func NewHavingClause(conds ...*Expression) *HavingClause {
	return &HavingClause{
		conditions: conds,
	}
}
