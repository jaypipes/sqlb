//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import (
	"github.com/jaypipes/sqlb/pkg/grammar"
	"github.com/jaypipes/sqlb/pkg/types"
)

type havingClause struct {
	conditions []*Expression
}

func (c *havingClause) ArgCount() int {
	argc := 0
	for _, condition := range c.conditions {
		argc += condition.ArgCount()
	}
	return argc
}

func (c *havingClause) Size(scanner types.Scanner) int {
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

func (c *havingClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
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
