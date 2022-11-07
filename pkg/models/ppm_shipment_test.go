package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchPPMShipment() {
	t := suite.T()

	ppm := testdatagen.MakeMinimalPPMShipment(suite.DB(), testdatagen.Assertions{})

	_, err := suite.DB().ValidateAndSave(&ppm)
	if err != nil {
		t.Errorf("could not save PPM: %v", err)
		return
	}

	retrievedPPM, _ := models.FetchPPMShipmentFromMTOShipmentID(suite.DB(), ppm.ShipmentID)

	suite.Equal(retrievedPPM.ID, ppm.ID)
	suite.Equal(retrievedPPM.ShipmentID, ppm.ShipmentID)
}
