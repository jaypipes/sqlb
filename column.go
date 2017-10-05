package sqlb

type Column struct {
    alias string
    cdef *ColumnDef
    tbl *Table
}

func (c *Column) from() selection {
    return c.tbl
}

func (c *Column) disableAliasScan() func() {
    origAlias := c.alias
    c.alias = ""
    return func() {c.alias = origAlias}
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
        size += len(c.tbl.tdef.name)
    }
    size += len(Symbols[SYM_PERIOD])
    size += len(c.cdef.name)
    if c.alias != "" {
        size += len(Symbols[SYM_AS]) + len(c.alias)
    }
    return size
}

func (c *Column) scan(b []byte, args []interface{}) (int, int) {
    bw := 0
    if c.tbl.alias != "" {
        bw += copy(b[bw:], c.tbl.alias)
    } else {
        bw += copy(b[bw:], c.tbl.tdef.name)
    }
    bw += copy(b[bw:], Symbols[SYM_PERIOD])
    bw += copy(b[bw:], c.cdef.name)
    if c.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], c.alias)
    }
    return bw, 0
}

func (c *Column) setAlias(alias string) {
    c.alias = alias
}

func (c *Column) As(alias string) *Column {
    return &Column{
        alias: alias,
        tbl: c.tbl,
        cdef: c.cdef,
    }
}

type ColumnDef struct {
    name string
    tdef *TableDef
}

func (cd *ColumnDef) from() selection {
    return cd.tdef
}

// A column definition isn't aliasable...
func (cd *ColumnDef) disableAliasScan() func() {
    return func() {}
}

func (cd *ColumnDef) Column() *Column {
    return &Column{
        cdef: cd,
        tbl: &Table{
            tdef: cd.tdef,
        },
    }
}

func (cd *ColumnDef) argCount() int {
    return 0
}

func (cd *ColumnDef) size() int {
    return len(cd.tdef.name) + len(Symbols[SYM_PERIOD]) + len(cd.name)
}

func (cd *ColumnDef) scan(b []byte, args []interface{}) (int, int) {
    bw := copy(b, cd.tdef.name)
    bw += copy(b[bw:], Symbols[SYM_PERIOD])
    bw += copy(b[bw:], cd.name)
    return bw, 0
}

// Generate an aliased Column from a ColumnDef
func (cd *ColumnDef) As(alias string) *Column {
    return &Column{
        cdef: cd,
        alias: alias,
        tbl: &Table{
            tdef: cd.tdef,
        },
    }
}
