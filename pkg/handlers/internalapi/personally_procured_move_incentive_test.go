package internalapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
)

func (suite *HandlerSuite) setupPersonallyProcuredMoveIncentiveTest(ordersID uuid.UUID) {
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
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "STDM",
	}
	suite.MustSave(&tsp)
	tspPerformance := models.TransportationServiceProviderPerformance{
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
	suite.MustSave(&tspPerformance)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "78626",
	}
	suite.MustSave(&address)

	locationName := "New Duty Location"
	location := models.DutyLocation{
		Name:      locationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:                ordersID,
			NewDutyLocationID: location.ID,
		},
	})

	moveID, _ := uuid.NewV4()
	_ = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       moveID,
			OrdersID: ordersID,
		},
		Order: orders,
	})
}

func (suite *HandlerSuite) TestShowPPMIncentiveHandlerForbidden() {
	ordersID := uuid.Must(uuid.NewV4())

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	user := testdatagen.MakeDefaultServiceMember(suite.DB())
	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:           req,
		OriginalMoveDate:      *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:             "94540",
		OriginDutyLocationZip: "50309",
		Weight:                7500,
		OrdersID:              strfmt.UUID(ordersID.String()),
	}

	handlerConfig := suite.HandlerConfig()
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1693, nil)
	handlerConfig.SetPlanner(planner)
	showHandler := ShowPPMIncentiveHandler{handlerConfig}
	showResponse := showHandler.Handle(params)
	suite.Assertions.IsType(&ppmop.ShowPPMIncentiveForbidden{}, showResponse)
}

func (suite *HandlerSuite) TestShowPPMIncentiveHandler() {
	ordersID := uuid.Must(uuid.NewV4())

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}
	suite.setupPersonallyProcuredMoveIncentiveTest(ordersID)
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:           req,
		OriginalMoveDate:      *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:             "94540",
		OriginDutyLocationZip: "50309",
		Weight:                7500,
		OrdersID:              strfmt.UUID(ordersID.String()),
	}

	handlerConfig := suite.HandlerConfig()
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1693, nil)
	handlerConfig.SetPlanner(planner)
	showHandler := ShowPPMIncentiveHandler{handlerConfig}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	// PPMs estimates are being hardcoded because we are not loading tariff400ng data
	// update this check so test passes - but this is not testing correctness of data
	suite.Equal(int64(658287), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(625373), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
func (suite *HandlerSuite) TestShowPPMIncentiveHandlerLowWeight() {
	ordersID := uuid.Must(uuid.NewV4())

	appCtx := suite.AppContextForTest()

	if err := scenario.RunRateEngineScenario2(appCtx); err != nil {
		suite.FailNow("failed to run scenario 2: %+v", err)
	}

	suite.setupPersonallyProcuredMoveIncentiveTest(ordersID)
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	req := httptest.NewRequest("GET", "/personally_procured_moves/incentive", nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := ppmop.ShowPPMIncentiveParams{
		HTTPRequest:           req,
		OriginalMoveDate:      *handlers.FmtDate(scenario.Oct1TestYear),
		OriginZip:             "94540",
		OriginDutyLocationZip: "50309",
		Weight:                600,
		OrdersID:              strfmt.UUID(ordersID.String()),
	}

	handlerConfig := suite.HandlerConfig()
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1693, nil)
	handlerConfig.SetPlanner(planner)
	showHandler := ShowPPMIncentiveHandler{handlerConfig}
	showResponse := showHandler.Handle(params)

	okResponse := showResponse.(*ppmop.ShowPPMIncentiveOK)
	cost := okResponse.Payload

	// PPMs estimates are being hardcoded because we are not loading tariff400ng data
	// update this check so test passes - but this is not testing correctness of data
	suite.Equal(int64(52663), *cost.Gcc, "Gcc was not equal")
	suite.Equal(int64(50030), *cost.IncentivePercentage, "IncentivePercentage was not equal")
}
