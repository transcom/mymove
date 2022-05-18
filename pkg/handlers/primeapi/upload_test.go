package primeapi

import (
	"fmt"
	"net/http/httptest"
	"testing"

	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"

	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/handlers"

	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateUploadHandler() {
	primeUser := testdatagen.MakeStubbedUser(suite.DB())

	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)

	testdatagen.MakeDefaultContractor(suite.DB())

	suite.T().Run("successful create upload", func(t *testing.T) {
		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)

		handler := CreateUploadHandler{
			handlerConfig,
			paymentrequest.NewPaymentRequestUploadCreator(handlerConfig.FileStorer()),
		}

		file := suite.Fixture("test.pdf")

		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: paymentRequest.ID.String(),
		}
		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadCreated{}, response)
	})

	suite.T().Run("create upload fail - invalid payment request ID format", func(t *testing.T) {
		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")

		handler := CreateUploadHandler{
			handlerConfig,
			paymentrequest.NewPaymentRequestUploadCreator(handlerConfig.FileStorer()),
		}

		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: badFormatID.String(),
		}
		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadUnprocessableEntity{}, response)
	})

	suite.T().Run("create upload fail - payment request not found", func(t *testing.T) {
		badFormatID := strfmt.UUID(uuid.Nil.String())

		handler := CreateUploadHandler{
			handlerConfig,
			paymentrequest.NewPaymentRequestUploadCreator(handlerConfig.FileStorer()),
		}

		file := suite.Fixture("test.pdf")

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: badFormatID.String(),
		}
		response := handler.Handle(params)

		suite.IsType(&uploadop.CreateUploadNotFound{}, response)
	})
}
