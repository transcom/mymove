package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/payment_requests"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestFetchPaymentRequestHandler() {
	expectedServiceItemName := "Test Service"
	expectedShipmentType := models.MTOShipmentTypeHHG

	move := testdatagen.MakeAvailableMove(suite.DB())
	// This should create all the other associated records we need.
	paymentServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
		Move: move,
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:  models.ServiceItemParamNameRequestedPickupDate,
			Type: models.ServiceItemParamTypeDate,
		},
	})
	paymentRequest := paymentServiceItemParam.PaymentServiceItem.PaymentRequest

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	suite.T().Run("successful fetch of payment request", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s", paymentRequest.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			paymentrequest.NewPaymentRequestFetcher(),
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestOK{}, response)
		okResponse := response.(*paymentrequestop.GetPaymentRequestOK)
		payload := okResponse.Payload
		paymentServiceItemParamPayload := payload.ServiceItems[0].PaymentServiceItemParams[0]

		suite.Equal(paymentRequest.ID.String(), payload.ID.String())
		suite.Equal(expectedServiceItemName, payload.ServiceItems[0].MtoServiceItemName)
		suite.EqualValues(expectedShipmentType, payload.ServiceItems[0].MtoShipmentType)

		suite.Equal(1, len(payload.ServiceItems))
		suite.Equal(paymentServiceItemParam.PaymentServiceItemID.String(), payload.ServiceItems[0].ID.String())
		suite.Equal(1, len(payload.ServiceItems[0].PaymentServiceItemParams))
		suite.Equal(paymentServiceItemParam.ID.String(), paymentServiceItemParamPayload.ID.String())
		suite.EqualValues(models.ServiceItemParamNameRequestedPickupDate, paymentServiceItemParamPayload.Key)
		suite.Equal(paymentServiceItemParam.Value, paymentServiceItemParamPayload.Value)
	})

	suite.T().Run("failed fetch for payment request - forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTOO,
		})
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s", paymentRequest.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTOO)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			paymentRequestFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestForbidden{}, response)
	})

	suite.T().Run("payment request not found", func(t *testing.T) {
		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(models.PaymentRequest{}, nil).Once()

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s", paymentRequest.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetPaymentRequestParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequest.ID.String()),
		}

		handler := GetPaymentRequestHandler{
			handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			paymentRequestFetcher,
		}
		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestGetPaymentRequestsForMoveHandler() {
	expectedServiceItemName := "Test Service"
	expectedShipmentType := models.MTOShipmentTypeHHG
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	move := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	// we need a mapping for the pickup address postal code to our user's gbloc
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		move.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	// This should create all the other associated records we need.
	paymentServiceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
		Move:           move,
		PaymentRequest: models.PaymentRequest{MoveTaskOrderID: move.ID, MoveTaskOrder: move},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:  models.ServiceItemParamNameRequestedPickupDate,
			Type: models.ServiceItemParamTypeDate,
		},
	})
	paymentRequest := paymentServiceItemParam.PaymentServiceItem.PaymentRequest
	paymentRequests := models.PaymentRequests{paymentRequest}

	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	suite.T().Run("Successful list fetch", func(t *testing.T) {
		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", move.Locator), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     move.Locator,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentrequest.NewPaymentRequestListFetcher(),
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestsForMoveOK{}, response)
		paymentRequestsResponse := response.(*paymentrequestop.GetPaymentRequestsForMoveOK)
		paymentRequestsPayload := paymentRequestsResponse.Payload
		paymentServiceItemParamPayload := paymentRequestsPayload[0].ServiceItems[0].PaymentServiceItemParams[0]

		suite.Equal(1, len(paymentRequestsPayload))
		suite.Equal(paymentRequests[0].ID.String(), paymentRequestsPayload[0].ID.String())
		suite.Equal(expectedServiceItemName, paymentRequestsPayload[0].ServiceItems[0].MtoServiceItemName)
		suite.EqualValues(expectedShipmentType, paymentRequestsPayload[0].ServiceItems[0].MtoShipmentType)

		suite.Equal(1, len(paymentRequestsPayload[0].ServiceItems))
		suite.Equal(paymentServiceItemParam.PaymentServiceItemID.String(), paymentRequestsPayload[0].ServiceItems[0].ID.String())
		suite.Equal(1, len(paymentRequestsPayload[0].ServiceItems[0].PaymentServiceItemParams))
		suite.Equal(paymentServiceItemParam.ID.String(), paymentServiceItemParamPayload.ID.String())
		suite.EqualValues(models.ServiceItemParamNameRequestedPickupDate, paymentServiceItemParamPayload.Key)
		suite.Equal(paymentServiceItemParam.Value, paymentServiceItemParamPayload.Value)
	})

	suite.T().Run("Failed list fetch - Not found error ", func(t *testing.T) {
		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestListByMove", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything).Return(nil, errors.New("not found")).Once()

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", "ABC123"), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     "ABC123",
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentRequestListFetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestNotFound{}, response)
	})

	suite.T().Run("Failed list fetch - Forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		paymentRequestListFetcher := &mocks.PaymentRequestListFetcher{}
		paymentRequestListFetcher.On("FetchPaymentRequestListByMove", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything).Return(&paymentRequests, nil).Once()

		request := httptest.NewRequest("GET", fmt.Sprintf("/moves/%s/payment-requests/", "ABC123"), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUserTOO)
		params := paymentrequestop.GetPaymentRequestsForMoveParams{
			HTTPRequest: request,
			Locator:     "ABC123",
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handler := GetPaymentRequestForMoveHandler{
			HandlerContext:            context,
			PaymentRequestListFetcher: paymentRequestListFetcher,
		}
		response := handler.Handle(params)
		suite.Assertions.IsType(&paymentrequestop.GetPaymentRequestsForMoveForbidden{}, response)
	})
}

