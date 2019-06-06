package query

import "github.com/transcom/mymove/pkg/services"

// QueryFilter contains the fields necessary to build a query filter clause
// Fields are private and methods are exposed to satisfy query building interfaces
type queryFilter struct {
	column     string
	comparator string
	value      interface{}
}

// Column returns the filter's column as a string
func (f queryFilter) Column() string {
	return f.column
}

// Comparator returns the filter's comparator as a string
func (f queryFilter) Comparator() string {
	return f.comparator
}

// Value returns the filter's value as a string
func (f queryFilter) Value() interface{} {
	return f.value
}

// NewQueryFilter is a builder for query filters to be used by handlers
// and talk to services that require query filters
func NewQueryFilter(column string, comparator string, value interface{}) services.QueryFilter {
	return queryFilter{
		column,
		comparator,
		value,
	}
}
