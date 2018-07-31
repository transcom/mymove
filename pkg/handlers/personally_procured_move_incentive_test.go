package handlers

import (
	"net/http/httptest"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestShowPPMIncentiveHandlerUnauthorised() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := testdatagen.MakeDefaultServiceMember(suite.db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.authenticateRequest(req, user)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          7500,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)

	_, succeeded := showResponse.(*ppmop.ShowPPMIncentiveUnauthorized)
	suite.True(succeeded, "ShowPpmIncentive allowed non-office user to call")

}

func (suite *HandlerSuite) TestShowPPMIncentiveHandler() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          7500,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.Equal(int64(637056), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(605203), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
func (suite *HandlerSuite) TestShowPPMIncentiveHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          600,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.Equal(int64(270252), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(256739), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
