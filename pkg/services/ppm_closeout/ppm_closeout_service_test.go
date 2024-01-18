package ppmcloseout

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMCloseoutSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *PPMCloseoutSuite) TestPPMCloseoutServiceSuite() {
	suite.Run("Able to return values from the DB for the PPM Closeout", func() {
		_, err := setUpMockPPMCloseout(suite)
		if err != nil {
			suite.NoError(err)
		}
	})
	suite.PopTestSuite.TearDown()
}

func setUpMockPPMCloseout(suite *PPMCloseoutSuite) (*models.PPMCloseout, error) {
	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.AppContextForTest().DB(), nil, nil)
	ppmCloseoutFetcher := NewPPMCloseoutFetcher()
	ppmCloseoutObj, err := ppmCloseoutFetcher.GetPPMCloseout(suite.AppContextForTest(), ppmShipment.ID)
	if err != nil {
		return nil, err
	}

	return ppmCloseoutObj, nil
}
