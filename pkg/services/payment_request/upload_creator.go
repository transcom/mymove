package paymentrequest

import (
	"fmt"
	"io"
	"path"
	"time"

	"github.com/spf13/afero"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

type paymentRequestUploadCreator struct {
	db            *pop.Connection
	logger        storage.Logger
	fileStorer    storage.FileStorer
	fileSizeLimit uploader.ByteSize
}

func NewPaymentRequestUploadCreator(db *pop.Connection, logger storage.Logger, fileStorer storage.FileStorer) services.PaymentRequestUploadCreator {
	return &paymentRequestUploadCreator{db, logger, fileStorer, uploader.MaxFileSizeLimit}
}

func (p *paymentRequestUploadCreator) convertFileReadCloserToAfero(file io.ReadCloser, paymentRequestID uuid.UUID) (afero.File, error) {

	fs := afero.NewMemMapFs()

	fileName, err := p.assembleUploadFilePathName(paymentRequestID)
	if err != nil {
		return nil, fmt.Errorf("could not assemble upload filepath name %w", err)
	}
	aferoFile, err := fs.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("afero.Create Failed in payment request upload creation: %w", err)
	}

	_, err = io.Copy(aferoFile, file)
	if err != nil {
		return nil, fmt.Errorf("Error copying to afero file %w", err)
	}

	return aferoFile, err
}

func (p *paymentRequestUploadCreator) assembleUploadFilePathName(paymentRequestID uuid.UUID) (string, error) {
	paymentRequest, err := models.FetchPaymentRequestByID(p.db, paymentRequestID)
	if err != nil {
		return "", err
	}
	filename := "timestamp-" + time.Now().String()
	uploadFilePath := fmt.Sprintf("/app/payment-request-uploads/mto-%s/payment-request-%s", paymentRequest.MoveTaskOrderID, paymentRequestID)
	uploadFileName := path.Join(uploadFilePath, filename)

	return uploadFileName, err
}

func (p *paymentRequestUploadCreator) CreateUpload(file io.ReadCloser, paymentRequestID uuid.UUID, userID uuid.UUID) (*models.Upload, error) {
	var upload *models.Upload
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		newUploader, err := uploader.NewUploader(tx, p.logger, p.fileStorer, p.fileSizeLimit)
		if err != nil {
			return fmt.Errorf("cannot create uploader in payment request uploadCreator: %w", err)
		}

		aferoFile, err := p.convertFileReadCloserToAfero(file, paymentRequestID)
		if err != nil {
			return fmt.Errorf("failure to convert payment request upload to afero file: %w", err)
		}

		var verrs *validate.Errors
		upload, verrs, err = newUploader.CreateUpload(userID, uploader.File{File: aferoFile}, uploader.AllowedTypesPaymentRequest)
		if err != nil {
			return fmt.Errorf("failure creating payment request upload: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating payment request upload: %w", verrs)
		}

		aferoFile.Close()

		var paymentRequest models.PaymentRequest
		err = tx.Find(&paymentRequest, paymentRequestID)
		if err != nil {
			return fmt.Errorf("could not find PaymentRequestID [%s]: %w", paymentRequestID, err)
		}
		// create proof of service doc
		proofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: paymentRequestID,
			PaymentRequest:   paymentRequest,
			UploadID:         upload.ID,
			Upload:           *upload,
		}

		verrs, err = tx.ValidateAndCreate(&proofOfServiceDoc)
		if err != nil {
			return fmt.Errorf("failure creating proof of service doc: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating proof of service doc: %w", verrs)
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return upload, nil
}
