package shipmentlineitem

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ShipmentLineItemServiceSuite) TestGetShipmentLineItems() {
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	serviceMemberUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	tspSession := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		OfficeUserID:    tspUser.ID,
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

	// When we don't have permission
	serviceMemberSession := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          serviceMemberUser.UserID,
		IDToken:         "fake token",
		ServiceMemberID: serviceMemberUser.ID,
	}

	_, err = fetcher.GetShipmentLineItemsByShipmentID(shipmentLineItem1.ShipmentID, &serviceMemberSession)
	suite.Equal(models.ErrFetchForbidden, err)

	assertions := testdatagen.Assertions{
		TspUser: models.TspUser{
			Email: "unused_test_email_for_sit@sit.com",
		},
	}

	tspUser2 := testdatagen.MakeTspUser(suite.DB(), assertions)
	tspSession1 := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser2.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser2.ID,
	}

	// TSP doesn't own the shipment
	shipmentLineItem := testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{})
	_, err = fetcher.GetShipmentLineItemsByShipmentID(shipmentLineItem.ShipmentID, &tspSession1)
	suite.Equal(models.ErrFetchForbidden, err)

}
