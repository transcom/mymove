package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	movingexpenseservice "github.com/transcom/mymove/pkg/services/moving_expense"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CREATE TEST
func (suite *HandlerSuite) TestCreateMovingExpenseHandler() {
	// Reusable objects
	movingExpenseCreator := movingexpenseservice.NewMovingExpenseCreator()

	type movingExpenseCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      movingexpenseops.CreateMovingExpenseParams
		handler     CreateMovingExpenseHandler
	}

	makeCreateSubtestData := func(authenticateRequest bool) (subtestData movingExpenseCreateSubtestData) {

		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
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
		subtestData := makeCreateSubtestData(true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseCreated{}, response)

		createdMovingExpense := response.(*movingexpenseops.CreateMovingExpenseCreated).Payload

		suite.NotEmpty(createdMovingExpense.ID.String())
		suite.Equal(createdMovingExpense.PpmShipmentID.String(), subtestData.ppmShipment.ID.String())
		suite.NotNil(createdMovingExpense.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true)

		params := subtestData.params
		// Missing PPM Shipment ID
		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		// user is unauthenticated to trigger 401
		subtestData := makeCreateSubtestData(false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseUnauthorized{}, response)
	})

	suite.Run("POST failure - 404 - Not Found - Wrong Service Member", func() {
		subtestData := makeCreateSubtestData(false)

		unauthorizedUser := factory.BuildServiceMember(suite.DB(), nil, nil)
		req := subtestData.params.HTTPRequest
		unauthorizedRequest := suite.AuthenticateRequest(req, unauthorizedUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedRequest

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.CreateMovingExpenseNotFound{}, response)
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.MovingExpenseCreator{}
		subtestData := makeCreateSubtestData(true)
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

//
// UPDATE test
//

func (suite *HandlerSuite) TestUpdateMovingExpenseHandler() {
	// Create Reusable objects
	ppmEstimator := mocks.PPMEstimator{}
	movingExpenseUpdater := movingexpenseservice.NewCustomerMovingExpenseUpdater(&ppmEstimator)

	type movingExpenseUpdateSubtestData struct {
		ppmShipment   models.PPMShipment
		movingExpense models.MovingExpense
		params        movingexpenseops.UpdateMovingExpenseParams
		handler       UpdateMovingExpenseHandler
	}
	makeUpdateSubtestData := func(authenticateRequest bool) (subtestData movingExpenseUpdateSubtestData) {
		// Fake data:
		subtestData.movingExpense = factory.BuildMovingExpense(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.movingExpense.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", subtestData.ppmShipment.ID.String(), subtestData.movingExpense.ID.String())
		req := httptest.NewRequest("PATCH", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		eTag := etag.GenerateEtag(subtestData.movingExpense.UpdatedAt)
		subtestData.params = movingexpenseops.
			UpdateMovingExpenseParams{
			HTTPRequest:     req,
			PpmShipmentID:   *handlers.FmtUUID(subtestData.ppmShipment.ID),
			MovingExpenseID: *handlers.FmtUUID(subtestData.movingExpense.ID),
			IfMatch:         eTag,
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingExpenseUpdater,
		}

		return subtestData
	}

	suite.Run("Successfully Update Moving Expense - Integration Test", func() {
		subtestData := makeUpdateSubtestData(true)

		params := subtestData.params
		// Add a Description
		contractedExpense := internalmessages.MovingExpenseType(models.MovingExpenseReceiptTypeContractedExpense)
		description := "Cost of moving items to a different location"
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{
			MovingExpenseType: &contractedExpense,
			Description:       &description,
		}
		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response)

		updatedMovingExpense := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload
		suite.Equal(subtestData.movingExpense.ID.String(), updatedMovingExpense.ID.String())
		suite.Equal(params.UpdateMovingExpense.Description, updatedMovingExpense.Description)
	})
	suite.Run("PATCH failure -400 - nil body", func() {
		subtestData := makeUpdateSubtestData(true)
		subtestData.params.UpdateMovingExpense = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseBadRequest{}, response)
	})

	suite.Run("PATCH failure - 401- permission denied - not authenticated", func() {
		subtestData := makeUpdateSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseUnauthorized{}, response)
	})

	suite.Run("PATCH failure - 403- permission denied - wrong application / user", func() {
		subtestData := makeUpdateSubtestData(false)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseForbidden{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{}
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.MovingExpenseID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
	})

	suite.Run("PATCH failure - 404- wrong service member", func() {
		subtestData := makeUpdateSubtestData(false)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{}
		params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, factory.BuildServiceMember(suite.DB(), nil, nil))

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpensePreconditionFailed{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{
			Amount: handlers.FmtInt64(0),
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 500 - server error", func() {
		mockUpdater := mocks.MovingExpenseUpdater{}
		subtestData := makeUpdateSubtestData(true)
		params := subtestData.params
		contractedExpense := internalmessages.MovingExpenseType(models.MovingExpenseReceiptTypeContractedExpense)
		description := "Cost of moving items to a different location"
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{
			MovingExpenseType: &contractedExpense,
			Description:       &description,
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateMovingExpense",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.MovingExpense"),
			mock.AnythingOfType("string"),
		).Return(nil, err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			&mockUpdater,
		}

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseInternalServerError{}, response)
		errResponse := response.(*movingexpenseops.UpdateMovingExpenseInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, *errResponse.Payload.Title, "Wrong Payload title")
	})
}

//
// DELETE test
//

func (suite *HandlerSuite) TestDeleteMovingExpenseHandler() {
	// Create Reusable objects
	movingExpenseDeleter := movingexpenseservice.NewMovingExpenseDeleter()

	type movingExpenseDeleteSubtestData struct {
		ppmShipment   models.PPMShipment
		movingExpense models.MovingExpense
		params        movingexpenseops.DeleteMovingExpenseParams
		handler       DeleteMovingExpenseHandler
	}
	makeDeleteSubtestData := func(authenticateRequest bool) (subtestData movingExpenseDeleteSubtestData) {
		// Fake data:
		subtestData.movingExpense = factory.BuildMovingExpense(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.movingExpense.PPMShipment
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expenses/%s", subtestData.ppmShipment.ID.String(), subtestData.movingExpense.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateRequest(req, serviceMember)
		}
		subtestData.params = movingexpenseops.
			DeleteMovingExpenseParams{
			HTTPRequest:     req,
			PpmShipmentID:   *handlers.FmtUUID(subtestData.ppmShipment.ID),
			MovingExpenseID: *handlers.FmtUUID(subtestData.movingExpense.ID),
		}

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		subtestData.handler = DeleteMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingExpenseDeleter,
		}

		return subtestData
	}

	suite.Run("Successfully Delete Moving Expense - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		subtestData := makeDeleteSubtestData(false)

		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseForbidden{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong service member user", func() {
		subtestData := makeDeleteSubtestData(false)

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, otherServiceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseForbidden{}, response)
	})
	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and moving expense ID don't match", func() {
		subtestData := makeDeleteSubtestData(false)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, serviceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)
		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.MovingExpenseID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.MovingExpenseDeleter{}
		subtestData := makeDeleteSubtestData(true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteMovingExpense",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
		).Return(err)

		// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
		handler := DeleteMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			&mockDeleter,
		}

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseInternalServerError{}, response)
	})
}
