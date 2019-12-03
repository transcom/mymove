package services

// ListFetcher is the exported interface for fetching multiple records
//go:generate mockery -name ListFetcher
type ListFetcher interface {
	FetchRecordList(model interface{}, filters []QueryFilter, associations QueryAssociations, pagination Pagination, ordering QueryOrder) error
	FetchRecordCount(model interface{}, filters []QueryFilter) (int, error)
}
