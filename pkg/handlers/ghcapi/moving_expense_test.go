package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	movingexpenseops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	movingexpenseservice "github.com/transcom/mymove/pkg/services/moving_expense"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestUpdateMovingExpenseHandlerUnit() {
	var ppmShipment models.PPMShipment

	suite.PreloadData(func() {
		userUploader, err := uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		ppmShipment.MovingExpenses = append(
			ppmShipment.MovingExpenses,
			testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
				ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
				PPMShipment:   ppmShipment,
				UserUploader:  userUploader,
			}),
		)
	})

	setUpRequestAndParams := func() movingexpenseops.UpdateMovingExpenseParams {
		movingExpense := ppmShipment.MovingExpenses[0]

		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

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

	suite.PreloadData(func() {
		var err error

		userUploader, err = uploader.NewUserUploader(suite.createS3HandlerConfig().FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)

		suite.FatalNoError(err)

		ppmShipment = testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{
			UserUploader: userUploader,
		})

		packingMaterialsType := models.MovingExpenseReceiptTypePackingMaterials
		movingExpense = testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
			ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
			PPMShipment:   ppmShipment,
			UserUploader:  userUploader,
			MovingExpense: models.MovingExpense{
				MovingExpenseType: &packingMaterialsType,
			},
		})

		ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, movingExpense)
	})

	setUpRequestAndParams := func(movingExpense models.MovingExpense) movingexpenseops.UpdateMovingExpenseParams {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/moving-expense/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

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
		return UpdateMovingExpenseHandler{
			suite.createS3HandlerConfig(),
			movingexpenseservice.NewOfficeMovingExpenseUpdater(),
		}
	}

	suite.Run("Success", func() {
		suite.Run("Can approve a moving expense", func() {
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
			storageExpenseType := models.MovingExpenseReceiptTypeStorage
			storageMovingExpense := testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{
				ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
				PPMShipment:   ppmShipment,
				UserUploader:  userUploader,
				MovingExpense: models.MovingExpense{
					MovingExpenseType: &storageExpenseType,
				},
			})

			ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, storageMovingExpense)

			params := setUpRequestAndParams(storageMovingExpense)

			newAmount := movingExpense.Amount.AddCents(1000)
			newSitStartDate := testdatagen.NextValidMoveDate
			newSitEndDate := testdatagen.NextValidMoveDate.AddDate(0, 0, 10)

			params.UpdateMovingExpense.Amount = newAmount.Int64()
			params.UpdateMovingExpense.SitStartDate = strfmt.Date(newSitStartDate)
			params.UpdateMovingExpense.SitEndDate = strfmt.Date(newSitEndDate)

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
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.HTTPRequest = suite.AuthenticateRequest(params.HTTPRequest, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

			handler := setUpHandler()

			response := handler.Handle(params)

			suite.IsType(&movingexpenseops.UpdateMovingExpenseForbidden{}, response)
		})

		suite.Run("Returns a NotFound response when the moving expense is not found", func() {
			params := setUpRequestAndParams(ppmShipment.MovingExpenses[0])

			params.MovingExpenseID = handlers.FmtUUIDValue(uuid.Must(uuid.NewV4()))

			handler := setUpHandler()

			response := handler.Handle(params)

			suite.IsType(&movingexpenseops.UpdateMovingExpenseNotFound{}, response)
		})

		suite.Run("Returns a PreconditionFailed response when the eTag doesn't match the expected eTag", func() {
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
