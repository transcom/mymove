package query

import "github.com/transcom/mymove/pkg/services"

// queryOrder contains the fields necessary to build a query order clause
// Fields are private and methods are exposed to satisfy query building interfaces
type queryOrder struct {
	column    *string
	sortOrder *bool
}

// Column returns the order column as a string
func (f queryOrder) Column() *string {
	return f.column
}

// SortOrder returns the sort ordering as a bool
// True, asc order
// False, desc order
func (f queryOrder) SortOrder() *bool {
	return f.sortOrder
}

// NewQueryOrder is a builder for query ordering to be used by handlers
// and talk to services that require query ordering
func NewQueryOrder(column *string, sortOrder *bool) services.QueryOrder {
	return queryOrder{
		column,
		sortOrder,
	}
}
