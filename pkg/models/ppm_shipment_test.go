package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestDeletePPMShipment() {
	ppmShipment := testdatagen.MakeStubbedPPMShipment(suite.DB())

	err := DeletePPMShipment(suite.DB(), &ppmShipment)
	suite.NoError(err)
}
