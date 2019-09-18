package pagination

import "github.com/transcom/mymove/pkg/services"

type pagination struct {
	page    int
	perPage int
}

func (p pagination) Page() int {
	return int(p.page)
}

func (p pagination) PerPage() int {
	return int(p.perPage)
}

func (p pagination) Offset() int {
	return int((p.Page() - 1) * p.PerPage())
}

func DefaultPage() int64 {
	return 1
}

func DefaultPerPage() int64 {
	return 25
}

func NewPagination(page *int64, perPage *int64) services.Pagination {
	if page == nil {
		return pagination{int(DefaultPage()), int(DefaultPerPage())}
	}

	pageValue, perPageValue := int(*page), int(*perPage)

	return pagination{pageValue, perPageValue}
}
