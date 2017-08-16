package sqlb

type Table struct {
    alias string
    def *TableDef
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
        copy(b[written:], []byte(SYM_AS))
        written += SYM_AS_LEN
        nalias := copy(b[written:], []byte(t.alias))
        written += nalias
    }
    return written
}
