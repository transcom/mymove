package serviceitem

func (suite *ServiceItemServiceSuite) TestFetchServiceItem() {
	serviceItemFetcher := NewServiceItemFetcher()

	suite.Run("check that service items is not empty", func() {
		result, err := serviceItemFetcher.FetchServiceItemList(suite.AppContextForTest())
		suite.NoError(err)
		suite.NotEmpty(result)
	})
}
