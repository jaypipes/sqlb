package sqlb

type Table struct {
    alias string
    tdef *TableDef
}

func (t *Table) selectionId() uint64 {
    if t.alias != "" {
        return toId(t.alias)
    }
    return t.tdef.selectionId()
}

func (t *Table) Column(name string) *Column {
    for _, cdef := range t.tdef.ColumnDefs() {
        if name == cdef.name {
            return &Column{
                tbl: t,
                cdef: cdef,
            }
        }
    }
    return nil
}

func (t *Table) Columns() []*Column {
    cdefs := t.tdef.ColumnDefs()
    ncols := len(cdefs)
    cols := make([]*Column, ncols)
    for x := 0; x < ncols; x++ {
        cols[x] = &Column{
            tbl: t,
            cdef: cdefs[x],
        }
    }
    return cols
}

func (t *Table) ArgCount() int {
    return 0
}

func (t *Table) Size() int {
    size := t.tdef.Size()
    if t.alias != "" {
        size += len(Symbols[SYM_AS]) + len(t.alias)
    }
    return size
}

func (t *Table) Scan(b []byte, args []interface{}) (int, int) {
    bw, _ := t.tdef.Scan(b, args)
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

type TableDef struct {
    name string
    schema string
    cdefs []*ColumnDef
}

func (td *TableDef) selectionId() uint64 {
    return toId(td.schema, td.name)
}

func (td *TableDef) Table() *Table {
    return &Table{tdef: td}
}

func (td *TableDef) ArgCount() int {
    return 0
}

func (td *TableDef) Size() int {
    return len(td.name)
}

func (td *TableDef) Scan(b []byte, args []interface{}) (int, int) {
    return copy(b, td.name), 0
}

// Generate an aliased Table from a TableDef
func (td *TableDef) As(alias string) *Table {
    return &Table{tdef: td, alias: alias}
}

func (td *TableDef) Column(name string) *Column {
    for _, cdef := range td.cdefs {
        if cdef.name == name {
            return &Column{
                cdef: cdef,
                tbl: &Table{
                    tdef: td,
                },
            }
        }
    }
    return nil
}

func (td *TableDef) ColumnDefs() []*ColumnDef {
    return td.cdefs
}
