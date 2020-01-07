package primeapi

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestCreateUploadHandler() {
	primeUser := testdatagen.MakeDefaultUser(suite.DB())
	uploadID, _ := uuid.FromString("e2e79f36-de9e-4a52-9566-47fa3834b359")

	paymentRequest := testdatagen.MakeDefaultPaymentRequest(suite.DB())

	suite.T().Run("successful create upload", func(t *testing.T) {
		upload := models.Upload{
			ID:          uploadID,
			DocumentID:  nil,
			Document:    models.Document{},
			UploaderID:  primeUser.ID,
			Filename:    "test.pdf",
			Bytes:       42330,
			ContentType: "application/json",
			Checksum:    "asdfsadfasdf",
			StorageKey:  "storagekeyvalue",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}

		paymentRequestUploadCreator := &mocks.PaymentRequestUploadCreator{}
		paymentRequestUploadCreator.On(
			"CreateUpload",
			mock.AnythingOfType("uploader.File"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(
			&upload, nil).Once()

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)

		handler := CreateUploadHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestUploadCreator,
		}
		file, err := os.Open("../fixtures/test.pdf")
		suite.NoError(err)

		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: paymentRequest.ID.String(),
		}
		response := handler.Handle(params)
		file.Close()

		suite.IsType(&uploadop.CreateUpload{}, response)
	})

	suite.T().Run("create upload fail - invalid payment request ID format", func(t *testing.T) {
		badFormatID := strfmt.UUID("gb7b134a-7c44-45f2-9114-bb0831cc5db3")
		paymentRequestUploadCreator := &mocks.PaymentRequestUploadCreator{}
		paymentRequestUploadCreator.On(
			"CreateUpload",
			mock.AnythingOfType("uploader.File"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(
			&models.Upload{}, nil).Once()

		handler := CreateUploadHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestUploadCreator,
		}
		file, err := os.Open("../fixtures/test.pdf")
		suite.NoError(err)

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)
		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             file,
			PaymentRequestID: badFormatID.String(),
		}
		response := handler.Handle(params)
		file.Close()
		suite.IsType(&uploadop.CreateUploadBadRequest{}, response)
	})

	suite.T().Run("create upload fail - invalid file format", func(t *testing.T) {
		paymentRequestUploadCreator := &mocks.PaymentRequestUploadCreator{}
		paymentRequestUploadCreator.On(
			"CreateUpload",
			mock.AnythingOfType("uploader.File"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(
			&models.Upload{}, nil).Once()

		handler := CreateUploadHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			paymentRequestUploadCreator,
		}
		badFile, err := os.Open("../fixtures/test.pdf") // TODO: make this a failing file
		suite.NoError(err)

		req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests/%s/uploads", paymentRequest.ID), nil)
		req = suite.AuthenticateUserRequest(req, primeUser)

		params := uploadop.CreateUploadParams{
			HTTPRequest:      req,
			File:             badFile,
			PaymentRequestID: paymentRequest.ID.String(),
		}
		response := handler.Handle(params)
		badFile.Close()

		suite.IsType(&uploadop.CreateUploadBadRequest{}, response)
	})
}
