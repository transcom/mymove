package shipmentlineitem

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ShipmentLineItemServiceSuite) TestGetShipmentLineItems() {
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	//officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	tspSession := auth.Session{
		ApplicationName: auth.TspApp,
		UserID: *tspUser.UserID,
		IDToken: "fake token",
		OfficeUserID: tspUser.ID,
	}

	shipmentLineItem1 := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{})
	_ = testdatagen.MakeCompleteShipmentLineItem(suite.DB(),
		testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment: shipmentLineItem1.Shipment,
			},
		})

	shipmentOfferAssertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			Shipment: shipmentLineItem1.Shipment,
		},
	}
	testdatagen.MakeShipmentOffer(suite.DB(), shipmentOfferAssertions)

	// Happy path. This should succeed.
	fetcher := NewShipmentLineItemFetcher(suite.DB())
	retrievedShipmentLineItems, err := fetcher.GetShipmentLineItemsByShipmentID(shipmentLineItem1.ShipmentID, &tspSession)
	suite.NoError(err)
	suite.Equal(2, len(retrievedShipmentLineItems))
	suite.Equal(shipmentLineItem1.ShipmentID, retrievedShipmentLineItems[0].ShipmentID)
	suite.Equal(shipmentLineItem1.ShipmentID, retrievedShipmentLineItems[1].ShipmentID)

}
