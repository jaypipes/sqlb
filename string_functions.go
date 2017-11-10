package sqlb

type trimFunc struct {
	sel     selection
	alias   string
	subject element
	dialect Dialect
}

// Sets the sqlFunc's dialect and pushes the dialect down into any of the
// sqlFunc's elements
func (f *trimFunc) setDialect(dialect Dialect) {
	f.dialect = dialect
	switch f.subject.(type) {
	case *value:
		v := f.subject.(*value)
		v.dialect = dialect
	}
}

func (f *trimFunc) from() selection {
	return f.sel
}

func (f *trimFunc) disableAliasScan() func() {
	origAlias := f.alias
	f.alias = ""
	return func() { f.alias = origAlias }
}

func (f *trimFunc) As(alias string) *trimFunc {
	aliased := &trimFunc{
		sel:     f.sel,
		alias:   alias,
		subject: f.subject,
		dialect: f.dialect,
	}
	return aliased
}

func (f *trimFunc) argCount() int {
	return f.subject.argCount()
}

func (f *trimFunc) size() int {
	size := len(Symbols[SYM_TRIM]) + len(Symbols[SYM_RPAREN])
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = articles.author"
	switch f.subject.(type) {
	case projection:
		reset := f.subject.(projection).disableAliasScan()
		defer reset()
	}
	size += f.subject.size()
	if f.alias != "" {
		size += len(Symbols[SYM_AS]) + len(f.alias)
	}
	return size
}

func (f *trimFunc) scan(b []byte, args []interface{}, curArg *int) int {
	// TODO(jaypipes): Handle dialect differences
	bw := copy(b, Symbols[SYM_TRIM])
	// We need to disable alias output for elements that are
	// projections. We don't want to output, for example,
	// "ON users.id AS user_id = articles.author"
	switch f.subject.(type) {
	case projection:
		reset := f.subject.(projection).disableAliasScan()
		defer reset()
	}
	bw += f.subject.scan(b[bw:], args, curArg)
	bw += copy(b[bw:], Symbols[SYM_RPAREN])
	if f.alias != "" {
		bw += copy(b[bw:], Symbols[SYM_AS])
		bw += copy(b[bw:], f.alias)
	}
	return bw
}

func Trim(p projection) *trimFunc {
	return &trimFunc{
		subject: p.(element),
		sel:     p.from(),
	}
}

func (c *Column) Trim() *trimFunc {
	f := Trim(c)
	f.setDialect(c.tbl.meta.dialect)
	return f
}
