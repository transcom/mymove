package internalapi

import (
	"net/http/httptest"

	"github.com/gofrs/uuid"
	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) TestShowPPMEstimateHandler() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *handlers.FmtDate(scenario.Oct1_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		WeightEstimate:  7500,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler{context}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(605203), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(668909), *cost.RangeMax, "RangeMax was not equal")
}

func (suite *HandlerSuite) TestShowPPMEstimateHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.MustSave(&user)

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.AuthenticateUserRequest(req, user)

	params := ppmop.ShowPPMEstimateParams{
		HTTPRequest:     req,
		PlannedMoveDate: *handlers.FmtDate(scenario.Oct1_2018),
		OriginZip:       "94540",
		DestinationZip:  "78626",
		WeightEstimate:  600,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler{context}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(256739), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(283765), *cost.RangeMax, "RangeMax was not equal")
}
