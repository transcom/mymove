package pagination

func (suite *PaginationServiceSuite) TestOffset() {
	suite.Run("should return the correct offset for a given page", func() {
		page, perPage := int64(4), int64(25)
		pagination := NewPagination(&page, &perPage)

		suite.Equal(75, pagination.Offset())
	})
}
