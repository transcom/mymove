package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	movingexpenseservice "github.com/transcom/mymove/pkg/services/moving_expense"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestCreateMovingExpenseHandler() {
	// Reusable objects
	movingExpenseCreator := movingexpenseservice.NewMovingExpenseCreator()

	type movingExpenseCreateSubtestData struct {
		ppmShipment models.PPMShipment
		params      movingexpenseops.CreateMovingExpenseParams
		handler     CreateMovingExpenseHandler
	}

	makeCreateSubtestData := func(authenticateRequest bool, closeoutForCustomerFeatureFlag bool) (subtestData movingExpenseCreateSubtestData) {

		subtestData.ppmShipment = factory.BuildPPMShipment(suite.DB(), nil, nil)
		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense", subtestData.ppmShipment.ID.String())
		req := httptest.NewRequest("POST", endpoint, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}

		subtestData.params = movingexpenseops.CreateMovingExpenseParams{
			HTTPRequest:   req,
			PpmShipmentID: *handlers.FmtUUID(subtestData.ppmShipment.ID),
		}

		closeoutForCustomerFF := services.FeatureFlag{
			Key:   "complete_ppm_closeout_for_customer",
			Match: false,
		}

		handlerConfig := suite.HandlerConfig()
		if !closeoutForCustomerFeatureFlag {
			mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(closeoutForCustomerFF, nil)
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)
		}

		subtestData.handler = CreateMovingExpenseHandler{
			handlerConfig,
			movingExpenseCreator,
		}

		return subtestData
	}
	suite.Run("Successfully Create Moving Expense - Integration Test", func() {
		subtestData := makeCreateSubtestData(true, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseCreated{}, response)

		createdMovingExpense := response.(*movingexpenseops.CreateMovingExpenseCreated).Payload

		suite.NotEmpty(createdMovingExpense.ID.String())
		suite.Equal(createdMovingExpense.PpmShipmentID.String(), subtestData.ppmShipment.ID.String())
		suite.NotNil(createdMovingExpense.DocumentID.String())
	})

	suite.Run("POST failure - 400- bad request", func() {
		subtestData := makeCreateSubtestData(true, true)

		params := subtestData.params
		// Missing PPM Shipment ID
		params.PpmShipmentID = ""

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseBadRequest{}, response)
	})

	suite.Run("POST failure -401 - Unauthorized - unauthenticated user", func() {
		// user is unauthenticated to trigger 401
		subtestData := makeCreateSubtestData(false, true)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseUnauthorized{}, response)
	})

	suite.Run("POST failure - 422 - FF off", func() {
		subtestData := makeCreateSubtestData(true, false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.CreateMovingExpenseUnprocessableEntity{}, response)
		errResponse := response.(*movingexpenseops.CreateMovingExpenseUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Moving expenses cannot be created unless the complete_ppm_closeout_for_customer feature flag is enabled.")
	})

	suite.Run("Post failure - 500 - Server Error", func() {
		mockCreator := mocks.MovingExpenseCreator{}
		subtestData := makeCreateSubtestData(true, true)
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

func (suite *HandlerSuite) TestUpdateMovingExpenseHandlerUnit() {
	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)

		ppmShipment.MovingExpenses = append(
			ppmShipment.MovingExpenses,
			factory.BuildMovingExpense(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
			}, nil),
		)
	})

	setUpRequestAndParams := func() movingexpenseops.UpdateMovingExpenseParams {
		movingExpense := ppmShipment.MovingExpenses[0]

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		return movingexpenseops.UpdateMovingExpenseParams{
			HTTPRequest:         req,
			PpmShipmentID:       handlers.FmtUUIDValue(ppmShipment.ID),
			MovingExpenseID:     handlers.FmtUUIDValue(movingExpense.ID),
			IfMatch:             etag.GenerateEtag(movingExpense.UpdatedAt),
			UpdateMovingExpense: &ghcmessages.UpdateMovingExpense{},
		}
	}

	setUpMockMovingExpenseUpdater := func(returnValues ...interface{}) services.MovingExpenseUpdater {
		mockMovingExpenseUpdater := &mocks.MovingExpenseUpdater{}

		mockMovingExpenseUpdater.On(
			"UpdateMovingExpense",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.MovingExpense"),
			mock.AnythingOfType("string"),
		).Return(returnValues...)

		return mockMovingExpenseUpdater
	}

	setUpHandler := func(movingExpenseUpdater services.MovingExpenseUpdater) UpdateMovingExpenseHandler {
		return UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingExpenseUpdater,
		}
	}

	suite.Run("Returns an error if the request is not coming from the office app", func() {
		params := setUpRequestAndParams()

		params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, nil)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseForbidden{}, response)
	})

	suite.Run("Returns a NotFound response if the updater returns a NotFoundError", func() {
		params := setUpRequestAndParams()

		updateError := apperror.NewNotFoundError(ppmShipment.MovingExpenses[0].ID, "moving expense not found")
		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, updateError)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
	})

	suite.Run("Returns an InternalServerError response if the updater returns a QueryError", func() {
		params := setUpRequestAndParams()

		updateError := apperror.NewQueryError("MovingExpense", nil, "error getting moving expense")
		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, updateError)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseInternalServerError{}, response)
	})

	suite.Run("Returns a PreconditionFailed response if the updater returns a PreconditionFailedError", func() {
		params := setUpRequestAndParams()

		updateError := apperror.NewPreconditionFailedError(ppmShipment.MovingExpenses[0].ID, nil)
		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, updateError)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		if suite.IsType(&movingexpenseops.UpdateMovingExpensePreconditionFailed{}, response) {
			payload := response.(*movingexpenseops.UpdateMovingExpensePreconditionFailed).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.Equal(updateError.Error(), *payload.Message)
		}
	})

	suite.Run("Returns an UnprocessableEntity response if the updater returns an InvalidInputError", func() {
		params := setUpRequestAndParams()

		verrs := validate.NewErrors()
		fieldWithErr := "field"
		fieldErrorMsg := "Field error"
		verrs.Add(fieldWithErr, fieldErrorMsg)

		updateError := apperror.NewInvalidInputError(ppmShipment.MovingExpenses[0].ID, nil, verrs, "Invalid input")

		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, updateError)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		if suite.IsType(&movingexpenseops.UpdateMovingExpenseUnprocessableEntity{}, response) {
			payload := response.(*movingexpenseops.UpdateMovingExpenseUnprocessableEntity).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.Equal(handlers.ValidationErrMessage, *payload.Title)

			suite.Len(payload.InvalidFields, 1)

			fieldErrors, ok := payload.InvalidFields[fieldWithErr]
			suite.True(ok, "Expected field error to be present")

			suite.Contains(fieldErrors, fieldErrorMsg)
		}
	})

	suite.Run("Returns an InternalServerError response if the updater returns an unexpected error", func() {
		params := setUpRequestAndParams()

		updateError := apperror.NewNotImplementedError("Not implemented")
		movingExpenseUpdater := setUpMockMovingExpenseUpdater(nil, updateError)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		suite.IsType(&movingexpenseops.UpdateMovingExpenseInternalServerError{}, response)
	})

	suite.Run("Returns an updated moving expense if the updater succeeds", func() {
		params := setUpRequestAndParams()

		params.UpdateMovingExpense = &ghcmessages.UpdateMovingExpense{
			Status: ghcmessages.PPMDocumentStatusAPPROVED,
		}

		updatedMovingExpense := ppmShipment.MovingExpenses[0]
		docStatusApproved := models.PPMDocumentStatusApproved
		updatedMovingExpense.Status = &docStatusApproved

		movingExpenseUpdater := setUpMockMovingExpenseUpdater(&updatedMovingExpense, nil)

		handler := setUpHandler(movingExpenseUpdater)

		response := handler.Handle(params)

		if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
			payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.Equal(updatedMovingExpense.ID.String(), payload.ID.String())

			if suite.NotNil(payload.Status) {
				suite.Equal(string(*updatedMovingExpense.Status), string(*payload.Status))
			}
		}
	})
}

