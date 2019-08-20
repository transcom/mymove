package invoice

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type testValues struct {
	name                string
	requiresPreapproval bool
	approved            models.ShipmentLineItemStatus
	invoiced            bool
	expectedCount       int
}

var tvs = []testValues{
	{
		"preapproval approved",
		true,
		models.ShipmentLineItemStatusAPPROVED,
		false,
		1,
	},
	{
		"preapproval not approved",
		true,
		models.ShipmentLineItemStatusSUBMITTED,
		false,
		0,
	},
	{
		"no preapproval",
		false,
		models.ShipmentLineItemStatusSUBMITTED,
		false,
		1,
	},
	{
		"already invoiced",
		false,
		models.ShipmentLineItemStatusSUBMITTED,
		true,
		0,
	},
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceForShipmentCall() {
	for _, tv := range tvs {
		suite.T().Run(tv.name, func(t *testing.T) {
			shipment := testdatagen.MakeDefaultShipment(suite.DB())
			lineItem := helperSetupLineItem(shipment, tv, suite.DB())
			suite.NotEqual(models.ShipmentLineItem{}.ID, lineItem.ID)

			f := FetchShipmentForInvoice{suite.DB()}
			actualShipment, err := f.Call(shipment.ID)
			suite.NoError(err)

			suite.Equal(tv.expectedCount, len(actualShipment.ShipmentLineItems))
			if tv.expectedCount != 0 {
				suite.Equal(lineItem.ID, actualShipment.ShipmentLineItems[0].ID)
			}
		})
	}

	suite.T().Run("multiple line items", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultShipment(suite.DB())
		for _, tv := range tvs {
			helperSetupLineItem(shipment, tv, suite.DB())
		}

		f := FetchShipmentForInvoice{suite.DB()}
		actualShipment, err := f.Call(shipment.ID)
		suite.NoError(err)

		suite.Equal(2, len(actualShipment.ShipmentLineItems))
	})

	suite.T().Run("tariff item association", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultShipment(suite.DB())
		tariffItem := testdatagen.MakeDefaultTariff400ngItem(suite.DB())
		suite.NotEqual(tariffItem.ID, models.Tariff400ngItem{}.ID)
		testdatagen.MakeCompleteShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment:        shipment,
				ShipmentID:      shipment.ID,
				Status:          models.ShipmentLineItemStatusAPPROVED,
				Tariff400ngItem: tariffItem,
			},
		})

		f := FetchShipmentForInvoice{suite.DB()}
		actualShipment, err := f.Call(shipment.ID)
		suite.NoError(err)

		suite.Equal(tariffItem.ID, actualShipment.ShipmentLineItems[0].Tariff400ngItem.ID)
	})
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceWith35AValid() {
	shipment := makeShipment(suite.DB())
	lineItem := makeLineItem35A(suite.DB(), shipment, true)

	f := FetchShipmentForInvoice{suite.DB()}
	actualShipment, err := f.Call(shipment.ID)
	suite.NoError(err)
	suite.Equal(1, len(actualShipment.ShipmentLineItems))
	suite.Equal(lineItem.ID, actualShipment.ShipmentLineItems[0].ID)
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceWith35AInvalid() {
	shipment := makeShipment(suite.DB())
	_ = makeLineItem35A(suite.DB(), shipment, false)

	f := FetchShipmentForInvoice{suite.DB()}
	actualShipment, err := f.Call(shipment.ID)
	suite.NoError(err)
	suite.Equal(0, len(actualShipment.ShipmentLineItems))
}

func (suite *InvoiceServiceSuite) TestFetchInvoiceWith35ALegacy() {
	shipment := makeShipment(suite.DB())
	lineItem := makeLineItem35ALegacy(suite.DB(), shipment)

	f := FetchShipmentForInvoice{suite.DB()}
	actualShipment, err := f.Call(shipment.ID)
	suite.NoError(err)
	suite.Equal(1, len(actualShipment.ShipmentLineItems))
	suite.Equal(lineItem.ID, actualShipment.ShipmentLineItems[0].ID)
}

func helperSetupLineItem(shipment models.Shipment, tv testValues, db *pop.Connection) models.ShipmentLineItem {
	assertions := testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:   shipment,
			ShipmentID: shipment.ID,
			Status:     tv.approved,
		},
		Tariff400ngItem: models.Tariff400ngItem{
			RequiresPreApproval: tv.requiresPreapproval,
		},
	}
	if tv.invoiced {
		invoice := testdatagen.MakeInvoice(db, testdatagen.Assertions{})
		assertions.ShipmentLineItem.InvoiceID = &invoice.ID
	}
	return testdatagen.MakeCompleteShipmentLineItem(db, assertions)
}

func makeShipment(DB *pop.Connection) models.Shipment {
	return testdatagen.MakeDefaultShipment(DB)
}

func makeLineItem35A(DB *pop.Connection, shipment models.Shipment, hasActAmt bool) models.ShipmentLineItem {
	acc35A := testdatagen.MakeTariff400ngItem(DB, testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			Code:                "35A",
			Item:                "Third Party Service",
			RequiresPreApproval: true,
		},
	})

	estAmt := unit.Cents(1234)
	assertions := testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:            shipment,
			ShipmentID:          shipment.ID,
			Status:              models.ShipmentLineItemStatusCONDITIONALLYAPPROVED,
			Tariff400ngItem:     acc35A,
			Location:            "ORIGIN",
			Description:         swag.String("This is a Description"),
			Reason:              swag.String("this is a reason"),
			EstimateAmountCents: &estAmt,
		},
	}

	if hasActAmt {
		actAmt := unit.Cents(1000)
		assertions.ShipmentLineItem.ActualAmountCents = &actAmt
		assertions.ShipmentLineItem.Status = models.ShipmentLineItemStatusAPPROVED
	}
	lineItem := testdatagen.MakeCompleteShipmentLineItem(DB, assertions)

	return lineItem
}

func makeLineItem35ALegacy(DB *pop.Connection, shipment models.Shipment) models.ShipmentLineItem {
	// shipment := testdatagen.MakeDefaultShipment(DB)

	acc35A := testdatagen.MakeTariff400ngItem(DB, testdatagen.Assertions{
		Tariff400ngItem: models.Tariff400ngItem{
			Code: "35A",
			Item: "Third Party Service",
		},
	})

	assertions := testdatagen.Assertions{
		ShipmentLineItem: models.ShipmentLineItem{
			Shipment:        shipment,
			ShipmentID:      shipment.ID,
			Status:          models.ShipmentLineItemStatusAPPROVED,
			Tariff400ngItem: acc35A,
		},
	}

	return testdatagen.MakeCompleteShipmentLineItem(DB, assertions)
}
