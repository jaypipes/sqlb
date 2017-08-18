package sqlb

var (
    SYM_QM = []byte("?")
    SYM_QM_LEN = 1
    SYM_AS = []byte(" AS ")
    SYM_AS_LEN = 4
    SYM_COMMA_WS = []byte(", ")
    SYM_COMMA_WS_LEN = 2
    SYM_SELECT = []byte("SELECT ")
    SYM_SELECT_LEN = 7
    SYM_FROM = []byte(" FROM ")
    SYM_FROM_LEN = 6
    SYM_LPAREN = []byte("(")
    SYM_LPAREN_LEN = 1
    SYM_RPAREN = []byte(")")
    SYM_RPAREN_LEN = 1
    SYM_IN = []byte(" IN (")
    SYM_IN_LEN = 5
    SYM_WHERE = []byte(" WHERE ")
    SYM_WHERE_LEN = 7
    SYM_AND = []byte(" AND ")
    SYM_AND_LEN = 5

    SYM_OP = map[Op][]byte{
        OP_EQUAL: []byte(" = "),
        OP_NEQUAL: []byte(" != "),
        OP_AND: []byte(" AND "),
        OP_OR: []byte(" OR "),
    }
)
