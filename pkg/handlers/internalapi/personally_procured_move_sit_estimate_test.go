package internalapi

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, nil).Once()
	showEstimateHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
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

func (suite *HandlerSuite) TestShowPPMSitEstimateHandlerWithError() {
	ordersID, ppmID, user := setupShowPPMSitEstimateHandlerData(suite)

	suite.T().Run("not found ppm ID fails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
		req = suite.AuthenticateRequest(req, user)

		notFoundID := uuid.FromStringOrNil("07b277d6-8ee5-4288-9e08-72449aa6643f")
		params := ppmop.ShowPPMSitEstimateParams{
			HTTPRequest:              req,
			OriginalMoveDate:         *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
			DaysInStorage:            4,
			OriginZip:                "77901",
			OrdersID:                 strfmt.UUID(notFoundID.String()),
			WeightEstimate:           3000,
			PersonallyProcuredMoveID: strfmt.UUID(ppmID.String()),
		}

		mockedSitCharge := int64(3000)
		mockedCost := rateengine.CostComputation{}
		estimateCalculator := &mocks.EstimateCalculator{}
		estimateCalculator.On("CalculateEstimates",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, nil).Once()
		showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
		showResponse := showHandler.Handle(params)

		suite.CheckResponseNotFound(showResponse)
	})

	suite.T().Run("not found orders ID fails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
		req = suite.AuthenticateRequest(req, user)

		notFoundID := uuid.FromStringOrNil("07b277d6-8ee5-4288-9e08-72449aa6643f")
		params := ppmop.ShowPPMSitEstimateParams{
			HTTPRequest:              req,
			OriginalMoveDate:         *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
			DaysInStorage:            4,
			OriginZip:                "77901",
			OrdersID:                 strfmt.UUID(ordersID.String()),
			WeightEstimate:           3000,
			PersonallyProcuredMoveID: strfmt.UUID(notFoundID.String()),
		}

		mockedSitCharge := int64(3000)
		mockedCost := rateengine.CostComputation{}
		estimateCalculator := &mocks.EstimateCalculator{}
		estimateCalculator.On("CalculateEstimates",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, nil).Once()
		showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
		showResponse := showHandler.Handle(params)

		suite.CheckResponseNotFound(showResponse)
	})

	suite.T().Run("missing original move date fails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
		req = suite.AuthenticateRequest(req, user)

		params := ppmop.ShowPPMSitEstimateParams{
			HTTPRequest:              req,
			DaysInStorage:            4,
			OriginZip:                "77901",
			OrdersID:                 strfmt.UUID(ordersID.String()),
			WeightEstimate:           3000,
			PersonallyProcuredMoveID: strfmt.UUID(ppmID.String()),
		}

		mockedSitCharge := int64(3000)
		mockedCost := rateengine.CostComputation{}
		estimateCalculator := &mocks.EstimateCalculator{}
		estimateCalculator.On("CalculateEstimates",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, nil).Once()
		showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
		showResponse := showHandler.Handle(params)

		suite.IsType(&ppmop.ShowPPMSitEstimateUnprocessableEntity{}, showResponse)
	})

	suite.T().Run("no estimated weight fails", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/estimates/ppm_sit", nil)
		req = suite.AuthenticateRequest(req, user)

		params := ppmop.ShowPPMSitEstimateParams{
			HTTPRequest:              req,
			OriginalMoveDate:         *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
			DaysInStorage:            4,
			OriginZip:                "77901",
			OrdersID:                 strfmt.UUID(ordersID.String()),
			PersonallyProcuredMoveID: strfmt.UUID(ppmID.String()),
		}

		mockedSitCharge := int64(3000)
		mockedCost := rateengine.CostComputation{}
		estimateCalculator := &mocks.EstimateCalculator{}
		estimateCalculator.On("CalculateEstimates",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, nil).Once()
		showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
		showResponse := showHandler.Handle(params)

		suite.IsType(&ppmop.ShowPPMSitEstimateUnprocessableEntity{}, showResponse)
	})
}

func (suite *HandlerSuite) TestShowPpmSitEstimateHandlerEstimateCalculationFails() {
	ordersID, ppmID, user := setupShowPPMSitEstimateHandlerData(suite)

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

	var mockedSitCharge int64
	mockedCost := rateengine.CostComputation{}
	estimateCalculator := &mocks.EstimateCalculator{}
	mockedError := errors.New("this is an error")
	estimateCalculator.On("CalculateEstimates",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*models.PersonallyProcuredMove"), mock.Anything).Return(mockedSitCharge, mockedCost, mockedError).Once()
	showHandler := ShowPPMSitEstimateHandler{handlers.NewHandlerContext(suite.DB(), suite.Logger()), estimateCalculator}
	showResponse := showHandler.Handle(params)

	suite.IsType(&handlers.ErrResponse{}, showResponse)
}

func setupShowPPMSitEstimateHandlerData(suite *HandlerSuite) (orderID uuid.UUID, ppmID uuid.UUID, user models.ServiceMember) {
	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "67401",
	}
	suite.MustSave(&address)

	locationName := "New Duty Location"
	location := models.DutyLocation{
		Name:      locationName,
		AddressID: address.ID,
		Address:   address,
	}
	suite.MustSave(&location)

	user = testdatagen.MakeDefaultServiceMember(suite.DB())
	ordersID := uuid.Must(uuid.NewV4())
	orders := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:                ordersID,
			NewDutyLocationID: location.ID,
			ServiceMember:     user,
			ServiceMemberID:   user.ID,
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
		},
	})

	return ordersID, ppm.ID, user
}
