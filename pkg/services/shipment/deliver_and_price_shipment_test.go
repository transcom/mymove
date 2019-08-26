package shipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ShipmentServiceSuite) TestDeliverAndPriceShipment() {
	suite.T().Run("shipment is delivered", func(t *testing.T) {
		numTspUsers := 1
		numShipments := 1
		numShipmentOfferSplit := []int{1}
		status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
		_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
		suite.FatalNoError(err)

		shipment := shipments[0]

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
		verrs, err := NewShipmentDeliverAndPricer(
			suite.DB(),
			engine,
			route.NewTestingPlanner(1044),
		).DeliverAndPriceShipment(deliveryDate, &shipment)

		suite.FatalNoError(err)
		suite.FatalFalse(verrs.HasAny())

		suite.Equal(models.ShipmentStatusDELIVERED, shipment.Status)

		suite.DB().Reload(&shipment)

		fetchedLineItems, err := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)
		suite.FatalNoError(err)
		// All items should be priced
		for _, item := range fetchedLineItems {
			suite.NotNil(item.AmountCents, item.Tariff400ngItem.Code)
		}
	})

	suite.T().Run("transaction rolls back when deliver of shipment fails", func(t *testing.T) {
		numTspUsers := 1
		numShipments := 1
		numShipmentOfferSplit := []int{1}
		invalidTransitionStatus := []models.ShipmentStatus{models.ShipmentStatusAPPROVED}
		_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, invalidTransitionStatus, models.SelectedMoveTypeHHG)
		suite.FatalNoError(err)

		shipment := shipments[0]

		deliveryDate := testdatagen.DateInsidePerformancePeriod
		planner := route.NewTestingPlanner(1044)

		engine := rateengine.NewRateEngine(suite.DB(), suite.logger)
		verrs, err := NewShipmentDeliverAndPricer(
			suite.DB(),
			engine,
			planner,
		).DeliverAndPriceShipment(deliveryDate, &shipment)

		suite.Empty(verrs.Errors)
		suite.Error(err)

		suite.DB().Reload(&shipment)
		suite.Equal(models.ShipmentStatusAPPROVED, shipment.Status)

		// No items should be priced
		fetchedLineItems, _ := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)
		for _, item := range fetchedLineItems {
			suite.Nil(item.AmountCents, item.Tariff400ngItem.Code)
		}
	})

	suite.T().Run("transaction rolls back when pricing fails", func(t *testing.T) {
		numTspUsers := 1
		numShipments := 1
		numShipmentOfferSplit := []int{1}
		status := []models.ShipmentStatus{models.ShipmentStatusINTRANSIT}
		_, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
		suite.FatalNoError(err)

		shipment := shipments[0]
		shipment.MoveID = uuid.UUID{} // make shipment unprice-able to force error

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

		deliveryDate := testdatagen.DateInsidePerformancePeriod
		engine := rateengine.NewRateEngine(suite.DB(), suite.logger)

		verrs, err := NewShipmentDeliverAndPricer(
			suite.DB(),
			engine,
			route.NewTestingPlanner(1044),
		).DeliverAndPriceShipment(deliveryDate, &shipment)

		suite.NotEmpty(verrs)
		suite.Error(err)

		suite.DB().Reload(&shipment)
		suite.Equal(models.ShipmentStatusINTRANSIT, shipment.Status)

		// No items should be priced
		fetchedLineItems, _ := models.FetchLineItemsByShipmentID(suite.DB(), &shipment.ID)
		for _, item := range fetchedLineItems {
			suite.Nil(item.AmountCents, item.Tariff400ngItem.Code)
		}
	})
}