func (suite *HandlerSuite) TestUpdateMovingExpenseHandlerIntegration() {
	var userUploader *uploader.UserUploader
	var ppmShipment models.PPMShipment
	var movingExpense models.MovingExpense

	setupData := func() {
		var err error

		userUploader, err = uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), userUploader, nil)

		packingMaterialsType := models.MovingExpenseReceiptTypePackingMaterials
		movingExpense = factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
				LinkOnly: true,
			},
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
			{
				Model: models.UserUpload{},
				ExtendedParams: &factory.UserUploadExtendedParams{
					UserUploader: userUploader,
					AppContext:   suite.AppContextForTest(),
				},
			},
			{
				Model: models.MovingExpense{
					MovingExpenseType: &packingMaterialsType,
				},
			},
		}, nil)

		ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, movingExpense)
	}

	setUpRequestAndParams := func(movingExpense models.MovingExpense) movingexpenseops.UpdateMovingExpenseParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := factory.BuildOfficeUser(nil, nil, nil)

		req = suite.AuthenticateOfficeRequest(req, officeUser)

		// We're setting the current weight ticket values in the default payload because the handler can't handle
		// partial updates currently. Tests can update fields they want to test changes to as needed.
		params := movingexpenseops.UpdateMovingExpenseParams{
			HTTPRequest:     req,
			PpmShipmentID:   handlers.FmtUUIDValue(ppmShipment.ID),
			MovingExpenseID: handlers.FmtUUIDValue(movingExpense.ID),
			IfMatch:         etag.GenerateEtag(movingExpense.UpdatedAt),
			UpdateMovingExpense: &ghcmessages.UpdateMovingExpense{
				Amount: movingExpense.Amount.Int64(),
			},
		}

		if movingExpense.Status != nil {
			params.UpdateMovingExpense.Status = ghcmessages.PPMDocumentStatus(*movingExpense.Status)
		}

		if movingExpense.Reason != nil {
			params.UpdateMovingExpense.Reason = *movingExpense.Reason
		}

		if movingExpense.SITStartDate != nil {
			params.UpdateMovingExpense.SitStartDate = strfmt.Date(*movingExpense.SITStartDate)
		}

		if movingExpense.SITEndDate != nil {
			params.UpdateMovingExpense.SitEndDate = strfmt.Date(*movingExpense.SITEndDate)
		}

		return params
	}

	setUpHandler := func() UpdateMovingExpenseHandler {
		sitEstimatedCost := models.CentPointer(unit.Cents(62500))
		ppmEstimator := mocks.PPMEstimator{}
		ppmEstimator.
			On(
				"CalculatePPMSITEstimatedCost",
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("*models.PPMShipment"),
			).
			Return(sitEstimatedCost, nil)
		return UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingexpenseservice.NewOfficeMovingExpenseUpdater(&ppmEstimator),
		}
	}

	suite.Run("Success", func() {
		suite.Run("Can approve a moving expense", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.UpdateMovingExpense.Status = ghcmessages.PPMDocumentStatusAPPROVED

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(ppmShipment.MovingExpenses[0].ID.String(), payload.ID.String())

				if suite.NotNil(payload.Status) {
					suite.Equal(string(models.PPMDocumentStatusApproved), string(*payload.Status))
				}
			}
		})

		suite.Run("Can exclude a moving expense", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			reason := "Not a valid receipt"
			params.UpdateMovingExpense.Status = ghcmessages.PPMDocumentStatusEXCLUDED
			params.UpdateMovingExpense.Reason = reason

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(ppmShipment.MovingExpenses[0].ID.String(), payload.ID.String())

				if suite.NotNil(payload.Status) {
					suite.Equal(string(models.PPMDocumentStatusExcluded), string(*payload.Status))
				}

				if suite.NotNil(payload.Reason) {
					suite.Equal(reason, string(*payload.Reason))
				}
			}
		})

		suite.Run("Can reject a moving expense", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			reason := "Over budget!"
			params.UpdateMovingExpense.Status = ghcmessages.PPMDocumentStatusREJECTED
			params.UpdateMovingExpense.Reason = reason

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(ppmShipment.MovingExpenses[0].ID.String(), payload.ID.String())

				if suite.NotNil(payload.Status) {
					suite.Equal(string(models.PPMDocumentStatusRejected), string(*payload.Status))
				}

				if suite.NotNil(payload.Reason) {
					suite.Equal(reason, string(*payload.Reason))
				}
			}
		})

		suite.Run("Can update a non-storage moving expense", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			newAmount := movingExpense.Amount.AddCents(1000)

			params.UpdateMovingExpense.Amount = newAmount.Int64()

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(movingExpense.ID.String(), payload.ID.String())

				if suite.NotNil(payload.Amount) {
					suite.Equal(newAmount.Int64(), *payload.Amount)
				}

				suite.Nil(payload.SitStartDate)
				suite.Nil(payload.SitEndDate)
			}
		})

		suite.Run("Can update a storage moving expense", func() {
			setupData()
			storageExpenseType := models.MovingExpenseReceiptTypeStorage
			storageMovingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
				{
					Model:    ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
					LinkOnly: true,
				},
				{
					Model:    ppmShipment,
					LinkOnly: true,
				},
				{
					Model: models.UserUpload{},
					ExtendedParams: &factory.UserUploadExtendedParams{
						UserUploader: userUploader,
						AppContext:   suite.AppContextForTest(),
					},
				},
				{
					Model: models.MovingExpense{
						MovingExpenseType: &storageExpenseType,
					},
				},
			}, nil)

			ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, storageMovingExpense)

			params := setUpRequestAndParams(storageMovingExpense)
			sitLocation := ghcmessages.SITLocationTypeORIGIN
			weightStored := 2000

			newAmount := movingExpense.Amount.AddCents(1000)
			newSitStartDate := testdatagen.NextValidMoveDate
			newSitEndDate := testdatagen.NextValidMoveDate.AddDate(0, 0, 10)

			params.UpdateMovingExpense.Amount = newAmount.Int64()
			params.UpdateMovingExpense.SitStartDate = strfmt.Date(newSitStartDate)
			params.UpdateMovingExpense.SitEndDate = strfmt.Date(newSitEndDate)
			params.UpdateMovingExpense.WeightStored = int64(weightStored)
			params.UpdateMovingExpense.SitLocation = &sitLocation

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseOK{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseOK).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(storageMovingExpense.ID.String(), payload.ID.String())

				if suite.NotNil(payload.Amount) {
					suite.Equal(newAmount.Int64(), *payload.Amount)
				}

				if suite.NotNil(payload.SitStartDate) {
					suite.True(newSitStartDate.Equal(handlers.FmtDatePtrToPop(payload.SitStartDate)))
				}

				if suite.NotNil(payload.SitEndDate) {
					suite.True(newSitEndDate.Equal(handlers.FmtDatePtrToPop(payload.SitEndDate)))
				}
			}
		})
	})

	suite.Run("Failure", func() {
		suite.Run("Returns a Forbidden response if the request doesn't come from the office app", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

			handler := setUpHandler()

			response := handler.Handle(params)

			suite.IsType(&movingexpenseops.UpdateMovingExpenseForbidden{}, response)
		})

		suite.Run("Returns a NotFound response when the moving expense is not found", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.MovingExpenseID = handlers.FmtUUIDValue(uuid.Must(uuid.NewV4()))

			handler := setUpHandler()

			response := handler.Handle(params)

			suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
		})

		suite.Run("Returns a PreconditionFailed response when the eTag doesn't match the expected eTag", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.IfMatch = "wrong eTag"

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpensePreconditionFailed{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpensePreconditionFailed).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				expectedErr := apperror.NewPreconditionFailedError(ppmShipment.MovingExpenses[0].ID, nil)

				suite.Equal(expectedErr.Error(), *payload.Message)
			}
		})

		suite.Run("Returns an UnprocessableEntity response when the requested updates aren't valid", func() {
			setupData()
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.UpdateMovingExpense = &ghcmessages.UpdateMovingExpense{
				Amount: movingExpense.Amount.Int64(),
				Status: ghcmessages.PPMDocumentStatusREJECTED,
				Reason: "",
			}

			handler := setUpHandler()

			response := handler.Handle(params)

			if suite.IsType(&movingexpenseops.UpdateMovingExpenseUnprocessableEntity{}, response) {
				payload := response.(*movingexpenseops.UpdateMovingExpenseUnprocessableEntity).Payload

				suite.NoError(payload.Validate(strfmt.Default))

				suite.Equal(handlers.ValidationErrMessage, *payload.Title)

				suite.Len(payload.InvalidFields, 1)

				fieldErrors, ok := payload.InvalidFields["Reason"]
				suite.True(ok, "Expected field error to be present")

				suite.Contains(fieldErrors, "reason is mandatory if the status is Excluded or Rejected")
			}
		})
	})
}