func (suite *HandlerSuite) TestUpdatePaymentRequestStatusHandler() {
	paymentRequestID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	officeUserUUID, _ := uuid.NewV4()
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true, OfficeUser: models.OfficeUser{ID: officeUserUUID}})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	paymentRequest := models.PaymentRequest{
		ID:        paymentRequestID,
		IsFinal:   false,
		Status:    models.PaymentRequestStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())

	suite.T().Run("successful status update of payment request", func(t *testing.T) {
		pendingPaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil, ETag: etag.GenerateEtag(pendingPaymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: statusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload
		suite.Equal(models.PaymentRequestStatusReviewed.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.T().Run("successful status update of rejected payment request", func(t *testing.T) {
		pendingPaymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(pendingPaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED", RejectionReason: nil, ETag: etag.GenerateEtag(paymentRequest.UpdatedAt)},
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: statusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)
		suite.Logger().Error("")
		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)

		payload := response.(*paymentrequestop.UpdatePaymentRequestStatusOK).Payload
		suite.Equal(models.PaymentRequestStatusReviewedAllRejected.String(), string(payload.Status))
		suite.NotNil(payload.ReviewedAt)
	})

	suite.T().Run("prevent handler from updating payment request status to unapproved statuses", func(t *testing.T) {
		nonApprovedPRStatuses := [...]ghcmessages.PaymentRequestStatus{
			ghcmessages.PaymentRequestStatusSENTTOGEX,
			ghcmessages.PaymentRequestStatusRECEIVEDBYGEX,
			ghcmessages.PaymentRequestStatusPAID,
			ghcmessages.PaymentRequestStatusEDIERROR,
			ghcmessages.PaymentRequestStatusPENDING,
			ghcmessages.PaymentRequestStatusDEPRECATED,
		}

		for _, nonApprovedPRStatus := range nonApprovedPRStatuses {
			pendingPaymentRequest := testdatagen.MakeStubbedPaymentRequest(suite.DB())

			paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
			paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
				pendingPaymentRequest.ID).Return(pendingPaymentRequest, nil).Once()

			req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", pendingPaymentRequest.ID), nil)
			req = suite.AuthenticateOfficeRequest(req, officeUser)
			params := paymentrequestop.UpdatePaymentRequestStatusParams{
				HTTPRequest:      req,
				Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: nonApprovedPRStatus, RejectionReason: nil, ETag: etag.GenerateEtag(paymentRequest.UpdatedAt)},
				PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
			}

			handler := UpdatePaymentRequestStatusHandler{
				HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
				PaymentRequestStatusUpdater: statusUpdater,
				PaymentRequestFetcher:       paymentRequestFetcher,
			}

			response := handler.Handle(params)
			suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity(), response)
		}
	})

	suite.T().Run("failed status update of payment request - forbidden", func(t *testing.T) {
		officeUserTOO := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
			RoleType: roles.RoleTypeTOO,
		})

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(&paymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTOO)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusForbidden(), response)

	})

	suite.T().Run("successful status update of prime-available payment request", func(t *testing.T) {
		availableMove := testdatagen.MakeAvailableMove(suite.DB())
		availablePaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			Move: availableMove,
		})
		availablePaymentRequestID := availablePaymentRequest.ID

		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(&availablePaymentRequest, nil).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(availablePaymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", availablePaymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(availablePaymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusOK(), response)
		suite.HasWebhookNotification(availablePaymentRequestID, traceID)
	})

	suite.T().Run("unsuccessful status update of payment request (500)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, errors.New("Something bad happened")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusInternalServerError(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, not found (404)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.NewNotFoundError(paymentRequest.ID, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusNotFound(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, precondition failed (412)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.PreconditionFailedError{}).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusPreconditionFailed(), response)

	})

	suite.T().Run("unsuccessful status update of payment request, validation errors (422)", func(t *testing.T) {
		paymentRequestStatusUpdater := &mocks.PaymentRequestStatusUpdater{}
		paymentRequestStatusUpdater.On("UpdatePaymentRequestStatus", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything, mock.Anything).Return(nil, apperror.NewInvalidInputError(paymentRequestID, nil, nil, "")).Once()

		paymentRequestFetcher := &mocks.PaymentRequestFetcher{}
		paymentRequestFetcher.On("FetchPaymentRequest", mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything).Return(paymentRequest, nil).Once()

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/payment_request/%s/status", paymentRequestID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.UpdatePaymentRequestStatusParams{
			HTTPRequest:      req,
			Body:             &ghcmessages.UpdatePaymentRequestStatusPayload{Status: "REVIEWED", RejectionReason: nil},
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := UpdatePaymentRequestStatusHandler{
			HandlerContext:              handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			PaymentRequestStatusUpdater: paymentRequestStatusUpdater,
			PaymentRequestFetcher:       paymentRequestFetcher,
		}

		response := handler.Handle(params)

		suite.IsType(paymentrequestop.NewUpdatePaymentRequestStatusUnprocessableEntity(), response)
	})

}

