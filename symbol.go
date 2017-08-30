package sqlb

type Symbol int
type scanInfo []Symbol

const (
    SYM_ELEMENT = iota // Marker for an element that self-scans into the SQL buffer
    SYM_QUEST_MARK
    SYM_PERIOD
    SYM_AS
    SYM_COMMA_WS
    SYM_SELECT
    SYM_FROM
    SYM_JOIN
    SYM_ON
    SYM_WHERE
    SYM_GROUP_BY
    SYM_ORDER_BY
    SYM_DESC
    SYM_LIMIT
    SYM_OFFSET
    SYM_LPAREN
    SYM_RPAREN
    SYM_IN
    SYM_AND
    SYM_OR
    SYM_EQUAL
    SYM_NEQUAL
    SYM_BETWEEN
    SYM_MAX
    SYM_MIN
    SYM_SUM
    SYM_AVG
    SYM_COUNT_STAR
    SYM_COUNT_DISTINCT
)

var (
    Symbols = map[Symbol][]byte{
        SYM_QUEST_MARK: []byte("?"),
        SYM_PERIOD: []byte("."),
        SYM_AS: []byte(" AS "),
        SYM_COMMA_WS: []byte(", "),
        SYM_SELECT: []byte("SELECT "),
        SYM_FROM: []byte(" FROM "),
        SYM_JOIN: []byte(" JOIN "),
        SYM_ON: []byte(" ON "),
        SYM_WHERE: []byte(" WHERE "),
        SYM_GROUP_BY: []byte(" GROUP BY "),
        SYM_ORDER_BY: []byte(" ORDER BY "),
        SYM_DESC: []byte(" DESC"),
        SYM_LIMIT: []byte(" LIMIT "),
        SYM_OFFSET: []byte(" OFFSET "),
        SYM_LPAREN: []byte("("),
        SYM_RPAREN: []byte(")"),
        SYM_IN: []byte(" IN ("),
        SYM_AND: []byte(" AND "),
        SYM_OR: []byte(" OR "),
        SYM_EQUAL: []byte(" = "),
        SYM_NEQUAL: []byte(" != "),
        SYM_BETWEEN: []byte(" BETWEEN "),
        SYM_MAX: []byte("MAX("),
        SYM_MIN: []byte("MIN("),
        SYM_SUM: []byte("SUM("),
        SYM_AVG: []byte("AVG("),
        SYM_COUNT_STAR: []byte("COUNT(*)"),
        SYM_COUNT_DISTINCT: []byte("COUNT(DISTINCT "),
    }
)
