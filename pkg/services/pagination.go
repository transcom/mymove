package services

type Pagination interface {
	Page() int
	PerPage() int
	Offset() int
}

type NewPagination func(page int64, perPage int64) Pagination
