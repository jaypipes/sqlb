package sqlb

func As(a aliasable, alias string) aliasable {
    a.setAlias(alias)
    return a
}
