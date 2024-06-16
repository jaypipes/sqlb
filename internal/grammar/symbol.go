//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package grammar

type Symbol int
type ScanInfo []Symbol

const (
	SYM_ELEMENT = iota // Marker for an element that self-scans into the SQL buffer
	SYM_SPACE
	SYM_QUEST_MARK
	SYM_DOLLAR
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
	SYM_HAVING
	SYM_ORDER_BY
	SYM_DESC
	SYM_LIMIT
	SYM_OFFSET
	SYM_INSERT
	SYM_DELETE
	SYM_VALUES
	SYM_UPDATE
	SYM_SET
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
	SYM_GREATER
	SYM_GREATER_EQUAL
	SYM_LESS
	SYM_LESS_EQUAL
	SYM_MAX
	SYM_MIN
	SYM_SUM
	SYM_AVG
	SYM_COUNT_STAR
	SYM_COUNT_DISTINCT
	SYM_CAST
	SYM_BTRIM
	SYM_TRIM
	SYM_LTRIM
	SYM_RTRIM
	SYM_LEADING
	SYM_TRAILING
	SYM_BOTH
	SYM_CHAR_LENGTH
	SYM_BIT_LENGTH
	SYM_ASCII
	SYM_REVERSE
	SYM_CONCAT
	SYM_CONCAT_WS
	SYM_NOW
	SYM_CURRENT_TIMESTAMP
	SYM_CURRENT_TIME
	SYM_CURRENT_DATE
	SYM_EXTRACT
	SYM_TYPE_CHAR
	SYM_TYPE_VARCHAR
	SYM_TYPE_BINARY
	SYM_TYPE_TEXT
	SYM_TYPE_INT
	SYM_TYPE_FLOAT
	SYM_TYPE_DECIMAL
	SYM_UNIT_MICROSECOND
	SYM_UNIT_SECOND
	SYM_UNIT_MINUTE
	SYM_UNIT_HOUR
	SYM_UNIT_DAY
	SYM_UNIT_WEEK
	SYM_UNIT_MONTH
	SYM_UNIT_QUARTER
	SYM_UNIT_YEAR
	SYM_UNIT_SECOND_MICROSECOND
	SYM_UNIT_MINUTE_MICROSECOND
	SYM_UNIT_MINUTE_SECOND
	SYM_UNIT_HOUR_MICROSECOND
	SYM_UNIT_HOUR_SECOND
	SYM_UNIT_HOUR_MINUTE
	SYM_UNIT_DAY_MICROSECOND
	SYM_UNIT_DAY_SECOND
	SYM_UNIT_DAY_MINUTE
	SYM_UNIT_DAY_HOUR
	SYM_UNIT_YEAR_MONTH
	SYM_PLACEHOLDER = 9999999999
)

