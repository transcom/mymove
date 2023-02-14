package internalapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
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
		suite.Equal(createdMovingExpense.PpmShipmentID.String(), subtestData.ppmShipment.ID.String())
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

		suite.IsType(&movingexpenseops.CreateMovingExpenseUnauthorized{}, response)
	})

	suite.Run("POST failure - 403- permission denied - can't create moving expense due to wrong applicant", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeCreateSubtestData(appCtx, false)
		// Create non-service member user
		serviceCounselorOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

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

//
// UPDATE test
//

func (suite *HandlerSuite) TestUpdateMovingExpenseHandler() {
	// Create Reusable objects
	movingExpenseUpdater := movingexpenseservice.NewCustomerMovingExpenseUpdater()

	type movingExpenseUpdateSubtestData struct {
		ppmShipment   models.PPMShipment
		movingExpense models.MovingExpense
		params        movingexpenseops.UpdateMovingExpenseParams
		handler       UpdateMovingExpenseHandler
	}
	makeUpdateSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData movingExpenseUpdateSubtestData) {
		db := appCtx.DB()

		// Fake data:
		subtestData.movingExpense = testdatagen.MakeMovingExpense(db, testdatagen.Assertions{})
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
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)

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
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		subtestData.params.UpdateMovingExpense = nil
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseBadRequest{}, response)
	})

	suite.Run("PATCH failure - 401- permission denied - not authenticated", func() {
		appCtx := suite.AppContextForTest()
		subtestData := makeUpdateSubtestData(appCtx, false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseUnauthorized{}, response)
	})

	suite.Run("PATCH failure - 403- permission denied - wrong application / user", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, false)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseForbidden{}, response)
	})

	suite.Run("PATCH failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{}
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.MovingExpenseID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{}
		params.IfMatch = "intentionally-bad-if-match-header-value"

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpensePreconditionFailed{}, response)
	})

	suite.Run("PATCH failure -422 - Invalid Input", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
		params := subtestData.params
		params.UpdateMovingExpense = &internalmessages.UpdateMovingExpense{
			Amount: handlers.FmtInt64(0),
		}

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 500 - server error", func() {
		mockUpdater := mocks.MovingExpenseUpdater{}
		appCtx := suite.AppContextForTest()

		subtestData := makeUpdateSubtestData(appCtx, true)
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
	makeDeleteSubtestData := func(appCtx appcontext.AppContext, authenticateRequest bool) (subtestData movingExpenseDeleteSubtestData) {
		db := appCtx.DB()

		// Fake data:
		subtestData.movingExpense = testdatagen.MakeMovingExpense(db, testdatagen.Assertions{})
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
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		appCtx := suite.AppContextForTest()
		subtestData := makeDeleteSubtestData(appCtx, false)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, false)

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseForbidden{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong service member user", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, false)

		otherServiceMember := testdatagen.MakeDefaultServiceMember(suite.DB())

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, otherServiceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseForbidden{}, response)
	})
	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and moving expense ID don't match", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, false)
		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		otherPPMShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			Order: models.Order{ServiceMemberID: serviceMember.ID},
		})

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, serviceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)
		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("DELETE failure - 404- not found", func() {
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.MovingExpenseID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.MovingExpenseDeleter{}
		appCtx := suite.AppContextForTest()

		subtestData := makeDeleteSubtestData(appCtx, true)
		params := subtestData.params

		err := errors.New("ServerError")

		mockDeleter.On("DeleteMovingExpense",
			mock.AnythingOfType("*appcontext.appContext"),
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
