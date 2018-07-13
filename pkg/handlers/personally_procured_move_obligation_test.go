package handlers

import (
	"net/http/httptest"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

// func (suite *HandlerSuite) TestShowPPMObligationHandlerUnauthorised() {
// 	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
// 		suite.FailNow("failed to run scenario 2: %+v", err)
// 	}

// 	user := testdatagen.MakeDefaultServiceMember(suite.db)

// 	req := httptest.NewRequest("GET", "/personally_procured_moves/obligation", nil)
// 	req = suite.authenticateRequest(req, user)

// 	params := ppmop.ShowPPMObligationParams{
// 		HTTPRequest:     req,
// 		PlannedMoveDate: *fmtDate(scenario.May15_2018),
// 		OriginZip:       "94540",
// 		DestinationZip:  "78626",
// 		Weight:          7500,
// 	}

// 	context := NewHandlerContext(suite.db, suite.logger)
// 	context.SetPlanner(route.NewTestingPlanner(1693))
// 	showHandler := ShowPPMObligationHandler(context)
// 	showResponse := showHandler.Handle(params)

// 	response := showResponse.(*ppmop.ShowPPMObligationUnauthorized)
// 	suite.checkErrorResponse(response, http.StatusUnauthorized, "Unauthorized")

// }

func (suite *HandlerSuite) TestShowPPMObligationHandler() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/obligation", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMObligationParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          7500,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMObligationHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMObligationOK)
	cost := okResponse.Payload

	suite.Equal(int64(637056), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(605203), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
func (suite *HandlerSuite) TestShowPPMObligationHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.db); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	req := httptest.NewRequest("GET", "/personally_procured_moves/obligation", nil)
	req = suite.authenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMObligationParams{
		HTTPRequest:     req,
		PlannedMoveDate: *fmtDate(scenario.May15_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		Weight:          600,
	}

	context := NewHandlerContext(suite.db, suite.logger)
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMObligationHandler(context)
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMObligationOK)
	cost := okResponse.Payload

	suite.Equal(int64(270252), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(256739), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