func (suite *HandlerSuite) TestShipmentsSITBalanceHandler() {
	officeUserTIO := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	suite.T().Run("successful response of the shipments SIT Balance handler", func(t *testing.T) {
		now := time.Now()

		reviewedPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				Status: models.PaymentRequestStatusReviewed,
			},
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				AvailableToPrimeAt: &now,
			},
		})

		move := reviewedPaymentRequest.MoveTaskOrder

		pendingPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				SequenceNumber: 2,
			},
			Move: move,
		})

		year, month, day := time.Now().Date()
		originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
		sitDaysAllowance := 120

		doasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &originEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
			MTOShipment: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				SITDaysAllowance: &sitDaysAllowance,
			},
			Move: move,
		})

		shipment := doasit.MTOShipment
		// Creates the payment service item for DOASIT w/ SIT start date param
		doasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: originEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentServiceItem: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
			PaymentRequest: reviewedPaymentRequest,
			MTOServiceItem: doasit,
			Move:           move,
		})

		paymentEndDate := originEntryDate.Add(time.Hour * 24 * 30)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: paymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "30",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: doasitParam.PaymentServiceItem,
			PaymentRequest:     reviewedPaymentRequest,
			MTOServiceItem:     doasit,
			Move:               move,
		})

		destinationEntryDate := time.Date(year, month, day-90, 0, 0, 0, 0, time.UTC)
		ddasit := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				Status:       models.MTOServiceItemStatusApproved,
				SITEntryDate: &destinationEntryDate,
			},
			ReService: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		// Creates the payment service item for DOASIT w/ SIT start date param
		ddasitParam := testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationEntryDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestStart,
			},
			PaymentRequest: pendingPaymentRequest,
			MTOServiceItem: ddasit,
			Move:           move,
		})

		destinationPaymentEndDate := destinationEntryDate.Add(time.Hour * 24 * 60)
		// Creates the SIT end date param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: destinationPaymentEndDate.Format("2006-01-02"),
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameSITPaymentRequestEnd,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		// Creates the NumberDaysSIT param for existing DOASIT payment request service item
		testdatagen.MakePaymentServiceItemParam(suite.DB(), testdatagen.Assertions{
			PaymentServiceItemParam: models.PaymentServiceItemParam{
				Value: "60",
			},
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key: models.ServiceItemParamNameNumberDaysSIT,
			},
			PaymentServiceItem: ddasitParam.PaymentServiceItem,
			PaymentRequest:     pendingPaymentRequest,
			MTOServiceItem:     ddasit,
			Move:               move,
		})

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s/shipments-payment-sit-balance", pendingPaymentRequest.ID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTIO)

		params := paymentrequestop.GetShipmentsPaymentSITBalanceParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(pendingPaymentRequest.ID.String()),
		}

		handler := ShipmentsSITBalanceHandler{
			HandlerContext:             handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ShipmentsPaymentSITBalance: paymentrequest.NewPaymentRequestShipmentsSITBalance(),
		}

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetShipmentsPaymentSITBalanceOK{}, response)

		payload := response.(*paymentrequestop.GetShipmentsPaymentSITBalanceOK).Payload

		suite.NotNil(payload)
		suite.Len(payload, 1)
		shipmentSITBalance := payload[0]

		suite.Equal(shipment.ID.String(), shipmentSITBalance.ShipmentID.String())
		suite.Equal(int64(120), shipmentSITBalance.TotalSITDaysAuthorized)
		suite.Equal(int64(60), shipmentSITBalance.PendingSITDaysInvoiced)
		suite.Equal(int64(30), shipmentSITBalance.TotalSITDaysRemaining)
		suite.Equal(destinationPaymentEndDate.Format("2006-01-02"), shipmentSITBalance.PendingBilledEndDate.String())
		suite.Equal(int64(30), *shipmentSITBalance.PreviouslyBilledDays)
	})

	suite.T().Run("returns 403 unauthorized when request is not made by TIO office user", func(t *testing.T) {
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		paymentRequestID := uuid.Must(uuid.NewV4())

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s/shipments-payment-sit-balance", paymentRequestID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := paymentrequestop.GetShipmentsPaymentSITBalanceParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := ShipmentsSITBalanceHandler{
			HandlerContext:             handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ShipmentsPaymentSITBalance: paymentrequest.NewPaymentRequestShipmentsSITBalance(),
		}

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetShipmentsPaymentSITBalanceForbidden{}, response)
	})

	suite.T().Run("returns 404 not found when payment request does not exist", func(t *testing.T) {
		paymentRequestID := uuid.Must(uuid.NewV4())

		req := httptest.NewRequest("GET", fmt.Sprintf("/payment-requests/%s/shipments-payment-sit-balance", paymentRequestID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUserTIO)

		params := paymentrequestop.GetShipmentsPaymentSITBalanceParams{
			HTTPRequest:      req,
			PaymentRequestID: strfmt.UUID(paymentRequestID.String()),
		}

		handler := ShipmentsSITBalanceHandler{
			HandlerContext:             handlers.NewHandlerContext(suite.DB(), suite.Logger()),
			ShipmentsPaymentSITBalance: paymentrequest.NewPaymentRequestShipmentsSITBalance(),
		}

		response := handler.Handle(params)

		suite.IsType(&paymentrequestop.GetShipmentsPaymentSITBalanceNotFound{}, response)
	})
}
