package sqlb

type Symbol int

const (
    SYM_ELEMENT = iota // Marker for an element that self-scans into the SQL buffer
    SYM_QUEST_MARK
    SYM_AS
    SYM_COMMA_WS
    SYM_SELECT
    SYM_FROM
    SYM_LPAREN
    SYM_RPAREN
    SYM_IN
    SYM_WHERE
    SYM_AND
    SYM_OR
    SYM_EQUAL
    SYM_NEQUAL
    SYM_BETWEEN
    SYM_LIMIT
    SYM_OFFSET
    SYM_ORDER_BY
    SYM_DESC
    SYM_GROUP_BY
)

var (
    Symbols = map[Symbol][]byte{
        SYM_QUEST_MARK: []byte("?"),
        SYM_AS: []byte(" AS "),
        SYM_COMMA_WS: []byte(", "),
        SYM_SELECT: []byte("SELECT "),
        SYM_FROM: []byte(" FROM "),
        SYM_LPAREN: []byte("("),
        SYM_RPAREN: []byte(")"),
        SYM_IN: []byte(" IN ("),
        SYM_WHERE: []byte(" WHERE "),
        SYM_AND: []byte(" AND "),
        SYM_OR: []byte(" OR "),
        SYM_EQUAL: []byte(" = "),
        SYM_NEQUAL: []byte(" != "),
        SYM_BETWEEN: []byte(" BETWEEN "),
        SYM_LIMIT: []byte(" LIMIT "),
        SYM_OFFSET: []byte(" OFFSET "),
        SYM_ORDER_BY: []byte(" ORDER BY "),
        SYM_DESC: []byte(" DESC"),
        SYM_GROUP_BY: []byte(" GROUP BY "),
    }
)
