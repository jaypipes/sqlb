package sqlb

type Table struct {
    alias string
    def *TableDef
}

func (t *Table) Columns() []*Column {
    cdefs := t.def.ColumnDefs()
    ncols := len(cdefs)
    cols := make([]*Column, ncols)
    for x := 0; x < ncols; x++ {
        cols[x] = &Column{def: cdefs[x]}
    }
    return cols
}

func (t *Table) ArgCount() int {
    return 0
}

func (t *Table) Size() int {
    size := t.def.Size()
    if t.alias != "" {
        size += len(Symbols[SYM_AS]) + len(t.alias)
    }
    return size
}

func (t *Table) Scan(b []byte, args []interface{}) (int, int) {
    bw, _ := t.def.Scan(b, args)
    if t.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], t.alias)
    }
    return bw, 0
}

func (t *Table) Alias(alias string) {
    t.alias = alias
}

func (t *Table) As(alias string) *Table {
    t.Alias(alias)
    return t
}