func (suite *HandlerSuite) TestDeleteMovingExpenseHandler() {
	// Create Reusable objects
	movingExpenseDeleter := movingexpenseservice.NewMovingExpenseDeleter()

	type movingExpenseDeleteSubtestData struct {
		ppmShipment   models.PPMShipment
		movingExpense models.MovingExpense
		params        movingexpenseops.DeleteMovingExpenseParams
		handler       DeleteMovingExpenseHandler
		officeUser    models.OfficeUser
	}
	makeDeleteSubtestData := func(authenticateRequest bool, closeoutForCustomerFeatureFlag bool) (subtestData movingExpenseDeleteSubtestData) {
		// Fake data:
		subtestData.movingExpense = factory.BuildMovingExpense(suite.DB(), nil, nil)
		subtestData.ppmShipment = subtestData.movingExpense.PPMShipment
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		subtestData.officeUser = officeUser

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expenses/%s", subtestData.ppmShipment.ID.String(), subtestData.movingExpense.ID.String())
		req := httptest.NewRequest("DELETE", endpoint, nil)
		if authenticateRequest {
			req = suite.AuthenticateOfficeRequest(req, officeUser)
		}
		subtestData.params = movingexpenseops.DeleteMovingExpenseParams{
			HTTPRequest:     req,
			PpmShipmentID:   *handlers.FmtUUID(subtestData.ppmShipment.ID),
			MovingExpenseID: *handlers.FmtUUID(subtestData.movingExpense.ID),
		}

		closeoutForCustomerFF := services.FeatureFlag{
			Key:   "complete_ppm_closeout_for_customer",
			Match: false,
		}

		handlerConfig := suite.HandlerConfig()
		if !closeoutForCustomerFeatureFlag {
			mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
			mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
				mock.Anything,
				mock.AnythingOfType("*appcontext.appContext"),
				mock.AnythingOfType("string"),
				mock.Anything,
			).Return(closeoutForCustomerFF, nil)
			handlerConfig.SetFeatureFlagFetcher(mockFeatureFlagFetcher)

			subtestData.handler = DeleteMovingExpenseHandler{
				handlerConfig,
				movingExpenseDeleter,
			}
		} else {
			// Use createS3HandlerConfig for the HandlerConfig because we are required to upload a doc
			subtestData.handler = DeleteMovingExpenseHandler{
				suite.createS3HandlerConfig(),
				movingExpenseDeleter,
			}
		}

		return subtestData
	}

	suite.Run("Successfully Delete Moving Expense - Integration Test", func() {
		subtestData := makeDeleteSubtestData(true, true)

		params := subtestData.params
		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNoContent{}, response)
	})

	suite.Run("DELETE failure - 401 - permission denied - not authenticated", func() {
		subtestData := makeDeleteSubtestData(false, true)
		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseUnauthorized{}, response)
	})

	suite.Run("DELETE failure - 403 - permission denied - wrong application / user", func() {
		subtestData := makeDeleteSubtestData(false, true)

		serviceMember := subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateRequest(req, serviceMember)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseForbidden{}, response)
	})

	suite.Run("DELETE failure - 404 - not found - ppm shipment ID and moving expense ID don't match", func() {
		subtestData := makeDeleteSubtestData(false, true)
		officeUser := subtestData.officeUser

		otherPPMShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.ppmShipment.Shipment.MoveTaskOrder.Orders,
				LinkOnly: true,
			},
		}, nil)

		subtestData.params.PpmShipmentID = *handlers.FmtUUID(otherPPMShipment.ID)
		req := subtestData.params.HTTPRequest
		unauthorizedReq := suite.AuthenticateOfficeRequest(req, officeUser)
		unauthorizedParams := subtestData.params
		unauthorizedParams.HTTPRequest = unauthorizedReq

		response := subtestData.handler.Handle(unauthorizedParams)
		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("DELETE failure - 404 - not found", func() {
		subtestData := makeDeleteSubtestData(true, true)
		params := subtestData.params
		// Wrong ID provided
		uuidString := handlers.FmtUUID(testdatagen.ConvertUUIDStringToUUID("e392b01d-3b23-45a9-8f98-e4d5b03c8a93"))
		params.MovingExpenseID = *uuidString

		response := subtestData.handler.Handle(params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseNotFound{}, response)
	})

	suite.Run("POST failure - 422 - FF off", func() {
		subtestData := makeDeleteSubtestData(true, false)

		response := subtestData.handler.Handle(subtestData.params)

		suite.IsType(&movingexpenseops.DeleteMovingExpenseUnprocessableEntity{}, response)
		errResponse := response.(*movingexpenseops.DeleteMovingExpenseUnprocessableEntity)

		suite.Contains(*errResponse.Payload.Detail, "Moving expenses cannot be deleted unless the complete_ppm_closeout_for_customer feature flag is enabled.")
	})

	suite.Run("DELETE failure - 500 - server error", func() {
		mockDeleter := mocks.MovingExpenseDeleter{}
		subtestData := makeDeleteSubtestData(true, true)
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
