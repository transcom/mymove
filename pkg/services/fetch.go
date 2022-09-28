package services

import "github.com/transcom/mymove/pkg/appcontext"

// Fetcher is the exported interface for fetching a record
//
//go:generate mockery --name Fetcher --disable-version-string
type Fetcher interface {
	FetchRecord(appCtx appcontext.AppContext, model interface{}, filters []QueryFilter) error
}

// ListFetcher is the exported interface for fetching multiple records
//
//go:generate mockery --name ListFetcher --disable-version-string
type ListFetcher interface {
	FetchRecordList(appCtx appcontext.AppContext, model interface{}, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) error
	FetchRecordCount(appCtx appcontext.AppContext, model interface{}, filters []QueryFilter) (int, error)
}
