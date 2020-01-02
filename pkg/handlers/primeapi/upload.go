package primeapi

import (
	"bytes"
	"fmt"
	"io"

	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/uploader"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	uploadop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/uploads"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func payloadForPaymentRequestUploadModel(u models.Upload) *primemessages.Upload {
	return &primemessages.Upload{
		Bytes:       &u.Bytes,
		ContentType: &u.ContentType,
		Filename:    &u.Filename,
	}
}

type CreateUploadHandler struct {
	handlers.HandlerContext
	services.PaymentRequestUploadCreator
}

func (h *CreateUploadHandler) convertFileReadCloserToAfero(file io.ReadCloser, logger handlers.Logger) (afero.File, error) {
	fs := afero.NewMemMapFs()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read payment request upload file into buffer: %w", err)
	}

	newStr := buf.String()
	aferoFile, err := fs.Create(newStr)
	if err != nil {
		return nil, fmt.Errorf("afero.Create Failed in payment request upload creation: %w", err)
	}
	defer aferoFile.Close()

	return aferoFile, err
}

func (h *CreateUploadHandler) Handle(params uploadop.CreateUploadParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	paymentRequestID, err := uuid.FromString(params.PaymentRequestID)
	if err != nil {
		logger.Error("error creating uuid from string", zap.Error(err))
	}
	stubbedUser := testdatagen.MakeOfficeUser(h.DB(), testdatagen.Assertions{
		User: models.User{
			LoginGovEmail: "fakeuploaduser@example.com",
		},
	})

	aferoFile, err := h.convertFileReadCloserToAfero(params.File, logger)
	if err != nil {
		logger.Error("error converting file to afero type", zap.Error(err))
		return uploadop.NewCreateUploadInternalServerError()
	}
	file := uploader.File{
		File: aferoFile,
	}
	uploadCreator := paymentrequest.NewPaymentRequestUploadCreator(h.DB(), logger, h.FileStorer())
	createdUpload, err := uploadCreator.CreateUpload(file, paymentRequestID, stubbedUser.ID)
	if err != nil {
		logger.Error("cannot create payment request upload", zap.Error(err))
		return uploadop.NewCreateUploadBadRequest()
	}

	returnPayload := payloadForPaymentRequestUploadModel(*createdUpload)
	return uploadop.NewCreateUploadCreated().WithPayload(returnPayload)
}
