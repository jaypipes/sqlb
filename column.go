package sqlb

type Column struct {
	alias string
	name  string
	tbl   *Table
}

func (c *Column) from() selection {
	return c.tbl
}

func (c *Column) disableAliasScan() func() {
	origAlias := c.alias
	c.alias = ""
	return func() { c.alias = origAlias }
}

func (c *Column) Column() *Column {
	return c
}

func (c *Column) argCount() int {
	return 0
}

func (c *Column) size() int {
	size := 0
	if c.tbl.alias != "" {
		size += len(c.tbl.alias)
	} else {
		size += len(c.tbl.name)
	}
	size += len(Symbols[SYM_PERIOD])
	size += len(c.name)
	if c.alias != "" {
		size += len(Symbols[SYM_AS]) + len(c.alias)
	}
	return size
}

func (c *Column) scan(b []byte, args []interface{}, curArg *int) int {
	bw := 0
	if c.tbl.alias != "" {
		bw += copy(b[bw:], c.tbl.alias)
	} else {
		bw += copy(b[bw:], c.tbl.name)
	}
	bw += copy(b[bw:], Symbols[SYM_PERIOD])
	bw += copy(b[bw:], c.name)
	if c.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], c.alias)
	}
	return bw
}

func (c *Column) setAlias(alias string) {
	c.alias = alias
}

func (c *Column) As(alias string) *Column {
	return &Column{
		alias: alias,
		name:  c.name,
		tbl:   c.tbl,
	}
}
