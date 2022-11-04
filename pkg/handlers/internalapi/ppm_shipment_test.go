package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	ppmops "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestSubmitPPMShipmentDocumentationHandlerUnit() {
	setUpPPMShipment := func() models.PPMShipment {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(
			suite.DB(),
			testdatagen.Assertions{
				Stub: true,
			},
		)

		ppmShipment.ID = uuid.Must(uuid.NewV4())
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt.AddDate(0, 0, 5)
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())

		return ppmShipment
	}

	setUpRequestAndParams := func(ppmShipment models.PPMShipment, authUser bool, setUpPayload bool) (*http.Request, ppmops.SubmitPPMShipmentDocumentationParams) {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())

		request := httptest.NewRequest("POST", endpoint, nil)

		if authUser {
			request = suite.AuthenticateRequest(request, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)
		}

		params := ppmops.SubmitPPMShipmentDocumentationParams{
			HTTPRequest:   request,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
			SavePPMShipmentSignedCertificationPayload: nil,
		}

		if setUpPayload {
			params.SavePPMShipmentSignedCertificationPayload = &internalmessages.SavePPMShipmentSignedCertification{
				CertificationText: handlers.FmtString("certification text"),
				Signature:         handlers.FmtString("signature"),
				Date:              handlers.FmtDate(time.Now()),
			}
		}

		return request, params
	}

	setUpPPMShipmentNewSubmitter := func(returnValue ...interface{}) services.PPMShipmentNewSubmitter {
		mockSubmitter := &mocks.PPMShipmentNewSubmitter{}

		mockSubmitter.On(
			"SubmitNewCustomerCloseOut",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockSubmitter
	}

	setUpHandler := func(submitter services.PPMShipmentNewSubmitter) SubmitPPMShipmentDocumentationHandler {
		return SubmitPPMShipmentDocumentationHandler{
			suite.HandlerConfig(),
			submitter,
		}
	}

	suite.Run("Returns an error if there is no session information", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, false, false)

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.SubmitPPMShipmentDocumentationUnauthorized{}, response)
	})

	suite.Run("Returns an error if the request isn't coming from the correct app", func() {
		ppmShipment := setUpPPMShipment()

		request, params := setUpRequestAndParams(ppmShipment, false, false)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params.HTTPRequest = request

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.SubmitPPMShipmentDocumentationForbidden{}, response)
	})

	suite.Run("Returns an error if the user ID is missing from the session", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.SubmitPPMShipmentDocumentationForbidden{}, response)
	})

	suite.Run("Returns an error if the PPMShipment ID in the url is invalid", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.ID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, "Invalid PPM shipment ID")
		}
	})

	suite.Run("Returns an error if there is no request body", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, "No body provided")
		}
	})

	suite.Run("Returns an error if the submitter service returns a BadDataError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewBadDataError("Bad data")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a NotFoundError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewNotFoundError(ppmShipment.ID, "Can't find PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a QueryError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewQueryError("PPMShipment", nil, "Error getting PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		suite.IsType(&ppmops.SubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns an error if the submitter service returns a InvalidInputError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		verrs := validate.NewErrors()
		fieldWithErr := "field"
		fieldErrorMsg := "Field error"
		verrs.Add(fieldWithErr, fieldErrorMsg)

		err := apperror.NewInvalidInputError(ppmShipment.ID, nil, verrs, "Invalid input")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationUnprocessableEntity{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationUnprocessableEntity)

			suite.Equal(handlers.ValidationErrMessage, *errResponse.Payload.Detail)

			fieldErrors, ok := errResponse.Payload.InvalidFields[fieldWithErr]
			suite.True(ok, "Expected field error to be present")
			suite.Contains(fieldErrors, fieldErrorMsg)
		}
	})

	suite.Run("Returns an error if the submitter service returns a ConflictError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewConflictError(ppmShipment.ID, "Can't route PPM shipment")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationConflict{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationConflict)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns an unexpected error", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewNotImplementedError("Not implemented")

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(nil, err))

		response := handler.Handle(params)

		suite.IsType(&ppmops.SubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns the PPM shipment if all goes well", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		expectedPPMShipment := ppmShipment
		expectedPPMShipment.Status = models.PPMShipmentStatusNeedsPaymentApproval
		expectedPPMShipment.SubmittedAt = models.TimePointer(time.Now())

		move := ppmShipment.Shipment.MoveTaskOrder
		certType := models.SignedCertificationTypePPMPAYMENT
		now := time.Now()
		signedCertification := models.SignedCertification{
			ID:                uuid.Must(uuid.NewV4()),
			SubmittingUserID:  move.Orders.ServiceMember.User.ID,
			MoveID:            move.ID,
			PpmID:             &ppmShipment.ID,
			CertificationType: &certType,
			CertificationText: *params.SavePPMShipmentSignedCertificationPayload.CertificationText,
			Signature:         *params.SavePPMShipmentSignedCertificationPayload.Signature,
			Date:              handlers.FmtDatePtrToPop(params.SavePPMShipmentSignedCertificationPayload.Date),
			CreatedAt:         now,
			UpdatedAt:         now,
		}

		expectedPPMShipment.SignedCertification = &signedCertification

		handler := setUpHandler(setUpPPMShipmentNewSubmitter(&expectedPPMShipment, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationOK{}, response) {
			okResponse := response.(*ppmops.SubmitPPMShipmentDocumentationOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(expectedPPMShipment.ID, returnedPPMShipment.ID)
			suite.EqualUUID(expectedPPMShipment.SignedCertification.ID, returnedPPMShipment.SignedCertification.ID)
		}
	})
}

