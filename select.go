package sqlb

type Selectable struct {
    alias string
    projected *ColumnList
    subjects []Scannable
}

func (s *Selectable) Alias(alias string) {
    s.alias = alias
}

func (s *Selectable) As(alias string) *Selectable {
    s.Alias(alias)
    return s
}

func (s *Selectable) Size() int {
    size := SYM_SELECT_LEN + SYM_FROM_LEN
    size += s.projected.Size()
    for _, subj := range s.subjects {
        size += subj.Size()
    }
    if s.alias != "" {
        size += SYM_AS_LEN + len(s.alias)
    }
    return size
}

func (s *Selectable) Scan(b []byte) int {
    idx := 0
    idx += copy(b[idx:], SYM_SELECT)
    idx += s.projected.Scan(b[idx:])
    idx += copy(b[idx:], SYM_FROM)
    for _, subj := range s.subjects {
        idx += subj.Scan(b[idx:])
    }
    if s.alias != "" {
        idx += copy(b[idx:], SYM_AS)
        idx += copy(b[idx:], s.alias)
    }
    return idx
}

func (s *Selectable) String() string {
    size := s.Size()
    b := make([]byte, size)
    s.Scan(b)
    return string(b)
}

func Select(items ...Scannable) *Selectable {
    // TODO(jaypipes): Make the memory allocation more efficient below by
    // looping through the elements and determining the number of Column struct
    // pointers to allocate instead of just making an empty array of Column
    // pointers.
    res := &Selectable{
        projected: &ColumnList{},
    }

    subjSet := make(map[Scannable]bool, 0)

    // For each scannable item we're received in the call, check what concrete
    // type they are and, depending on which type they are, either add them to
    // the returned Selectable's projected ColumnList or query the underlying
    // table metadata to generate a list of all columns in that table.
    for _, item := range items {
        switch item.(type) {
            case *Column:
                v := item.(*Column)
                res.projected.columns = append(res.projected.columns, v)
                subjSet[v.def.table] = true
            case *ColumnList:
                v := item.(*ColumnList)
                for _, c := range v.Columns() {
                    res.projected.columns = append(res.projected.columns, c)
                    subjSet[c.def.table] = true
                }
            case *Table:
                v := item.(*Table)
                for _, c := range v.Columns() {
                    res.projected.columns = append(res.projected.columns, c)
                }
                subjSet[v] = true
            case *TableDef:
                v := item.(*TableDef)
                for _, cd := range v.ColumnDefs() {
                    c := &Column{def: cd}
                    res.projected.columns = append(res.projected.columns, c)
                }
                t := &Table{def: v}
                subjSet[t] = true
            case *ColumnDef:
                v := item.(*ColumnDef)
                c := &Column{def: v}
                res.projected.columns = append(res.projected.columns, c)
                subjSet[v.table] = true
        }
    }
    subjects := make([]Scannable, len(subjSet))
    x := 0
    for scannable, _ := range subjSet {
        subjects[x] = scannable
        x++
    }
    res.subjects = subjects
    return res
}
