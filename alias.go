package sqlb

type Aliasable interface {
    Alias(string)
}

func As(a Aliasable, alias string) Aliasable {
    a.Alias(alias)
    return a
}