func (suite *HandlerSuite) TestSubmitPPMShipmentDocumentationHandlerIntegration() {
	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	mtoShipmentRouter := mtoshipment.NewShipmentRouter()
	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(mtoShipmentRouter)

	submitter := ppmshipment.NewPPMShipmentNewSubmitter(signedCertificationCreator, ppmShipmentRouter)

	setUpParamsAndHandler := func(ppmShipment models.PPMShipment, payload *internalmessages.SavePPMShipmentSignedCertification) (ppmops.SubmitPPMShipmentDocumentationParams, SubmitPPMShipmentDocumentationHandler) {
		endpoint := fmt.Sprintf("/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())

		request := httptest.NewRequest("POST", endpoint, nil)

		request = suite.AuthenticateRequest(request, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		params := ppmops.SubmitPPMShipmentDocumentationParams{
			HTTPRequest:   request,
			PpmShipmentID: handlers.FmtUUIDValue(ppmShipment.ID),
			SavePPMShipmentSignedCertificationPayload: payload,
		}

		handler := SubmitPPMShipmentDocumentationHandler{
			suite.HandlerConfig(),
			submitter,
		}

		return params, handler
	}

	suite.Run("Returns an error if the PPM shipment is not found", func() {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		ppmShipment.ID = uuid.Must(uuid.NewV4())

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, "not found while looking for PPMShipment")
		}
	})

	suite.Run("Returns an error if the SignedCertification has any errors", func() {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString("certification text"),
			Signature:         handlers.FmtString("signature"),
		})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationUnprocessableEntity{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationUnprocessableEntity)

			suite.Equal(handlers.ValidationErrMessage, *errResponse.Payload.Detail)

			fieldErrors, ok := errResponse.Payload.InvalidFields["Date"]
			suite.True(ok, "Expected date error to be present")
			suite.Contains(fieldErrors, "Date is required")
		}
	})

	suite.Run("Returns an error if the PPM shipment is not in the right status", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{})

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString("certification text"),
			Signature:         handlers.FmtString("signature"),
			Date:              handlers.FmtDate(time.Now()),
		})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationConflict{}, response) {
			errResponse := response.(*ppmops.SubmitPPMShipmentDocumentationConflict)

			suite.Contains(
				*errResponse.Payload.Detail,
				fmt.Sprintf(
					"PPM shipment can't be set to %s because it's not in the %s status.",
					models.PPMShipmentStatusNeedsPaymentApproval,
					models.PPMShipmentStatusWaitingOnCustomer,
				),
			)
		}
	})

	suite.Run("Can successfully submit a PPM shipment for close out", func() {
		ppmShipment := testdatagen.MakePPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), testdatagen.Assertions{})

		certText := "certification text"
		signature := "signature"
		signDate := time.Now()

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString(certText),
			Signature:         handlers.FmtString(signature),
			Date:              handlers.FmtDate(signDate),
		})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.SubmitPPMShipmentDocumentationOK{}, response) {
			okResponse := response.(*ppmops.SubmitPPMShipmentDocumentationOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusNeedsPaymentApproval), string(returnedPPMShipment.Status))
			suite.NotNil(returnedPPMShipment.SubmittedAt)

			suite.NotNil(returnedPPMShipment.SignedCertification)
			suite.NotNil(returnedPPMShipment.SignedCertification.ID)

			suite.EqualUUID(
				ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.User.ID,
				returnedPPMShipment.SignedCertification.SubmittingUserID,
			)

			suite.EqualUUID(ppmShipment.Shipment.MoveTaskOrder.ID, returnedPPMShipment.SignedCertification.MoveID)

			if suite.NotNil(returnedPPMShipment.SignedCertification.PpmID) {
				suite.EqualUUID(ppmShipment.ID, *returnedPPMShipment.SignedCertification.PpmID)
			}

			suite.Equal(
				string(models.SignedCertificationTypePPMPAYMENT),
				string(returnedPPMShipment.SignedCertification.CertificationType),
			)

			if suite.NotNil(returnedPPMShipment.SignedCertification.CertificationText) {
				suite.Equal(certText, *returnedPPMShipment.SignedCertification.CertificationText)
			}

			if suite.NotNil(returnedPPMShipment.SignedCertification.Signature) {
				suite.Equal(signature, *returnedPPMShipment.SignedCertification.Signature)
			}

			suite.True(
				signDate.Equal(handlers.FmtDatePtrToPop(returnedPPMShipment.SignedCertification.Date)),
				"Expected sign dates to be equal",
			)
		}
	})
}

