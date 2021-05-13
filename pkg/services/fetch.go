package services

// Fetcher is the exported interface for fetching a record
//go:generate mockery --name Fetcher
type Fetcher interface {
	FetchRecord(model interface{}, filters []QueryFilter) error
}

// ListFetcher is the exported interface for fetching multiple records
//go:generate mockery --name ListFetcher
type ListFetcher interface {
	FetchRecordList(model interface{}, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) error
	FetchRecordCount(model interface{}, filters []QueryFilter) (int, error)
}
