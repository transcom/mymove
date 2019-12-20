package fetch

import (
	"github.com/transcom/mymove/pkg/services"
)

type fetcherQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type fetcher struct {
	builder fetcherQueryBuilder
}

// FetchRecord uses the passed query builder to fetch a record
func (o *fetcher) FetchRecord(model interface{}, filters []services.QueryFilter) error {
	error := o.builder.FetchOne(model, filters)
	return error
}

// NewFetcher returns an implementation of ListFetcher
func NewFetcher(builder fetcherQueryBuilder) services.Fetcher {
	return &fetcher{builder}
}
