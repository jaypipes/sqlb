package sqlb

type havingClause struct {
	conditions []*Expression
}

func (c *havingClause) argCount() int {
	argc := 0
	for _, condition := range c.conditions {
		argc += condition.argCount()
	}
	return argc
}

func (c *havingClause) size(scanner *sqlScanner) int {
	size := 0
	nconditions := len(c.conditions)
	if nconditions > 0 {
		size += len(scanner.format.SeparateClauseWith)
		size += len(Symbols[SYM_HAVING])
		size += len(Symbols[SYM_AND]) * (nconditions - 1)
		for _, condition := range c.conditions {
			size += condition.size(scanner)
		}
	}
	return size
}

func (c *havingClause) scan(scanner *sqlScanner, b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if len(c.conditions) > 0 {
		bw += copy(b[bw:], scanner.format.SeparateClauseWith)
		bw += copy(b[bw:], Symbols[SYM_HAVING])
		for x, condition := range c.conditions {
			if x > 0 {
				bw += copy(b[bw:], Symbols[SYM_AND])
			}
			bw += condition.scan(scanner, b[bw:], args, curArg)
		}
	}
	return bw
}
