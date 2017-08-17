package sqlb

type Element interface {
    // Returns the number of bytes that the scannable element would consume as
    // a SQL string
    Size() int
    // Returns the number of interface{} arguments that the element will add to
    // the slice of interface{} arguments passed to Scan()
    ArgCount() int
    // Scan takes two slices -- one for the slice of bytes that the
    // implementation should copy its string representation to and another for
    // the slice of interface{} values that the element should add its
    // arguments to -- and returns two ints, one for the number of bytes that
    // it copied into the byte slice and another for the number of arguments
    // copied into the arg slice
    Scan([]byte, []interface{}) (int, int)
}
