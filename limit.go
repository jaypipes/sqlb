package sqlb

type LimitClause struct {
    limit int
    offset *int
}

func (lc *LimitClause) argCount() int {
    if lc.offset == nil {
        return 1
    }
    return 2
}

func (lc *LimitClause) size() int {
    size := len(Symbols[SYM_LIMIT]) + len(Symbols[SYM_QUEST_MARK])
    if lc.offset != nil {
        size += len(Symbols[SYM_OFFSET]) + len(Symbols[SYM_QUEST_MARK])
    }
    return size
}

func (lc *LimitClause) scan(b []byte, args []interface{}) (int, int) {
    var bw, ac int
    bw += copy(b[bw:], Symbols[SYM_LIMIT])
    bw += copy(b[bw:], Symbols[SYM_QUEST_MARK])
    args[ac] = lc.limit
    ac++
    if lc.offset != nil {
        bw += copy(b[bw:], Symbols[SYM_OFFSET])
        bw += copy(b[bw:], Symbols[SYM_QUEST_MARK])
        args[ac] = *lc.offset
        ac++
    }
    return bw, ac
}
