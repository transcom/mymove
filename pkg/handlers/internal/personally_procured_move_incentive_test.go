package internal

import (
	"net/http/httptest"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestShowPPMIncentiveHandlerForbidden() {
	if err := scenario.RunRateEngineScenario2(suite.parent.Db); err != nil {
		suite.parent.FailNow("failed to run scenario 2: %+v", err)
	}

	user := testdatagen.MakeDefaultServiceMember(suite.parent.Db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.parent.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *utils.FmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          7500,
	}

	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)
	suite.parent.Assertions.IsType(&ppmop.ShowPPMIncentiveForbidden{}, showResponse)
}

func (suite *HandlerSuite) TestShowPPMIncentiveHandler() {
	if err := scenario.RunRateEngineScenario2(suite.parent.Db); err != nil {
		suite.parent.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *utils.FmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          7500,
	}

	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.parent.Equal(int64(637056), *cost.Gcc, "Gcc was not equal")
	suite.parent.Equal(int64(605203), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
func (suite *HandlerSuite) TestShowPPMIncentiveHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.parent.Db); err != nil {
		suite.parent.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.parent.Db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.parent.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:     req,
		PlannedMoveDate: *utils.FmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          600,
	}

	context := utils.NewHandlerContext(suite.parent.Db, suite.parent.Logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.parent.Equal(int64(270252), *cost.Gcc, "Gcc was not equal")
	suite.parent.Equal(int64(256739), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
