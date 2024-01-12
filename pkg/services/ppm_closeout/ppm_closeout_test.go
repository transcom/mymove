package ppmcloseout

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMCloseoutSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *PPMCloseoutSuite) TestPPMClouseoutServiceSuite(appCtx appcontext.AppContext) {
	suite.Run("Can return a PPM CloseOut object", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(appCtx.DB(), nil, nil)
		var ppmCloseout services.PPMCloseoutFetcher = &ppmCloseoutFetcher{}
		ppmCloseoutObj, err := ppmCloseout.GetPPMCloseout(appCtx, ppmShipment.ID)
		suite.Nil(err)
		suite.NotNil(ppmCloseoutObj)
		suite.NotNil(ppmCloseoutObj.ID)
	})
}
