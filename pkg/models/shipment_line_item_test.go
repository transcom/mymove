package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchLineItem() {
	//Setup
	lineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	//make more items that don't relate to the first
	testdatagen.MakeDefaultShipmentLineItem(suite.db)
	testdatagen.MakeDefaultShipmentLineItem(suite.db)

	//Do
	accs, err := models.FetchLineItemsByShipmentID(suite.db, &lineItem.ShipmentID)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(lineItem.ID, accs[0].ID)
}
