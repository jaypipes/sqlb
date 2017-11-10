package sqlb

type Table struct {
	alias   string
	meta    *Meta
	name    string
	columns []*Column
}

// Sets the statement's dialect and pushes the dialect down into any of the
// statement's sub-clauses
func (t *Table) setDialect(dialect Dialect) {
	for _, c := range t.columns {
		c.setDialect(dialect)
	}
}

// Return a pointer to a Column with a name or alias matching the supplied
// string, or nil if no such column is known
func (t *Table) C(name string) *Column {
	for _, c := range t.columns {
		if c.name == name || c.alias == name {
			return c
		}
	}
	return nil
}

func (t *Table) NewColumn(name string) *Column {
	c := t.C(name)
	if c != nil {
		return c
	}
	c = &Column{
		name: name,
		tbl:  t,
	}
	t.columns = append(t.columns, c)
	return c
}

func (t *Table) projections() []projection {
	res := make([]projection, len(t.columns))
	for x, c := range t.columns {
		res[x] = c
	}
	return res
}

func (t *Table) argCount() int {
	return 0
}

func (t *Table) size() int {
	size := len(t.name)
	if t.alias != "" {
		size += len(Symbols[SYM_AS]) + len(t.alias)
	}
	return size
}

func (t *Table) scan(b []byte, args []interface{}, curArg *int) int {
	bw := copy(b, t.name)
	if t.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], t.alias)
	}
	return bw
}

func (t *Table) As(alias string) *Table {
	cols := make([]*Column, len(t.columns))
	tbl := &Table{
		alias: alias,
		name:  t.name,
		meta:  t.meta,
	}
	for x, c := range t.columns {
		cols[x] = &Column{
			alias: c.alias,
			name:  c.name,
			tbl:   tbl,
		}
	}
	tbl.columns = cols
	return tbl
}
