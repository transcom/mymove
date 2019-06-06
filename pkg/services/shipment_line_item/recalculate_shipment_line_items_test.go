package shipmentlineitem

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"

	"time"
)

func (suite *ShipmentLineItemServiceSuite) TestRecalculateShipmentLineItems() {

	planner := route.NewTestingPlanner(1001)
	statuses := []models.ShipmentStatus{
		models.ShipmentStatusDELIVERED,
	}

	// Generate the data we'll need to do a recalculation.
	tspUsers, shipments, _, _ := testdatagen.CreateShipmentOfferData(suite.DB(), 1, 1, []int{1}, statuses, models.SelectedMoveTypeHHG)

	tspSession := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUsers[0].UserID,
		IDToken:         "fake token",
		OfficeUserID:    tspUsers[0].ID,
	}
	// Create a line item for us to use. The Assertion makes sure we're using the shipment we made above.
	// It also makes sure there isn't an invoice as that's a condition that will prevent a recalc.
	shipmentLineItem1 := testdatagen.MakeCompleteShipmentLineItem(suite.DB(),
		testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Invoice:    models.Invoice{},
				InvoiceID:  nil,
				Shipment:   shipments[0],
				ShipmentID: shipments[0].ID,
			},
		},
	)

	// Create date range
	recalculateRange := models.ShipmentRecalculate{
		ShipmentUpdatedAfter:  time.Date(1970, time.January, 01, 0, 0, 0, 0, time.UTC),
		ShipmentUpdatedBefore: time.Now(),
		Active:                true,
	}
	suite.MustCreate(suite.DB(), &recalculateRange)
	// Make sure we have fuel prices to use for the calc.
	testdatagen.MakeDefaultFuelEIADieselPrices(suite.DB())

	// Happy path
	recalculator := NewShipmentLineItemRecalculator(suite.DB(), suite.logger)
	_, err := recalculator.RecalculateShipmentLineItems(shipmentLineItem1.ShipmentID, &tspSession, planner)
	suite.NoError(err)

}
