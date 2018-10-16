package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchGovBillOfLadingExtractor() {
	SourceTransOffice := testdatagen.MakeDefaultTransportationOffice(suite.db)
	DestinationTransOffice := testdatagen.MakeDefaultTransportationOffice(suite.db)

	packDate := time.Now()
	pickupDate := time.Now().AddDate(0, 0, 1)
	deliveryDate := time.Now().AddDate(0, 0, 2)
	edipi := "123456"
	gblNumber := "ABC12345"
	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			SourceGBLOC:                 &SourceTransOffice.Gbloc,
			DestinationGBLOC:            &DestinationTransOffice.Gbloc,
			PmSurveyPlannedDeliveryDate: &deliveryDate,
			PmSurveyPlannedPickupDate:   &pickupDate,
			PmSurveyPlannedPackDate:     &packDate,
			GBLNumber:                   &gblNumber,
		},
		ServiceMember: models.ServiceMember{
			Edipi: &edipi,
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("123"),
			SAC:                 models.StringPointer("456"),
			TAC:                 models.StringPointer("78901234"),
		},
	})
	testdatagen.MakeServiceAgent(suite.db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			ShipmentID: shipment.ID,
			Shipment:   &shipment,
		},
	})

	tsp := testdatagen.MakeDefaultTSP(suite.db)
	testdatagen.MakeShipmentOffer(suite.db, testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			ShipmentID:                      shipment.ID,
			Shipment:                        shipment,
			TransportationServiceProviderID: tsp.ID,
		},
	})

	gbl, err := models.FetchGovBillOfLadingExtractor(suite.db, shipment.ID)

	suite.NoError(err)

	suite.Equal(SourceTransOffice.Gbloc, gbl.IssuingOfficeGBLOC)
	suite.Equal(DestinationTransOffice.Gbloc, gbl.DestinationGbloc)
}
