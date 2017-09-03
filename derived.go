package sqlb

// A type of selection that represents a SELECT in the FROM clause
type derivedTable struct {
    alias string
    from *selectClause
}

func (dt *derivedTable) selectionId() uint64 {
    return toId(dt.alias)
}

func (dt *derivedTable) projections() []projection {
    nprojs := len(dt.from.projs)
    projs := make([]projection, nprojs)
    for x := 0; x < nprojs; x++ {
        p := dt.from.projs[x]
        switch p.(type) {
            case *Column:
            case *ColumnDef:
        }
    }
    return projs
}

func (dt *derivedTable) argCount() int {
    return 0
}

func (dt *derivedTable) size() int {
    size := dt.from.size()
    size += (len(Symbols[SYM_LPAREN]) + len(Symbols[SYM_RPAREN]) +
             len(Symbols[SYM_AS]) + len(dt.alias))
    return size
}

func (dt *derivedTable) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_LPAREN])
    sbw, sac := dt.from.scan(b[bw:], args[ac:])
    bw += sbw
    ac += sac
    bw += copy(b[bw:], Symbols[SYM_RPAREN])
    bw += copy(b[bw:], Symbols[SYM_AS])
    bw += copy(b[bw:], dt.alias)
    return bw, ac
}

// A derivedColumn is a type of projection that is produced from a derived
// table (SELECT in the FROM clause). What makes a derived column unique is
// that it uses the alias of the underlying column as its name in the outer
// projection.
//
// The inner projection is a column against an underlying table or table def.
// The outer projection will have the selection alias of the derived table and
// the name of the projection will be the alias or name of the underlying
// column. For example, given the following SQL:
//
// SELECT <outer> FROM (
//   SELECT users.id, users.name FROM users
// ) AS u
//
// <outer> should contain:
//
// &derivedColumn{dt: dt, c: &Column{name: "id", tbl: users}},
// &derivedColumn{dt: dt, c: &Column{name: "name", tbl: users}}
//
// when scanned into <outer>, that should produce:
//
// []byte("u.id, u.name")
//
// However, let's say that the inner projections have been
// aliased, like so:
//
// SELECT <outer> FROM (
//   SELECT
//     users.id AS user_id,
//     users.name AS user_name
//   FROM users
// ) AS u
//
// <outer> should instead contain:
//
// &derivedColumn{dt: dt, c: &Column{alias: "user_id". name: "id", tbl: users}},
// &derivedColumn{dt: dt, c: &Column{alias: "user_name", name: "name", tbl: users}}
//
// which, when scanned into <outer>, should produce:
//
// []byte("u.user_id, u.user_name")
//
// Finally, the derivedColumn can itself have an alias, which can result in the
// outermost projection looking like so:
//
// SELECT u.user_name AS uname FROM (
//   SELECT users.name AS user_name
//   FROM users
// ) AS u
type derivedColumn struct {
    alias string  // This is the outermost alias
    c *Column
    dt *derivedTable
}

func (dc *derivedColumn) from() selection {
    return dc.dt
}

func (dc *derivedColumn) projectionId() uint64 {
    args := make([]string, 2)
    args[0] = dc.dt.alias
    if dc.alias != "" {
        args[1] = dc.alias
    } else if dc.c.alias != "" {
        args[1] = dc.c.alias
    } else {
        args[1] = dc.c.cdef.name
    }
    return toId(args...)
}

func (dc *derivedColumn) disableAliasScan() func() {
    origAlias := dc.alias
    dc.alias = ""
    return func() {dc.alias = origAlias}
}

func (dc *derivedColumn) argCount() int {
    return 0
}

func (dc *derivedColumn) size() int {
    size := dc.dt.size()
    size += len(Symbols[SYM_PERIOD])
    size += dc.c.size()
    if dc.alias != "" {
        size += len(Symbols[SYM_AS]) + len(dc.alias)
    }
    return size
}

func (dc *derivedColumn) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], dc.dt.alias)
    bw += copy(b[bw:], Symbols[SYM_PERIOD])
    reset := dc.c.disableAliasScan()
    defer reset()
    cbw, cac := dc.c.scan(b[bw:], args[ac:])
    bw += cbw
    ac += cac
    if dc.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], dc.alias)
    }
    return bw, ac
}
