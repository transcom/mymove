package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/go-openapi/strfmt"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) setupPersonallyProcuredMoveIncentiveTest() {
	originZip3 := models.Tariff400ngZip3{
		Zip3:          "503",
		BasepointCity: "Des Moines",
		State:         "IA",
		ServiceArea:   "296",
		RateArea:      "US53",
		Region:        "7",
	}
	suite.MustSave(&originZip3)
	destinationZip3 := models.Tariff400ngZip3{
		Zip3:          "956",
		BasepointCity: "Sacramento",
		State:         "CA",
		ServiceArea:   "68",
		RateArea:      "US87",
		Region:        "2",
	}
	suite.MustSave(&destinationZip3)
	destinationZip5 := models.Tariff400ngZip5RateArea{
		Zip5:     "95630",
		RateArea: "US87",
	}
	suite.MustSave(&destinationZip5)
	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Des Moines, IA",
		ServiceArea:        "296",
		LinehaulFactor:     unit.Cents(263),
		ServiceChargeCents: unit.Cents(489),
		ServicesSchedule:   3,
		EffectiveDateLower: scenario.May15TestYear,
		EffectiveDateUpper: scenario.May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1447),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)
	destinationServiceArea := models.Tariff400ngServiceArea{
		Name:               "Sacramento, CA",
		ServiceArea:        "68",
		LinehaulFactor:     unit.Cents(78),
		ServiceChargeCents: unit.Cents(452),
		ServicesSchedule:   3,
		EffectiveDateLower: scenario.May15TestYear,
		EffectiveDateUpper: scenario.May14FollowingYear,
		SIT185ARateCents:   unit.Cents(1642),
		SIT185BRateCents:   unit.Cents(70),
		SITPDSchedule:      3,
	}
	suite.MustSave(&destinationServiceArea)
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "US53",
		DestinationRegion: "6",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl1)
	tdl2 := models.TrafficDistributionList{
		SourceRateArea:    "US87",
		DestinationRegion: "2",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl2)
	tdl3 := models.TrafficDistributionList{
		SourceRateArea:    "US53",
		DestinationRegion: "2",
		CodeOfService:     "2",
	}
	suite.MustSave(&tdl3)
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	suite.MustSave(&tsp)
	tspPerformance1 := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          scenario.Oct1TestYear,
		PerformancePeriodEnd:            scenario.Dec31TestYear,
		RateCycleStart:                  scenario.Oct1TestYear,
		RateCycleEnd:                    scenario.May14FollowingYear,
		TrafficDistributionListID:       tdl1.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance1)
	tspPerformance2 := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          scenario.Oct1TestYear,
		PerformancePeriodEnd:            scenario.Dec31TestYear,
		RateCycleStart:                  scenario.Oct1TestYear,
		RateCycleEnd:                    scenario.May14FollowingYear,
		TrafficDistributionListID:       tdl2.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance2)
	tspPerformance3 := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          scenario.Oct1TestYear,
		PerformancePeriodEnd:            scenario.Dec31TestYear,
		RateCycleStart:                  scenario.Oct1TestYear,
		RateCycleEnd:                    scenario.May14FollowingYear,
		TrafficDistributionListID:       tdl3.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     swag.Int(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance3)
}

func (suite *HandlerSuite) TestShowPPMIncentiveHandlerForbidden() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	defaultPPM := testdatagen.MakeDefaultPPM(suite.DB())
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "296",
		},
	})

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateRequest(req, defaultPPM.Move.Orders.ServiceMember)

	params := ppmop.ShowPPMIncentiveParams{
		PersonallyProcuredMoveID: strfmt.UUID(defaultPPM.ID.String()),
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:                "94540",
		DestinationZip:           "78626",
		Weight:                   7500,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler{context}
	showResponse := showHandler.Handle(params)
	suite.Assertions.IsType(&ppmop.ShowPPMIncentiveForbidden{}, showResponse)
}

func (suite *HandlerSuite) TestShowPPMIncentiveHandler() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}
	suite.setupPersonallyProcuredMoveIncentiveTest()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	defaultPPM := testdatagen.MakeDefaultPPM(suite.DB())

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		PersonallyProcuredMoveID: strfmt.UUID(defaultPPM.ID.String()),
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:                "94540",
		DestinationZip:           "78626",
		Weight:                   7500,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler{context}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.Equal(int64(637056), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(605203), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
func (suite *HandlerSuite) TestShowPPMIncentiveHandlerLowWeight() {
	if err := scenario.RunRateEngineScenario2(suite.DB()); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	suite.setupPersonallyProcuredMoveIncentiveTest()
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	defaultPPM := testdatagen.MakeDefaultPPM(suite.DB())
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "296",
		},
	})

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		PersonallyProcuredMoveID: strfmt.UUID(defaultPPM.ID.String()),
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:                "94540",
		DestinationZip:           "78626",
		Weight:                   600,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(route.NewTestingPlanner(1693))
	showHandler := ShowPPMIncentiveHandler{context}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	suite.Equal(int64(270252), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(256739), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
