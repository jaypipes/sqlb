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
        copy(b[written:], []byte(SYM_AS))
        written += SYM_AS_LEN
        nalias := copy(b[written:], []byte(c.alias))
        written += nalias
    }
    return written
}
