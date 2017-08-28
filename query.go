package sqlb

type Query struct {
    b []byte
    args []interface{}
    sel *selectClause
}

func (q *Query) String() string {
    size := q.sel.size()
    argc := q.sel.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.sel.scan(q.b, q.args)
    return string(q.b)
}

func (q *Query) StringArgs() (string, []interface{}) {
    size := q.sel.size()
    argc := q.sel.argCount()
    if len(q.args) != argc  {
        q.args = make([]interface{}, argc)
    }
    if len(q.b) != size {
        q.b = make([]byte, size)
    }
    q.sel.scan(q.b, q.args)
    return string(q.b), q.args
}
