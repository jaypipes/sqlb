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

func (t *Table) Size() int {
    size := t.def.Size()
    if t.alias != "" {
        size += SYM_AS_LEN + len(t.alias)
    }
    return size
}

func (t *Table) Scan(b []byte) int {
    written := t.def.Scan(b)
    if t.alias != "" {
        copy(b[written:], SYM_AS)
        written += SYM_AS_LEN
        nalias := copy(b[written:], t.alias)
        written += nalias
    }
    return written
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

func (tl *TableList) Size() int {
    size := 0
    ntables := len(tl.tables)
    for _, t := range tl.tables {
        size += t.Size()
    }
    size += (SYM_COMMA_WS_LEN * (ntables - 1))  // Add in the commas
    return size
}

func (tl *TableList) Scan(b []byte) int {
    ntables := len(tl.tables)
    written := 0
    for x, t := range tl.tables {
        written += t.Scan(b[written:])
        if x != (ntables - 1) {
            copy(b[written:], SYM_COMMA_WS)
            written += SYM_COMMA_WS_LEN
        }
    }
    return written
}
