package services

// QueryOrder describes the "order by" clause in a sql query
//go:generate mockery --name QueryOrder --disable-version-string
type QueryOrder interface {
	Column() *string
	SortOrder() *bool
}

// NewQueryOrder describes a function that creates a new QueryOrder object
type NewQueryOrder func(column *string, sortOrder *bool) QueryOrder
