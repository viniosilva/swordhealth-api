package repository

type WhereOpt interface {
	Query() string
	Values() []interface{}
}

type whereOpt struct {
	query  string
	values []interface{}
}

func (impl *whereOpt) Query() string {
	return impl.query
}

func (impl *whereOpt) Values() []interface{} {
	return impl.values
}

func SetWhere(query string, values []interface{}) WhereOpt {
	return &whereOpt{
		query:  query,
		values: values,
	}
}