var (
	Symbols = map[Symbol][]byte{
		SYM_QUEST_MARK:              []byte("?"),
		SYM_SPACE:                   []byte(" "),
		SYM_DOLLAR:                  []byte("$"),
		SYM_PERIOD:                  []byte("."),
		SYM_AS:                      []byte(" AS "),
		SYM_COMMA_WS:                []byte(", "),
		SYM_SELECT:                  []byte("SELECT "),
		SYM_FROM:                    []byte("FROM "),
		SYM_JOIN:                    []byte("JOIN "),
		SYM_LEFT_JOIN:               []byte("LEFT JOIN "),
		SYM_CROSS_JOIN:              []byte("CROSS JOIN "),
		SYM_ON:                      []byte(" ON "),
		SYM_WHERE:                   []byte("WHERE "),
		SYM_GROUP_BY:                []byte("GROUP BY "),
		SYM_HAVING:                  []byte("HAVING "),
		SYM_ORDER_BY:                []byte("ORDER BY "),
		SYM_DESC:                    []byte(" DESC"),
		SYM_LIMIT:                   []byte("LIMIT "),
		SYM_OFFSET:                  []byte(" OFFSET "),
		SYM_INSERT:                  []byte("INSERT INTO "),
		SYM_VALUES:                  []byte(") VALUES ("),
		SYM_DELETE:                  []byte("DELETE FROM "),
		SYM_UPDATE:                  []byte("UPDATE "),
		SYM_SET:                     []byte(" SET "),
		SYM_LPAREN:                  []byte("("),
		SYM_RPAREN:                  []byte(")"),
		SYM_IN:                      []byte(" IN ("),
		SYM_AND:                     []byte(" AND "),
		SYM_OR:                      []byte(" OR "),
		SYM_EQUAL:                   []byte(" = "),
		SYM_NEQUAL:                  []byte(" != "),
		SYM_BETWEEN:                 []byte(" BETWEEN "),
		SYM_IS_NULL:                 []byte(" IS NULL"),
		SYM_IS_NOT_NULL:             []byte(" IS NOT NULL"),
		SYM_GREATER:                 []byte(" > "),
		SYM_GREATER_EQUAL:           []byte(" >= "),
		SYM_LESS:                    []byte(" < "),
		SYM_LESS_EQUAL:              []byte(" <= "),
		SYM_MAX:                     []byte("MAX("),
		SYM_MIN:                     []byte("MIN("),
		SYM_SUM:                     []byte("SUM("),
		SYM_AVG:                     []byte("AVG("),
		SYM_COUNT_STAR:              []byte("COUNT(*)"),
		SYM_COUNT_DISTINCT:          []byte("COUNT(DISTINCT "),
		SYM_CAST:                    []byte("CAST("),
		SYM_BTRIM:                   []byte("BTRIM("),
		SYM_TRIM:                    []byte("TRIM("),
		SYM_LTRIM:                   []byte("LTRIM("),
		SYM_RTRIM:                   []byte("RTRIM("),
		SYM_LEADING:                 []byte("LEADING"),
		SYM_TRAILING:                []byte("TRAILING"),
		SYM_BOTH:                    []byte("BOTH"),
		SYM_CHAR_LENGTH:             []byte("CHAR_LENGTH("),
		SYM_BIT_LENGTH:              []byte("BIT_LENGTH("),
		SYM_ASCII:                   []byte("ASCII("),
		SYM_REVERSE:                 []byte("REVERSE("),
		SYM_CONCAT:                  []byte("CONCAT("),
		SYM_CONCAT_WS:               []byte("CONCAT_WS("),
		SYM_NOW:                     []byte("NOW()"),
		SYM_CURRENT_TIMESTAMP:       []byte("CURRENT_TIMESTAMP()"),
		SYM_CURRENT_TIME:            []byte("CURRENT_TIME()"),
		SYM_CURRENT_DATE:            []byte("CURRENT_DATE()"),
		SYM_EXTRACT:                 []byte("EXTRACT("),
		SYM_TYPE_CHAR:               []byte("CHAR"),
		SYM_TYPE_VARCHAR:            []byte("VARCHAR"),
		SYM_TYPE_TEXT:               []byte("TEXT"),
		SYM_TYPE_INT:                []byte("INT"),
		SYM_TYPE_FLOAT:              []byte("FLOAT"),
		SYM_TYPE_DECIMAL:            []byte("DECIMAL"),
		SYM_TYPE_BINARY:             []byte("BINARY"),
		SYM_UNIT_MICROSECOND:        []byte("MICROSECOND"),
		SYM_UNIT_SECOND:             []byte("SECOND"),
		SYM_UNIT_MINUTE:             []byte("MINUTE"),
		SYM_UNIT_HOUR:               []byte("HOST"),
		SYM_UNIT_DAY:                []byte("DAY"),
		SYM_UNIT_WEEK:               []byte("WEEK"),
		SYM_UNIT_MONTH:              []byte("MONTH"),
		SYM_UNIT_QUARTER:            []byte("QUARTER"),
		SYM_UNIT_YEAR:               []byte("YEAR"),
		SYM_UNIT_SECOND_MICROSECOND: []byte("SECOND_MICROSECOND"),
		SYM_UNIT_MINUTE_MICROSECOND: []byte("MINUTE_MICROSECOND"),
		SYM_UNIT_MINUTE_SECOND:      []byte("MINUTE_SECOND"),
		SYM_UNIT_HOUR_MICROSECOND:   []byte("HOUR_MICROSECOND"),
		SYM_UNIT_HOUR_SECOND:        []byte("HOUR_SECOND"),
		SYM_UNIT_HOUR_MINUTE:        []byte("HOUR_MINUTE"),
		SYM_UNIT_DAY_MICROSECOND:    []byte("DAY_MICROSECOND"),
		SYM_UNIT_DAY_SECOND:         []byte("DAY_SECOND"),
		SYM_UNIT_DAY_MINUTE:         []byte("DAY_MINUTE"),
		SYM_UNIT_DAY_HOUR:           []byte("DAY_HOUR"),
		SYM_UNIT_YEAR_MONTH:         []byte("YEAR_MONTH"),
	}
)
