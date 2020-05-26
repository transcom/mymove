package internalapi

import (
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithDcos() {
	ordersID, ppmID := setupShowPPMSitEstimateHandlerData(suite)
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:              req,
		OriginalMoveDate:         *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:            4,
		OriginZip:                "77901",
		OrdersID:                 strfmt.UUID(ordersID.String()),
		WeightEstimate:           3000,
		PersonallyProcuredMoveID: strfmt.UUID(ppmID.String()),
	}
	// And: ShowPPMSitEstimateHandler is queried
	// temp values
	mockedSitCharge := int64(3000)
	mockedCost := rateengine.CostComputation{}
	estimateCalculator := &mocks.EstimateCalculator{}
	estimateCalculator.On("CalculateEstimates",
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything, suite.TestLogger()).Return(&mockedSitCharge, mockedCost, nil).Once()
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger()), estimateCalculator}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(55376)
	if *sitCost.Estimate != expectedSitCost {
		suite.T().Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}

func (suite *HandlerSuite) TestShowPPMSitEstimateHandler2cos() {
	helperShowPPMSitEstimateHandler(suite, "2")
}

func helperCreateZip3AndServiceAreas(suite *HandlerSuite) (models.Tariff400ngServiceArea, models.Tariff400ngServiceArea) {
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "779", RateArea: "US68", BasepointCity: "Victoria", State: "TX", ServiceArea: "748", Region: "6"})
	suite.MustSave(&models.Tariff400ngZip3{Zip3: "674", Region: "5", BasepointCity: "Salina", State: "KS", RateArea: "US58", ServiceArea: "320"})

	originServiceArea := models.Tariff400ngServiceArea{
		Name:               "Victoria, TX",
		ServiceArea:        "748",
		ServicesSchedule:   3,
		LinehaulFactor:     39,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1402),
		SIT185BRateCents:   unit.Cents(53),
		SITPDSchedule:      3,
	}
	suite.MustSave(&originServiceArea)

	destServiceArea := models.Tariff400ngServiceArea{
		Name:               "Salina, KS",
		ServiceArea:        "320",
		ServicesSchedule:   2,
		LinehaulFactor:     43,
		ServiceChargeCents: 350,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
		SIT185ARateCents:   unit.Cents(1292),
		SIT185BRateCents:   unit.Cents(51),
		SITPDSchedule:      2,
	}
	suite.MustSave(&destServiceArea)

	return originServiceArea, destServiceArea
}
func setupShowPPMSitEstimateHandlerData(suite *HandlerSuite) (orderID uuid.UUID, ppmID uuid.UUID) {
	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "67401",
	}
	suite.MustSave(&address)

	stationName := "New Duty Station"
	station := models.DutyStation{
		Name:        stationName,
		Affiliation: internalmessages.AffiliationAIRFORCE,
		AddressID:   address.ID,
		Address:     address,
	}
	suite.MustSave(&station)

	ordersID := uuid.Must(uuid.NewV4())
	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:               ordersID,
			NewDutyStationID: station.ID,
		},
	})

	moveID, _ := uuid.NewV4()
	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       moveID,
			OrdersID: ordersID,
		},
		Order: orders,
	})

	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID: moveID,
			Move:   move,
			//PickupPostalCode: &pickupZip,
			//OriginalMoveDate: &moveDate,
			//WeightEstimate:   &weightEstimate,
		},
	})

	return ordersID, ppm.ID
}
func helperShowPPMSitEstimateHandler(suite *HandlerSuite, codeOfService string) {
	t := suite.T()

	// Given: a TDL, TSP and TSP performance with SITRate for relevant location and date
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US68",
			DestinationRegion: "5",
			CodeOfService:     codeOfService,
		},
	}) // Victoria, TX to Salina, KS
	tsp := testdatagen.MakeDefaultTSP(suite.DB())

	_, destServiceArea := helperCreateZip3AndServiceAreas(suite)

	itemRate210A := models.Tariff400ngItemRate{
		Code:               "210A",
		Schedule:           &destServiceArea.SITPDSchedule,
		WeightLbsLower:     2000,
		WeightLbsUpper:     4000,
		RateCents:          unit.Cents(57600),
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&itemRate210A)

	itemRate225A := models.Tariff400ngItemRate{
		Code:               "225A",
		Schedule:           &destServiceArea.ServicesSchedule,
		WeightLbsLower:     2000,
		WeightLbsUpper:     4000,
		RateCents:          unit.Cents(9900),
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&itemRate225A)

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          testdatagen.PerformancePeriodStart,
		PerformancePeriodEnd:            testdatagen.PerformancePeriodEnd,
		RateCycleStart:                  testdatagen.PeakRateCycleStart,
		RateCycleEnd:                    testdatagen.PeakRateCycleEnd,
		TrafficDistributionListID:       tdl.ID,
		TransportationServiceProviderID: tsp.ID,
		QualityBand:                     models.IntPointer(1),
		BestValueScore:                  90,
		LinehaulRate:                    unit.NewDiscountRateFromPercent(50.5),
		SITRate:                         unit.NewDiscountRateFromPercent(50),
	}
	suite.MustSave(&tspPerformance)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "67401",
	}
	suite.MustSave(&address)

	stationName := "New Duty Station"
	station := models.DutyStation{
		Name:        stationName,
		Affiliation: internalmessages.AffiliationAIRFORCE,
		AddressID:   address.ID,
		Address:     address,
	}
	suite.MustSave(&station)

	ordersID := uuid.Must(uuid.NewV4())
	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:               ordersID,
			NewDutyStationID: station.ID,
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

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:      req,
		OriginalMoveDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:    4,
		OriginZip:        "77901",
		OrdersID:         strfmt.UUID(ordersID.String()),
		WeightEstimate:   3000,
	}
	// And: ShowPPMSitEstimateHandler is queried
	// temp values
	mockedSitCharge := int64(3000)
	mockedCost := rateengine.CostComputation{}
	estimateCalculator := &mocks.EstimateCalculator{}
	estimateCalculator.On("CalculateEstimates",
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything, suite.TestLogger()).Return(&mockedSitCharge, mockedCost, nil).Once()
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger()), estimateCalculator}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(55376)
	if *sitCost.Estimate != expectedSitCost {
		t.Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
}

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithError() {
	orderID := uuid.Must(uuid.NewV4())

	// Given: A PPM Estimate request with all relevant records except TSP performance
	helperCreateZip3AndServiceAreas(suite)

	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains required auth values
	req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
	req = suite.AuthenticateRequest(req, user)

	params := ppmop.ShowPPMSitEstimateParams{
		HTTPRequest:      req,
		OriginalMoveDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		DaysInStorage:    4,
		OriginZip:        "77901",
		OrdersID:         strfmt.UUID(orderID.String()),
		WeightEstimate:   3000,
	}
	// And: ShowPPMSitEstimateHandler is queried
	// temp values
	mockedSitCharge := int64(3000)
	mockedCost := rateengine.CostComputation{}
	estimateCalculator := &mocks.EstimateCalculator{}
	estimateCalculator.On("CalculateEstimates",
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything, suite.TestLogger()).Return(&mockedSitCharge, mockedCost, nil).Once()
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger()), estimateCalculator}
	showResponse := showHandler.Handle(params)

	// Then: Expect bad request response
	suite.CheckResponseNotFound(showResponse)
}
