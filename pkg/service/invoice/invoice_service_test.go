package invoice

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type InvoiceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *InvoiceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func helperShipment(suite *InvoiceServiceSuite) models.Shipment {
	return helperShipmentUsingScac(suite, "ABBV")
}

func helperShipmentUsingScac(suite *InvoiceServiceSuite, scac string) models.Shipment {
	var weight unit.Pound
	weight = 2000
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			NetWeight: &weight,
		},
	})
	err := shipment.AssignGBLNumber(suite.DB())
	suite.MustSave(&shipment)
	suite.NoError(err, "could not assign GBLNumber")

	// Create an accepted shipment offer and the associated TSP.
	supplierID := scac + "2708" //scac + payee code -- ABBV2708

	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID,
		},
	})

	tspp := testdatagen.MakeTSPPerformance(suite.DB(), testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			TransportationServiceProvider:   tsp,
			TransportationServiceProviderID: tsp.ID,
		},
	})

	shipmentOffer := testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			Shipment:                                   shipment,
			Accepted:                                   swag.Bool(true),
			TransportationServiceProvider:              tsp,
			TransportationServiceProviderID:            tsp.ID,
			TransportationServiceProviderPerformance:   tspp,
			TransportationServiceProviderPerformanceID: tspp.ID,
		},
	})
	shipment.ShipmentOffers = models.ShipmentOffers{shipmentOffer}

	return shipment
}
