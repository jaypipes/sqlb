package sqlb

type Columnar interface {
    Column() *Column
}

type Column struct {
    alias string
    cdef *ColumnDef
    tbl *Table
}

func (c *Column) Column() *Column {
    return c
}

func (c *Column) ArgCount() int {
    return 0
}

func (c *Column) Size() int {
    size := c.cdef.Size()
    if c.alias != "" {
        size += len(Symbols[SYM_AS]) + len(c.alias)
    }
    return size
}

func (c *Column) Scan(b []byte, args []interface{}) (int, int) {
    bw, _ := c.cdef.Scan(b, args)
    if c.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], c.alias)
    }
    return bw, 0
}

func (c *Column) Alias(alias string) {
    c.alias = alias
}

func (c *Column) As(alias string) *Column {
    c.Alias(alias)
    return c
}

func isColumn(el Element) bool {
    switch el.(type) {
    case *Column:
        return true
    default:
        return false
    }
}

type ColumnDef struct {
    name string
    tdef *TableDef
}

func (cd *ColumnDef) Column() *Column {
    return &Column{
        cdef: cd,
        tbl: &Table{
            tdef: cd.tdef,
        },
    }
}

func (cd *ColumnDef) ArgCount() int {
    return 0
}

func (cd *ColumnDef) Size() int {
    return len(cd.name)
}

func (cd *ColumnDef) Scan(b []byte, args []interface{}) (int, int) {
    return copy(b, cd.name), 0
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
