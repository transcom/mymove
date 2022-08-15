package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	movingexpenseservice "github.com/transcom/mymove/pkg/services/moving_expense"
	"github.com/transcom/mymove/pkg/testdatagen"
)

//
//CREATE TEST
//

// ADD Create test
func (suite *HandlerSuite) TestCreateMovingExpenseHandler() {
	// Reusable objects
	movingExpenseCreator := movingexpenseservice.NewMovingExpenseCreator()

	type movingExpenseCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      movingexpenseops.CreateMovingExpenseParams
		handler     CreateMovingExpenseHandler
	}

	makeCreateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData movingExpenseCreateSubtestData) {
		db := appCtx.DB()

		subtestData.ppmShipment = testdatagen.MakePPMShipment(db, testdatagen.Assertions{})
		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = movingexpenseops.CreateMovingExpenseParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		subtestData.handler = CreateMovingExpenseHandler{
			suite.HandlerConfig(),
			movingExpenseCreator,
		}

		return subtestData
	}
	suite.Run("Successfully Create Moving Expense - Integration Test", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseOK{}, response)

		createdMovingExpense := response.(*movingexpenseops.CreateMovingExpenseOK).Payload

		suite.NotEmpty(createdMovingExpense.ID.String())
		suite.NotNil(createdMovingExpense.PpmShipmentID.String())
		suite.NotNil(createdMovingExpense.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)

		params := subtestData.params
		// Missing PPM Shipment ID
		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		appCtx := suite.AppContextForTest()
		// user is unauthenticated to trigger 401
		subtestData := makeCreateSubtestData(appCtx, false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateWeightTicketUnauthorized{}, response)
	})

	suite.Run("POST failure - 403- permission denied - can't create moving expense due to wrong applicant", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)
		// Create non-service member user
		serviceCounselorOfficeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, serviceCounselorOfficeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.CreateMovingExpenseForbidden{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.MovingExpenseCreator{}
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, true)
		params := subtestData.params
		serverErr := errors.New("ServerError")

		// return a server error
		mockCreator.On("CreateMovingExpense",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(nil, serverErr)

		handler := CreateMovingExpenseHandler{
			suite.HandlerConfig(),
			&mockCreator,
		}

		response := handler.Handle(params)
		// Check the type to test the server error
		suite.IsType(&movingexpenseops.CreateMovingExpenseInternalServerError{}, response)
	})
}

// ADD Update test
