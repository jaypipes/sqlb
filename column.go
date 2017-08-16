package sqlb

type Column struct {
    alias string
    def *ColumnDef
}

func (c *Column) Size() int {
    size := c.def.Size()
    if c.alias != "" {
        size += SYM_AS_LEN + len(c.alias)
    }
    return size
}

func (c *Column) Scan(b []byte) int {
    written := c.def.Scan(b)
    if c.alias != "" {
        copy(b[written:], SYM_AS)
        written += SYM_AS_LEN
        nalias := copy(b[written:], c.alias)
        written += nalias
    }
    return written
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

func (cl *ColumnList) Scan(b []byte) int {
    ncols := len(cl.columns)
    written := 0
    for x, c := range cl.columns {
        written += c.Scan(b[written:])
        if x != (ncols - 1) {
            copy(b[written:], SYM_COMMA_WS)
            written += SYM_COMMA_WS_LEN
        }
    }
    return written
}
