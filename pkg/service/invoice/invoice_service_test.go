package invoice

import (
	"log"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type InvoiceServiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *InvoiceServiceSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceServiceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}

func (suite *InvoiceServiceSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func helperShipment(suite *InvoiceServiceSuite) models.Shipment {
	return helperShipmentUsingScac(suite, "ABCD")
}

func helperShipmentUsingScac(suite *InvoiceServiceSuite, scac string) models.Shipment {
	var weight unit.Pound
	weight = 2000
	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			NetWeight: &weight,
		},
	})
	err := shipment.AssignGBLNumber(suite.db)
	suite.mustSave(&shipment)
	suite.NoError(err, "could not assign GBLNumber")

	// Create an accepted shipment offer and the associated TSP.
	supplierID := scac + "1234" //scac + payee code -- ABCD1234

	tsp := testdatagen.MakeTSP(suite.db, testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID,
		},
	})

	tspp := testdatagen.MakeTSPPerformance(suite.db, testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			TransportationServiceProvider:   tsp,
			TransportationServiceProviderID: tsp.ID,
		},
	})

	shipmentOffer := testdatagen.MakeShipmentOffer(suite.db, testdatagen.Assertions{
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
