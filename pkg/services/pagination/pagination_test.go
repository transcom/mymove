package pagination

import "testing"

func (suite *PaginationServiceSuite) TestOffset() {
	suite.T().Run("should return the correct offset for a given page", func(t *testing.T) {
		page, perPage := int64(4), int64(25)
		pagination := NewPagination(&page, &perPage)

		suite.Equal(75, pagination.Offset())
	})
}
