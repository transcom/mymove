package pagination

import "github.com/transcom/mymove/pkg/services"

type pagination struct {
	page    int
	perPage int
}

// Page represents the page number
func (p pagination) Page() int {
	return int(p.page)
}

// PerPage returns number per page
func (p pagination) PerPage() int {
	return int(p.perPage)
}

// Offset returns the offset in the pagination
func (p pagination) Offset() int {
	return int((p.Page() - 1) * p.PerPage())
}

// DefaultPage returns the default page
func DefaultPage() int64 {
	return 1
}

// DefaultPerPage returns the default per page
func DefaultPerPage() int64 {
	return 25
}

// NewPagination creates a new pagination object
func NewPagination(page *int64, perPage *int64) services.Pagination {
	if page == nil {
		return pagination{int(DefaultPage()), int(DefaultPerPage())}
	}

	pageValue, perPageValue := int(*page), int(*perPage)

	return pagination{pageValue, perPageValue}
}
