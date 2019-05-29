package shipment

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ShipmentServiceSuite) TestDeliverPriceShipmentCall() {
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
	offerList, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.FatalNoError(err)

	shipment := shipments[0]

	authorizedStartDate := shipment.ActualPickupDate
	actualStartDate := authorizedStartDate.Add(testdatagen.OneDay)
	sit := testdatagen.MakeStorageInTransit(suite.DB(), testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			ShipmentID:          shipment.ID,
			Shipment:            shipment,
			EstimatedStartDate:  *authorizedStartDate,
			AuthorizedStartDate: authorizedStartDate,
			ActualStartDate:     &actualStartDate,
			Status:              models.StorageInTransitStatusINSIT,
		},
	})

	// And an unpriced, approved pre-approval
	testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     models.ShipmentLineItemStatusAPPROVED,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: true,
		},
	})

	// Make sure there's a FuelEIADieselPrice
	assertions := testdatagen.Assertions{}
	assertions.FuelEIADieselPrice.BaselineRate = 6
	testdatagen.MakeFuelEIADieselPrices(suite.DB(), assertions)

	deliveryDate := testdatagen.DateInsidePerformancePeriod
	engine := rateengine.NewRateEngine(suite.DB(), suite.logger)
	verrs, err := DeliverAndPriceShipment{
		DB:      suite.DB(),
		Engine:  engine,
		Planner: route.NewTestingPlanner(1044),
	}.Call(deliveryDate, &shipment, offerList[0].TransportationServiceProviderID)

	suite.FatalNoError(err)
	suite.FatalFalse(verrs.HasAny())

	suite.Equal(models.ShipmentStatusDELIVERED, shipment.Status)

	sits, _ := models.FetchStorageInTransitsOnShipment(suite.DB(), shipment.ID)
	sit = sits[0]

	suite.Equal(models.StorageInTransitStatusDELIVERED, sit.Status)
	suite.NotNil(sit.OutDate)

	updatedShipment, err := models.FetchShipmentByTSP(suite.DB(), offerList[0].TransportationServiceProviderID, shipment.ID)
	suite.Equal(updatedShipment.ActualDeliveryDate, sit.OutDate)

	fetchedLineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)
	suite.FatalNoError(err)
	// All items should be priced
	for _, item := range fetchedLineItems {
		suite.NotNil(item.AmountCents, item.Tariff400ngItem.Code)
	}
}
