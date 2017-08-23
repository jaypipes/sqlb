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
)

var (
    // A static table containing information used in constructing the
    // expression's SQL string during Scan() calls
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
    }
)

type sqlFunc struct {
    alias string
    scanInfo scanInfo
    elements []Element
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

func (e *sqlFunc) ArgCount() int {
    ac := 0
    for _, el := range e.elements {
        ac += el.ArgCount()
    }
    return ac
}

func (f *sqlFunc) Size() int {
    size := 0
    elidx := 0
    for _, sym := range f.scanInfo {
        if sym == SYM_ELEMENT {
            el := f.elements[elidx]
            elidx++
            size += el.Size()
        } else {
            size += len(Symbols[sym])
        }
    }
    if f.alias != "" {
        size += len(Symbols[SYM_AS]) + len(f.alias)
    }
    return size
}

func (f *sqlFunc) Scan(b []byte, args []interface{}) (int, int) {
    bw, ac := 0, 0
    elidx := 0
    for _, sym := range f.scanInfo {
        if sym == SYM_ELEMENT {
            el := f.elements[elidx]
            elidx++
            ebw, eac := el.Scan(b[bw:], args[ac:])
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

func Max(el Element) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_MAX],
        elements: []Element{el},
    }
}

func (c *Column) Max() *sqlFunc {
    return Max(c)
}

func (c *ColumnDef) Max() *sqlFunc {
    return Max(c)
}

func Min(el Element) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_MIN],
        elements: []Element{el},
    }
}

func (c *Column) Min() *sqlFunc {
    return Min(c)
}

func (c *ColumnDef) Min() *sqlFunc {
    return Min(c)
}

func Sum(el Element) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_SUM],
        elements: []Element{el},
    }
}

func (c *Column) Sum() *sqlFunc {
    return Sum(c)
}

func (c *ColumnDef) Sum() *sqlFunc {
    return Sum(c)
}

func Avg(el Element) *sqlFunc {
    return &sqlFunc{
        scanInfo: funcScanTable[FUNC_AVG],
        elements: []Element{el},
    }
}

func (c *Column) Avg() *sqlFunc {
    return Avg(c)
}

func (c *ColumnDef) Avg() *sqlFunc {
    return Avg(c)
}
