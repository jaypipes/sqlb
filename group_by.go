//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//
package sqlb

import "github.com/jaypipes/sqlb/pkg/types"

type groupByClause struct {
	cols []types.Projection
}

func (gb *groupByClause) ArgCount() int {
	argc := 0
	return argc
}

func (gb *groupByClause) Size(scanner types.Scanner) int {
	size := 0
	size += len(scanner.FormatOptions().SeparateClauseWith)
	size += len(Symbols[SYM_GROUP_BY])
	ncols := len(gb.cols)
	for _, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		size += c.Size(scanner)
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (gb *groupByClause) Scan(scanner types.Scanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.FormatOptions().SeparateClauseWith)
	bw += copy(b[bw:], Symbols[SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.DisableAliasScan()
		defer reset()
		bw += c.Scan(scanner, b[bw:], args, curArg)
		if x != (ncols - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	return bw
}
