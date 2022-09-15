package services

// QueryFilter is an interface to allow passing filter values into query interfaces
// Ex `FetchMany` takes a list of filters
//go:generate mockery --name QueryFilter --disable-version-string
type QueryFilter interface {
	Column() string
	Comparator() string
	Value() interface{}
}

// NewQueryFilter is a function type definition for building a QueryFilter
// Should allow handlers to parse query params into QueryFilters for services
type NewQueryFilter func(column string, comparator string, value interface{}) QueryFilter

// QueryAssociation is an interface to allow passing association values into query interfaces
//go:generate mockery --name QueryAssociation --disable-version-string
type QueryAssociation interface {
	Field() string
}

// NewQueryAssociation is a function type definition for building a QueryAssociation
// Should allow services to pass associated data values for querying into QueryAssociations
type NewQueryAssociation func(field string) QueryAssociation

// QueryAssociations is an interface to allow
//go:generate mockery --name QueryAssociations --disable-version-string
type QueryAssociations interface {
	StringGetAssociations() []string
	Preload() bool
}

// NewQueryAssociations is a function type definition for building a QueryAssociations
// Should allow services to pass in query associations
type NewQueryAssociations func(associations []QueryAssociation) QueryAssociations
