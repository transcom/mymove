package serviceitem

import "github.com/transcom/mymove/pkg/models"

func (suite *ServiceItemServiceSuite) TestFetchServiceItem() {
	serviceItemFetcher := NewServiceItemFetcher()

	suite.Run("check that service items is not empty", func() {
		result, err := serviceItemFetcher.FetchServiceItemList(suite.AppContextForTest())
		suite.NoError(err)
		suite.NotEmpty(result)
	})

	suite.Run("POEFSC is auto approved for HHG/International shipments", func() {

		result, _ := serviceItemFetcher.FetchServiceItemList(suite.AppContextForTest())
		var poefscServiceItem models.ReServiceItem
		for _, v := range result {
			if models.ReServiceCodePOEFSC == v.ReService.Code && models.MTOShipmentTypeHHG == v.ShipmentType && models.MarketCodeInternational == v.MarketCode {
				poefscServiceItem = v
				break
			}
		}
		suite.Equal(true, poefscServiceItem.IsAutoApproved)
	})

	suite.Run("PODFSC is auto approved for UB/International shipments", func() {

		result, _ := serviceItemFetcher.FetchServiceItemList(suite.AppContextForTest())
		var podfscServiceItem models.ReServiceItem
		for _, v := range result {
			if models.ReServiceCodePODFSC == v.ReService.Code && models.MTOShipmentTypeUnaccompaniedBaggage == v.ShipmentType && models.MarketCodeInternational == v.MarketCode {
				podfscServiceItem = v
				break
			}
		}
		suite.Equal(true, podfscServiceItem.IsAutoApproved)
	})

	suite.Run("DOFSIT is NOT auto approved for UB/International shipments", func() {

		result, _ := serviceItemFetcher.FetchServiceItemList(suite.AppContextForTest())
		var dofsitServiceItem models.ReServiceItem
		for _, v := range result {
			if models.ReServiceCodeDOFSIT == v.ReService.Code && models.MTOShipmentTypeUnaccompaniedBaggage == v.ShipmentType && models.MarketCodeInternational == v.MarketCode {
				dofsitServiceItem = v
				break
			}
		}
		suite.Equal(false, dofsitServiceItem.IsAutoApproved)
	})
}
