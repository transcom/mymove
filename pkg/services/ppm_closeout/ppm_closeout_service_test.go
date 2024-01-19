package ppmcloseout

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMCloseoutSuite struct {
	*testingsuite.PopTestSuite
}

func TestPPMCloseoutServiceSuite(t *testing.T) {
	ts := &PPMCloseoutSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

// func setUpMockPPMCloseout(suite *PPMCloseoutSuite) (*models.PPMCloseout, error) {
// 	mockedPlanner := &mocks.Planner{}
// 	ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.AppContextForTest().DB(), nil, nil)
// 	ppmCloseoutFetcher := NewPPMCloseoutFetcher(mockedPlanner)
// 	ppmCloseoutObj, err := ppmCloseoutFetcher.GetPPMCloseout(suite.AppContextForTest(), ppmShipment.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ppmCloseoutObj, nil
// }
