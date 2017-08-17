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
        size += SYM_AS_LEN + len(t.alias)
    }
    return size
}

func (t *Table) Scan(b []byte, args []interface{}) (int, int) {
    bw, _ := t.def.Scan(b, args)
    if t.alias != "" {
        bw += copy(b[bw:], SYM_AS)
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

type TableList struct {
    tables []*Table
}

func (tl *TableList) ArgCount() int {
    return 0
}

func (tl *TableList) Size() int {
    size := 0
    ntables := len(tl.tables)
    for _, t := range tl.tables {
        size += t.Size()
    }
    size += (SYM_COMMA_WS_LEN * (ntables - 1))  // Add in the commas
    return size
}

func (tl *TableList) Scan(b []byte, args []interface{}) (int, int) {
    ntables := len(tl.tables)
    bw := 0
    for x, t := range tl.tables {
        tbw, _ := t.Scan(b[bw:], args)
        bw += tbw
        if x != (ntables - 1) {
            bw += copy(b[bw:], SYM_COMMA_WS)
        }
    }
    return bw, 0
}
