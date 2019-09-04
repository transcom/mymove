package pagination

import "testing"

func (suite *PaginationServiceSuite) TestOffset() {
	suite.T().Run("should return the correct offset for a given page", func(t *testing.T) {
		pagination := NewPagination(4, 25)

		suite.Equal(75, pagination.Offset())
	})
}
