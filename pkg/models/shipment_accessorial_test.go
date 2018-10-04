package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchAccessorial() {
	//Setup
	accessorial := testdatagen.MakeDefaultShipmentAccessorial(suite.db)
	//make more items that don't relate to the first
	testdatagen.MakeDefaultShipmentAccessorial(suite.db)
	testdatagen.MakeDefaultShipmentAccessorial(suite.db)

	//Do
	accs, err := models.FetchAccessorialsByShipmentID(suite.db, &accessorial.ShipmentID)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(accessorial.ID, accs[0].ID)
}
