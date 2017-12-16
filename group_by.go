package sqlb

type groupByClause struct {
	cols []projection
}

func (gb *groupByClause) argCount() int {
	argc := 0
	return argc
}

func (gb *groupByClause) size(scanner *sqlScanner) int {
	size := 0
	size += len(scanner.format.SeparateClauseWith)
	size += len(Symbols[SYM_GROUP_BY])
	ncols := len(gb.cols)
	for _, c := range gb.cols {
		reset := c.disableAliasScan()
		defer reset()
		size += c.size(scanner)
	}
	return size + (len(Symbols[SYM_COMMA_WS]) * (ncols - 1)) // the commas...
}

func (gb *groupByClause) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	bw += copy(b[bw:], scanner.format.SeparateClauseWith)
	bw += copy(b[bw:], Symbols[SYM_GROUP_BY])
	ncols := len(gb.cols)
	for x, c := range gb.cols {
		reset := c.disableAliasScan()
		defer reset()
		bw += c.scan(scanner, b[bw:], args, curArg)
		if x != (ncols - 1) {
			bw += copy(b[bw:], Symbols[SYM_COMMA_WS])
		}
	}
	return bw
}
