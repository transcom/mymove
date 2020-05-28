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

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerSuccess() {
	ordersID, ppmID, user := setupShowPPMSitEstimateHandlerData(suite)

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
	mockedSitCharge := int64(55000)
	mockedCost := rateengine.CostComputation{}
	estimateCalculator := &mocks.EstimateCalculator{}
	estimateCalculator.On("CalculateEstimates",
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything, suite.TestLogger()).Return(mockedSitCharge, mockedCost, nil).Once()
	showEstimateHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger()), estimateCalculator}
	showResponse := showEstimateHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.IsType(&ppmop.ShowPPMSitEstimateOK{}, showResponse)
	okResponse := showResponse.(*ppmop.ShowPPMSitEstimateOK)
	sitCost := okResponse.Payload

	// And: Returned SIT cost to be as expected
	expectedSitCost := int64(55000)
	if *sitCost.Estimate != expectedSitCost {
		suite.T().Errorf("Expected move ppm SIT cost to be '%v', instead is '%v'", expectedSitCost, *sitCost.Estimate)
	}
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
func setupShowPPMSitEstimateHandlerData(suite *HandlerSuite) (orderID uuid.UUID, ppmID uuid.UUID, user models.ServiceMember) {
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

	user = testdatagen.MakeDefaultServiceMember(suite.DB())
	ordersID := uuid.Must(uuid.NewV4())
	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:               ordersID,
			NewDutyStationID: station.ID,
			ServiceMember:    user,
			ServiceMemberID:  user.ID,
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

	return ordersID, ppm.ID, user
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
