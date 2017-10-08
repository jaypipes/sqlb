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
    SYM_LEFT_JOIN
    SYM_CROSS_JOIN
    SYM_ON
    SYM_WHERE
    SYM_GROUP_BY
    SYM_ORDER_BY
    SYM_DESC
    SYM_LIMIT
    SYM_OFFSET
    SYM_INSERT
    SYM_DELETE
    SYM_VALUES
    SYM_LPAREN
    SYM_RPAREN
    SYM_IN
    SYM_AND
    SYM_OR
    SYM_EQUAL
    SYM_NEQUAL
    SYM_BETWEEN
    SYM_IS_NULL
    SYM_IS_NOT_NULL
    SYM_MAX
    SYM_MIN
    SYM_SUM
    SYM_AVG
    SYM_COUNT_STAR
    SYM_COUNT_DISTINCT
    SYM_CAST
    // TODO(jaypipes): Note that in PostreSQL, this is BTRIM(). We need a way
    // to translate to DB dialects based on a driver.
    SYM_TRIM
    SYM_CHAR_LENGTH
    SYM_BIT_LENGTH
    SYM_ASCII
    SYM_REVERSE
    SYM_CONCAT
    SYM_CONCAT_WS
    SYM_TYPE_CHAR
    SYM_TYPE_VARCHAR
    SYM_TYPE_BINARY
    SYM_TYPE_TEXT
    SYM_TYPE_INT
    SYM_TYPE_FLOAT
    SYM_TYPE_DECIMAL
    SYM_PLACEHOLDER = 9999999999
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
        SYM_LEFT_JOIN: []byte(" LEFT JOIN "),
        SYM_CROSS_JOIN: []byte(" CROSS JOIN "),
        SYM_ON: []byte(" ON "),
        SYM_WHERE: []byte(" WHERE "),
        SYM_GROUP_BY: []byte(" GROUP BY "),
        SYM_ORDER_BY: []byte(" ORDER BY "),
        SYM_DESC: []byte(" DESC"),
        SYM_LIMIT: []byte(" LIMIT "),
        SYM_OFFSET: []byte(" OFFSET "),
        SYM_INSERT: []byte("INSERT INTO "),
        SYM_DELETE: []byte("DELETE FROM "),
        SYM_VALUES: []byte(") VALUES ("),
        SYM_LPAREN: []byte("("),
        SYM_RPAREN: []byte(")"),
        SYM_IN: []byte(" IN ("),
        SYM_AND: []byte(" AND "),
        SYM_OR: []byte(" OR "),
        SYM_EQUAL: []byte(" = "),
        SYM_NEQUAL: []byte(" != "),
        SYM_BETWEEN: []byte(" BETWEEN "),
        SYM_IS_NULL: []byte(" IS NULL"),
        SYM_IS_NOT_NULL: []byte(" IS NOT NULL"),
        SYM_MAX: []byte("MAX("),
        SYM_MIN: []byte("MIN("),
        SYM_SUM: []byte("SUM("),
        SYM_AVG: []byte("AVG("),
        SYM_COUNT_STAR: []byte("COUNT(*)"),
        SYM_COUNT_DISTINCT: []byte("COUNT(DISTINCT "),
        SYM_CAST: []byte("CAST("),
        SYM_TRIM: []byte("TRIM("),
        SYM_CHAR_LENGTH: []byte("CHAR_LENGTH("),
        SYM_BIT_LENGTH: []byte("BIT_LENGTH("),
        SYM_ASCII: []byte("ASCII("),
        SYM_REVERSE: []byte("REVERSE("),
        SYM_CONCAT: []byte("CONCAT("),
        SYM_CONCAT_WS: []byte("CONCAT_WS("),
        SYM_TYPE_CHAR: []byte("CHAR"),
        SYM_TYPE_VARCHAR: []byte("VARCHAR"),
        SYM_TYPE_TEXT: []byte("TEXT"),
        SYM_TYPE_INT: []byte("INT"),
        SYM_TYPE_FLOAT: []byte("FLOAT"),
        SYM_TYPE_DECIMAL: []byte("DECIMAL"),
        SYM_TYPE_BINARY: []byte("BINARY"),
    }
)