func (suite *HandlerSuite) TestResubmitPPMShipmentDocumentationHandlerUnit() {
	setUpPPMShipment := func() models.PPMShipment {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(
			suite.DB(),
			testdatagen.Assertions{
				Stub: true,
			},
		)

		ppmShipment.ID = uuid.Must(uuid.NewV4())
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt.AddDate(0, 0, 4)
		ppmShipment.SignedCertification.ID = uuid.Must(uuid.NewV4())

		ppmShipment.SignedCertification.CreatedAt = ppmShipment.UpdatedAt.AddDate(0, 0, -1)
		ppmShipment.SignedCertification.UpdatedAt = ppmShipment.UpdatedAt.AddDate(0, 0, 1)
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())

		return ppmShipment
	}

	setUpRequestAndParams := func(ppmShipment models.PPMShipment, authUser bool, setUpPayload bool) (*http.Request, ppmops.ResubmitPPMShipmentDocumentationParams) {
		endpoint := fmt.Sprintf(
			"/ppm-shipments/%s/submit-ppm-shipment-documentation/%s",
			ppmShipment.ID.String(),
			ppmShipment.SignedCertification.ID.String(),
		)

		request := httptest.NewRequest("PUT", endpoint, nil)

		if authUser {
			request = suite.AuthenticateRequest(request, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)
		}

		params := ppmops.ResubmitPPMShipmentDocumentationParams{
			HTTPRequest:           request,
			PpmShipmentID:         handlers.FmtUUIDValue(ppmShipment.ID),
			SignedCertificationID: handlers.FmtUUIDValue(ppmShipment.SignedCertification.ID),
			SavePPMShipmentSignedCertificationPayload: nil,
		}

		if setUpPayload {
			params.SavePPMShipmentSignedCertificationPayload = &internalmessages.SavePPMShipmentSignedCertification{
				CertificationText: handlers.FmtString("certification text"),
				Signature:         handlers.FmtString("signature"),
				Date:              handlers.FmtDate(time.Now()),
			}
		}

		return request, params
	}

	setUpPPMShipmentUpdatedSubmitter := func(returnValue ...interface{}) services.PPMShipmentUpdatedSubmitter {
		mockSubmitter := &mocks.PPMShipmentUpdatedSubmitter{}

		mockSubmitter.On(
			"SubmitUpdatedCustomerCloseOut",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockSubmitter
	}

	setUpHandler := func(submitter services.PPMShipmentUpdatedSubmitter) ResubmitPPMShipmentDocumentationHandler {
		return ResubmitPPMShipmentDocumentationHandler{
			suite.HandlerConfig(),
			submitter,
		}
	}

	suite.Run("Returns an error if there is no session information", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, false, false)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationUnauthorized{}, response)
	})

	suite.Run("Returns an error if the request isn't coming from the correct app", func() {
		ppmShipment := setUpPPMShipment()

		request, params := setUpRequestAndParams(ppmShipment, false, false)

		officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params.HTTPRequest = request

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationForbidden{}, response)
	})

	suite.Run("Returns an error if the user ID is missing from the session", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.UserID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationForbidden{}, response)
	})

	suite.Run("Returns an error if the PPMShipment ID in the url is invalid", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.ID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, "Invalid PPM shipment ID")
		}
	})

	suite.Run("Returns an error if the SignedCertification ID in the url is invalid", func() {
		ppmShipment := setUpPPMShipment()
		ppmShipment.SignedCertification.ID = uuid.Nil

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, "Invalid signed certification ID")
		}
	})

	suite.Run("Returns an error if there is no request body", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, false)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, "No body provided")
		}
	})

	suite.Run("Returns an error if the submitter service returns a BadDataError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewBadDataError("Bad data")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationBadRequest{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationBadRequest)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a NotFoundError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewNotFoundError(ppmShipment.ID, "Can't find PPM shipment")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a QueryError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewQueryError("PPMShipment", nil, "Error getting PPM shipment")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns an error if the submitter service returns a InvalidInputError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		verrs := validate.NewErrors()
		fieldWithErr := "field"
		fieldErrorMsg := "Field error"
		verrs.Add(fieldWithErr, fieldErrorMsg)

		err := apperror.NewInvalidInputError(ppmShipment.ID, nil, verrs, "Invalid input")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationUnprocessableEntity{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationUnprocessableEntity)

			suite.Equal(handlers.ValidationErrMessage, *errResponse.Payload.Detail)

			fieldErrors, ok := errResponse.Payload.InvalidFields[fieldWithErr]
			suite.True(ok, "Expected field error to be present")
			suite.Contains(fieldErrors, fieldErrorMsg)
		}
	})

	suite.Run("Returns an error if the submitter service returns a ConflictError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewConflictError(ppmShipment.ID, "Can't route PPM shipment")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationConflict{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationConflict)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns a PreconditionFailedError", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewPreconditionFailedError(ppmShipment.ID, nil)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationPreconditionFailed{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationPreconditionFailed)

			suite.Contains(*errResponse.Payload.Detail, err.Error())
		}
	})

	suite.Run("Returns an error if the submitter service returns an unexpected error", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		err := apperror.NewNotImplementedError("Not implemented")

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(nil, err))

		response := handler.Handle(params)

		suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationInternalServerError{}, response)
	})

	suite.Run("Returns the PPM shipment if all goes well", func() {
		ppmShipment := setUpPPMShipment()

		_, params := setUpRequestAndParams(ppmShipment, true, true)

		expectedPPMShipment := ppmShipment
		expectedPPMShipment.Status = models.PPMShipmentStatusNeedsPaymentApproval

		suite.FatalNotNil(expectedPPMShipment.SubmittedAt)

		originalSubmittedAt := expectedPPMShipment.SubmittedAt

		newSignDate := time.Now()
		expectedPPMShipment.SignedCertification.Date = newSignDate
		params.SavePPMShipmentSignedCertificationPayload.Date = handlers.FmtDate(newSignDate)

		handler := setUpHandler(setUpPPMShipmentUpdatedSubmitter(&expectedPPMShipment, nil))

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationOK{}, response) {
			okResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(expectedPPMShipment.ID, returnedPPMShipment.ID)
			suite.EqualUUID(expectedPPMShipment.SignedCertification.ID, returnedPPMShipment.SignedCertification.ID)

			if suite.NotNil(returnedPPMShipment.SubmittedAt) {
				suite.True(
					originalSubmittedAt.Equal(handlers.FmtDateTimePtrToPop(returnedPPMShipment.SubmittedAt)),
					"Expected SubmittedAt to be unchanged",
				)
			}
		}
	})
}

