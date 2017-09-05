package sqlb

import (
    "fmt"
)

type funcId int

const (
    FUNC_MAX funcId = iota
    FUNC_MIN
    FUNC_SUM
    FUNC_AVG
    FUNC_COUNT_STAR
    FUNC_COUNT_DISTINCT
    FUNC_CAST
    FUNC_TRIM
    FUNC_CHAR_LENGTH
    FUNC_BIT_LENGTH
    FUNC_ASCII
)

var (
    // A static table containing information used in constructing the
    // expression's SQL string during scan() calls
    funcScanTable = map[funcId]scanInfo{
        FUNC_MAX: scanInfo{
            SYM_MAX, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_MIN: scanInfo{
            SYM_MIN, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_SUM: scanInfo{
            SYM_SUM, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_AVG: scanInfo{
            SYM_AVG, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_COUNT_STAR: scanInfo{
            SYM_COUNT_STAR,
        },
        FUNC_COUNT_DISTINCT: scanInfo{
            SYM_COUNT_DISTINCT, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_CAST: scanInfo{
            SYM_CAST, SYM_ELEMENT, SYM_AS, SYM_PLACEHOLDER, SYM_RPAREN,
        },
        FUNC_TRIM: scanInfo{
            SYM_TRIM, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_CHAR_LENGTH: scanInfo{
            SYM_CHAR_LENGTH, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_BIT_LENGTH: scanInfo{
            SYM_BIT_LENGTH, SYM_ELEMENT, SYM_RPAREN,
        },
        FUNC_ASCII: scanInfo{
            SYM_ASCII, SYM_ELEMENT, SYM_RPAREN,
        },
    }
)

type sqlFunc struct {
    sel selection
    alias string
    scanInfo scanInfo
    elements []element
}

func (f *sqlFunc) from() selection {
    return f.sel
}

func (f *sqlFunc) disableAliasScan() func() {
    origAlias := f.alias
    f.alias = ""
    return func() {f.alias = origAlias}
}

func (f *sqlFunc) projectionId() uint64 {
    // Each construction of a function is unique, so here we cheat and just
    // return the hash of the struct's address in memory
    return toId(fmt.Sprintf("%p", f))
}

func (f *sqlFunc) Alias(alias string) {
    f.alias = alias
}

func (f *sqlFunc) As(alias string) *sqlFunc {
    f.Alias(alias)
    return f
}

func (e *sqlFunc) argCount() int {
    ac := 0
    for _, el := range e.elements {
        ac += el.argCount()
    }
    return ac
}

func (f *sqlFunc) size() int {
    size := 0
    elidx := 0
    for _, sym := range f.scanInfo {
        if sym == SYM_ELEMENT {
            el := f.elements[elidx]
            // We need to disable alias output for elements that are
            // projections. We don't want to output, for example,
            // "ON users.id AS user_id = articles.author"
            switch el.(type) {
            case projection:
                reset := el.(projection).disableAliasScan()
                defer reset()
            }
            elidx++
            size += el.size()
        } else {
            size += len(Symbols[sym])
        }
    }
    if f.alias != "" {
        size += len(Symbols[SYM_AS]) + len(f.alias)
    }
    return size
}

func (f *sqlFunc) scan(b []byte, args []interface{}) (int, int) {
    bw, ac := 0, 0
    elidx := 0
    for _, sym := range f.scanInfo {
        if sym == SYM_ELEMENT {
            el := f.elements[elidx]
            // We need to disable alias output for elements that are
            // projections. We don't want to output, for example,
            // "ON users.id AS user_id = articles.author"
            switch el.(type) {
            case projection:
                reset := el.(projection).disableAliasScan()
                defer reset()
            }
            elidx++
            ebw, eac := el.scan(b[bw:], args[ac:])
            bw += ebw
            ac += eac
        } else {
            bw += copy(b[bw:], Symbols[sym])
        }
    }
    if f.alias != "" {
        bw += copy(b[bw:], Symbols[SYM_AS])
        bw += copy(b[bw:], f.alias)
    }
    return bw, ac
}

func Max(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_MAX],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Max() *sqlFunc {
    return Max(c)
}

func (c *ColumnDef) Max() *sqlFunc {
    return Max(c)
}

func Min(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_MIN],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Min() *sqlFunc {
    return Min(c)
}

func (c *ColumnDef) Min() *sqlFunc {
    return Min(c)
}

func Sum(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_SUM],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Sum() *sqlFunc {
    return Sum(c)
}

func (c *ColumnDef) Sum() *sqlFunc {
    return Sum(c)
}

func Avg(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_AVG],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Avg() *sqlFunc {
    return Avg(c)
}

func (c *ColumnDef) Avg() *sqlFunc {
    return Avg(c)
}

func Count() *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_COUNT_STAR],
    }
}

func CountDistinct(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_COUNT_DISTINCT],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func Cast(p projection, stype SqlType) *sqlFunc {
    si := make([]Symbol, len(funcScanTable[FUNC_CAST]))
    copy(si, funcScanTable[FUNC_CAST])
    // Replace the placeholder with the SQL type's appropriate []byte
    // representation
    si[3] = sqlTypeToSymbol[stype]
    return &sqlFunc{
        scanInfo: si,
        elements: []element{p.(element)},
    }
}

func Trim(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_TRIM],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Trim() *sqlFunc {
    return Trim(c)
}

func (c *ColumnDef) Trim() *sqlFunc {
    return Trim(c)
}

func CharLength(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_CHAR_LENGTH],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) CharLength() *sqlFunc {
    return CharLength(c)
}

func (c *ColumnDef) CharLength() *sqlFunc {
    return CharLength(c)
}

func BitLength(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_BIT_LENGTH],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) BitLength() *sqlFunc {
    return BitLength(c)
}

func (c *ColumnDef) BitLength() *sqlFunc {
    return BitLength(c)
}

func Ascii(p projection) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_ASCII],
        elements: []element{p.(element)},
        sel: p.from(),
    }
}

func (c *Column) Ascii() *sqlFunc {
    return Ascii(c)
}

func (c *ColumnDef) Ascii() *sqlFunc {
    return Ascii(c)
}
