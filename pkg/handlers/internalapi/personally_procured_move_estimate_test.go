package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

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

	defaultPPM := testdatagen.MakeDefaultPPM(suite.DB())
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "296",
		},
	})

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.AuthenticateRequest(req, defaultPPM.Move.Orders.ServiceMember)

	params := ppmop.ShowPPMEstimateParams{
		PersonallyProcuredMoveID: strfmt.UUID(defaultPPM.ID.String()),
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:                "94540",
		DestinationZip:           "78626",
		WeightEstimate:           7500,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	//context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler{context}
	showHandler.SetPlanner(route.NewTestingPlanner(1693))
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(600359), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(663555), *cost.RangeMax, "RangeMax was not equal")
}

func (suite *HandlerSuite) TestShowPPMEstimateHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	defaultPPM := testdatagen.MakeDefaultPPM(suite.DB())
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "296",
		},
	})

	req := httptest.NewRequest("GET", "/estimates/ppm", nil)
	req = suite.AuthenticateRequest(req, defaultPPM.Move.Orders.ServiceMember)

	params := ppmop.ShowPPMEstimateParams{
		PersonallyProcuredMoveID: strfmt.UUID(defaultPPM.ID.String()),
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:                "94540",
		DestinationZip:           "78626",
		WeightEstimate:           600,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMEstimateHandler{context}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMEstimateOK)
	cost := okResponse.Payload

	suite.Equal(int64(256352), *cost.RangeMin, "RangeMin was not equal")
	suite.Equal(int64(283336), *cost.RangeMax, "RangeMax was not equal")
}
