package models_test

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestFetchGovBillOfLadingExtractor() {
	shipment := testdatagen.MakeDefaultShipment(suite.db)
	serviceAgent := testdatagen.MakeServiceAgent(suite.db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			ShipmentID: shipment.ID,
			Shipment:   shipment,
		},
	})
}
