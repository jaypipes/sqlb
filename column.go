package sqlb

type Column struct {
    alias string
    def *ColumnDef
}

func (c *Column) ArgCount() int {
    return 0
}

func (c *Column) Size() int {
    size := c.def.Size()
    if c.alias != "" {
        size += SYM_AS_LEN + len(c.alias)
    }
    return size
}

func (c *Column) Scan(b []byte, args []interface{}) (int, int) {
    bw, _ := c.def.Scan(b, args)
    if c.alias != "" {
        bw += copy(b[bw:], SYM_AS)
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

type ColumnList struct {
    columns []*Column
}

func (cl *ColumnList) ArgCount() int {
    return 0
}

func (cl *ColumnList) Columns() []*Column {
    return cl.columns
}

func (cl *ColumnList) Size() int {
    size := 0
    ncols := len(cl.columns)
    for _, c := range cl.columns {
        size += c.Size()
    }
    size += (SYM_COMMA_WS_LEN * (ncols - 1))
    return size
}

func (cl *ColumnList) Scan(b []byte, args []interface{}) (int, int) {
    ncols := len(cl.columns)
    bw  := 0
    for x, c := range cl.columns {
        cbw, _ := c.Scan(b[bw:], args)
        bw += cbw
        if x != (ncols - 1) {
            bw += copy(b[bw:], SYM_COMMA_WS)
        }
    }
    return bw, 0
}
