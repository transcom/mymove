package handlers

import (
	"net/http/httptest"

	"github.com/gobuffalo/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestShowPPMEstimateHandler() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.authenticateRequest(req, user)

	params := ppmop.ShowPPMEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		WeightEstimate:  7500,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(605203), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(668909), *cost.RangeMax, "RangeMax was not equal")
}

func (suite *HandlerSuite) TestShowPPMEstimateHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.authenticateRequest(req, user)

	params := ppmop.ShowPPMEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		WeightEstimate:  600,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(256739), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(283765), *cost.RangeMax, "RangeMax was not equal")
}
