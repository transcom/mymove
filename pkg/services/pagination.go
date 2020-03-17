package services

// Pagination represents an interface for pagination
type Pagination interface {
	Page() int
	PerPage() int
	Offset() int
}

// NewPagination creates a new Pagination interface
type NewPagination func(page *int64, perPage *int64) Pagination
