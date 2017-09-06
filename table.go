package sqlb

type Table struct {
    alias string
    tdef *TableDef
}

func (t *Table) idParts() []string {
    if t.alias != "" {
        return []string{t.alias}
    }
    return t.tdef.idParts()
}

func (t *Table) Column(name string) *Column {
    for _, cdef := range t.tdef.cdefs {
        if name == cdef.name {
            return &Column{
                tbl: t,
                cdef: cdef,
            }
        }
    }
    return nil
}

func (t *Table) projections() []projection {
    cdefs := t.tdef.cdefs
    ncols := len(cdefs)
    cols := make([]projection, ncols)
    for x := 0; x < ncols; x++ {
        cols[x] = &Column{
            tbl: t,
            cdef: cdefs[x],
        }
    }
    return cols
}

func (t *Table) argCount() int {
    return 0
}

func (t *Table) size() int {
    size := t.tdef.size()
    if t.alias != "" {
        size += len(Symbols[SYM_AS]) + len(t.alias)
    }
    return size
}

func (t *Table) scan(b []byte, args []interface{}) (int, int) {
    bw, _ := t.tdef.scan(b, args)
    if t.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], t.alias)
    }
    return bw, 0
}

func (t *Table) setAlias(alias string) {
    t.alias = alias
}

func (t *Table) As(alias string) *Table {
    t.setAlias(alias)
    return t
}

type TableDef struct {
    meta *Meta
    name string
    cdefs []*ColumnDef
}

func (td *TableDef) idParts() []string {
    return []string{td.meta.schemaName, td.name}
}

func (td *TableDef) Table() *Table {
    return &Table{tdef: td}
}

func (td *TableDef) argCount() int {
    return 0
}

func (td *TableDef) size() int {
    return len(td.name)
}

func (td *TableDef) scan(b []byte, args []interface{}) (int, int) {
    return copy(b, td.name), 0
}

// Generate an aliased Table from a TableDef
func (td *TableDef) As(alias string) *Table {
    return &Table{tdef: td, alias: alias}
}

// Return a pointer to a ColumnDef with a name matching the supplied string, or
// nil if no such column is known
func (td *TableDef) ColumnDef(name string) *ColumnDef {
    for _, cdef := range td.cdefs {
        if cdef.name == name {
            return cdef
        }
    }
    return nil
}

func (td *TableDef) NewColumnDef(name string) *ColumnDef {
    cd := td.ColumnDef(name)
    if cd != nil {
        return cd
    }
    cd = &ColumnDef{
        name: name,
        tdef: td,
    }
    td.cdefs = append(td.cdefs, cd)
    return cd
}

// Returns a pointer to a Column representing an aliasable version of the table
// column with the supplied name, or nil if no such column definition is known
func (td *TableDef) Column(name string) *Column {
    cd := td.ColumnDef(name)
    if cd == nil {
        return nil
    }
    return &Column{
        cdef: cd,
        tbl: td.Table(),
    }
}

func (td *TableDef) projections() []projection {
    res := make([]projection, len(td.cdefs))
    for x, cdef := range td.cdefs {
        res[x] = cdef
    }
    return res
}
