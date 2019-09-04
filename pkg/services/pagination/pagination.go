package pagination

import "github.com/transcom/mymove/pkg/services"

type pagination struct {
	page    int64
	perPage int64
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

func NewPagination(page int64, perPage int64) services.Pagination {
	return pagination{page, perPage}
}
