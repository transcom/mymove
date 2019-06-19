package services

// QueryFilter is an interface to allow passing filter values into query interfaces
// Ex `FetchMany` takes a list of filters
//go:generate mockery -name QueryFilter
type QueryFilter interface {
	Column() string
	Comparator() string
	Value() interface{}
}

// NewQueryFilter is a function type definition for building a QueryFilter
// Should allow handlers to parse query params into QueryFilters for services
type NewQueryFilter func(column string, comparator string, value interface{}) QueryFilter