func (suite *HandlerSuite) TestResubmitPPMShipmentDocumentationHandlerIntegration() {
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	mtoShipmentRouter := mtoshipment.NewShipmentRouter()
	ppmShipmentRouter := ppmshipment.NewPPMShipmentRouter(mtoShipmentRouter)

	submitter := ppmshipment.NewPPMShipmentUpdatedSubmitter(signedCertificationUpdater, ppmShipmentRouter)

	setUpParamsAndHandler := func(ppmShipment models.PPMShipment, payload *internalmessages.SavePPMShipmentSignedCertification) (ppmops.ResubmitPPMShipmentDocumentationParams, ResubmitPPMShipmentDocumentationHandler) {
		endpoint := fmt.Sprintf(
			"/ppm-shipments/%s/submit-ppm-shipment-documentation/%s",
			ppmShipment.ID.String(),
			ppmShipment.SignedCertification.ID.String(),
		)

		request := httptest.NewRequest("PUT", endpoint, nil)

		request = suite.AuthenticateRequest(request, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		eTag := etag.GenerateEtag(ppmShipment.SignedCertification.UpdatedAt)
		params := ppmops.ResubmitPPMShipmentDocumentationParams{
			HTTPRequest:           request,
			PpmShipmentID:         handlers.FmtUUIDValue(ppmShipment.ID),
			SignedCertificationID: handlers.FmtUUIDValue(ppmShipment.SignedCertification.ID),
			SavePPMShipmentSignedCertificationPayload: payload,
			IfMatch: eTag,
		}

		handler := ResubmitPPMShipmentDocumentationHandler{
			suite.HandlerConfig(),
			submitter,
		}

		return params, handler
	}

	suite.Run("Returns an error if the PPM shipment is not found", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(suite.DB(), testdatagen.Assertions{})

		ppmShipment.ID = uuid.Must(uuid.NewV4())

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, "not found while looking for PPMShipment")
		}
	})

	suite.Run("Returns an error if the signed certification is not found", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(suite.DB(), testdatagen.Assertions{})

		ppmShipment.SignedCertification.ID = uuid.Must(uuid.NewV4())

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, "not found while looking for SignedCertification")
		}
	})

	suite.Run("Returns an error if the signed certification is not found", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(suite.DB(), testdatagen.Assertions{})

		ppmShipment.SignedCertification.ID = uuid.Must(uuid.NewV4())

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationNotFound{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationNotFound)

			suite.Contains(*errResponse.Payload.Detail, "not found while looking for SignedCertification")
		}
	})

	suite.Run("Returns an error if the PPM shipment is not in the right status", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsPaymentApproval(suite.DB(), testdatagen.Assertions{})

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString("certification text"),
			Signature:         handlers.FmtString("signature"),
			Date:              handlers.FmtDate(time.Now()),
		})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationConflict{}, response) {
			errResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationConflict)

			suite.Contains(
				*errResponse.Payload.Detail,
				fmt.Sprintf(
					"PPM shipment can't be set to %s because it's not in the %s status.",
					models.PPMShipmentStatusNeedsPaymentApproval,
					models.PPMShipmentStatusWaitingOnCustomer,
				),
			)
		}
	})

	suite.Run("Can successfully resubmit a PPM shipment for close out", func() {
		ppmShipment := testdatagen.MakePPMShipmentThatNeedsToBeResubmitted(suite.DB(), testdatagen.Assertions{})

		newCertText := "new certification text"
		newSignature := "new signature"
		newSignDate := time.Now().AddDate(0, 0, 1)

		params, handler := setUpParamsAndHandler(ppmShipment, &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString(newCertText),
			Signature:         handlers.FmtString(newSignature),
			Date:              handlers.FmtDate(newSignDate),
		})

		response := handler.Handle(params)

		if suite.IsType(&ppmops.ResubmitPPMShipmentDocumentationOK{}, response) {
			okResponse := response.(*ppmops.ResubmitPPMShipmentDocumentationOK)
			returnedPPMShipment := okResponse.Payload

			suite.NoError(returnedPPMShipment.Validate(strfmt.Default))

			suite.EqualUUID(ppmShipment.ID, returnedPPMShipment.ID)
			suite.Equal(string(models.PPMShipmentStatusNeedsPaymentApproval), string(returnedPPMShipment.Status))

			if suite.NotNil(returnedPPMShipment.SubmittedAt) {
				// For some reason, if we compare the time.Time objects directly, the test fails, but if we truncate to
				// the nanosecond level, they are equal.
				truncatedOriginalSubmittedAt := ppmShipment.SubmittedAt.Truncate(time.Nanosecond)
				truncatedNewSubmittedAt := handlers.FmtDateTimePtrToPop(returnedPPMShipment.SubmittedAt).Truncate(time.Nanosecond)
				suite.Equal(truncatedOriginalSubmittedAt, truncatedNewSubmittedAt)
			}

			suite.NotNil(returnedPPMShipment.SignedCertification)
			suite.EqualUUID(ppmShipment.SignedCertification.ID, returnedPPMShipment.SignedCertification.ID)

			suite.EqualUUID(
				ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember.User.ID,
				returnedPPMShipment.SignedCertification.SubmittingUserID,
			)

			suite.EqualUUID(ppmShipment.Shipment.MoveTaskOrder.ID, returnedPPMShipment.SignedCertification.MoveID)

			if suite.NotNil(returnedPPMShipment.SignedCertification.PpmID) {
				suite.EqualUUID(ppmShipment.ID, *returnedPPMShipment.SignedCertification.PpmID)
			}

			suite.Equal(
				string(models.SignedCertificationTypePPMPAYMENT),
				string(returnedPPMShipment.SignedCertification.CertificationType),
			)

			if suite.NotNil(returnedPPMShipment.SignedCertification.CertificationText) {
				suite.Equal(newCertText, *returnedPPMShipment.SignedCertification.CertificationText)
			}

			if suite.NotNil(returnedPPMShipment.SignedCertification.Signature) {
				suite.Equal(newSignature, *returnedPPMShipment.SignedCertification.Signature)
			}

			suite.True(
				newSignDate.Equal(handlers.FmtDatePtrToPop(returnedPPMShipment.SignedCertification.Date)),
				"Expected sign dates to be equal",
			)
		}
	})
}
