package sqlb

type aliasable interface {
    setAlias(string)
}

type element interface {
    // Returns the number of bytes that the scannable element would consume as
    // a SQL string
    size() int
    // Returns the number of interface{} arguments that the element will add to
    // the slice of interface{} arguments passed to Scan()
    argCount() int
    // scan takes two slices -- one for the slice of bytes that the
    // implementation should copy its string representation to and another for
    // the slice of interface{} values that the element should add its
    // arguments to -- and returns two ints, one for the number of bytes that
    // it copied into the byte slice and another for the number of arguments
    // copied into the arg slice
    scan([]byte, []interface{}) (int, int)
}

// A projection is something that produces a scalar value. A column, column
// definition, function, etc. When appearing in the SELECT clause's projection
// list, the projection will output itself using the "AS alias" extended
// notation. When outputting in GROUP BY, ORDER BY or ON clauses, the
// projection will not include the alias extension
type projection interface {
    from() selection
    projectionId() uint64
    // projections must also implement element
    size() int
    argCount() int
    scan([]byte, []interface{}) (int, int)
    // disables the outputting of the "AS alias" extended output. Returns a
    // function that resets the outputting of the "AS alias" extended output
    disableAliasScan() func()
}

// A selection is something that produces rows. A table, table definition,
// view, subselect, etc
type selection interface {
    projections() []projection
    selectionId() uint64
    // selections must also implement element
    size() int
    argCount() int
    scan([]byte, []interface{}) (int, int)
}
