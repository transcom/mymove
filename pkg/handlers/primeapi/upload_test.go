package primeapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/models"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateUploadHandler() {

	fakeS3 := storageTest.NewFakeS3Storage(true)

	setupTestData := func() (CreateUploadHandler, models.PaymentRequest) {
		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		handler := CreateUploadHandler{
			handlerConfig,
			paymentrequest.NewPaymentRequestUploadCreator(handlerConfig.FileStorer()),
		}
		paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())
		testdatagen.MakeDefaultContractor(suite.DB())
		return handler, paymentRequest
	}

	suite.Run("successful create upload", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, paymentRequest := setupTestData()
		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)

		file := suite.Fixture("test.pdf")

		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: paymentRequest.ID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadCreated{}, response)
		responsePayload := response.(*uploadop.CreateUploadCreated).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})

	suite.Run("create upload fail - invalid payment request ID format", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, paymentRequest := setupTestData()

		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")
		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: badFormatID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadUnprocessableEntity{}, response)

		// Validate outgoing payload
		// TODO: Can't validate the response because of the issue noted below. Figure out a way to
		//   either alter the service or relax the swagger requirements.
		// responsePayload := response.(*uploadop.CreateUploadCreated).Payload
		// suite.NoError(responsePayload.Validate(strfmt.Default))
		// Handler is not setting any validation errors so InvalidFields won't be added to the payload.
	})

	suite.Run("create upload fail - payment request not found", func() {
		primeUser := factory.BuildUser(nil, nil, nil)
		handler, paymentRequest := setupTestData()

		badFormatID, _ := uuid.NewV4()
		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: badFormatID.String(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadNotFound{}, response)
		responsePayload := response.(*uploadop.CreateUploadNotFound).Payload

		// Validate outgoing payload
		suite.NoError(responsePayload.Validate(strfmt.Default))
	})
}
