package handlers

import (
	"net/http/httptest"
	"time"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestShowPPMEstimateHandler() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user, _ := testdatagen.MakeServiceMember(suite.db)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.authenticateRequest(req, user)

	date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
	params := ppmop.ShowPPMEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(date),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		WeightEstimate:  7500,
	}
	// And: show Queue is queried
	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler(context)
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	suite.Equal(int64(605204), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(668910), *cost.RangeMax, "RangeMax was not equal")
}
