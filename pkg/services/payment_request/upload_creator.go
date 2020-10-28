package paymentrequest

import (
	"fmt"
	"io"
	"path"
	"time"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

const (
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

type paymentRequestUploadCreator struct {
	db            *pop.Connection
	logger        Logger
	fileStorer    storage.FileStorer
	fileSizeLimit uploader.ByteSize
}

// NewPaymentRequestUploadCreator returns a new payment request upload creator
func NewPaymentRequestUploadCreator(db *pop.Connection, logger Logger, fileStorer storage.FileStorer) services.PaymentRequestUploadCreator {
	return &paymentRequestUploadCreator{db, logger, fileStorer, uploader.MaxFileSizeLimit}
}

func (p *paymentRequestUploadCreator) assembleUploadFilePathName(paymentRequestID uuid.UUID, filename string) (string, error) {
	var paymentRequest models.PaymentRequest
	err := p.db.Where("id=$1", paymentRequestID).First(&paymentRequest)
	if err != nil {
		return "", services.NewNotFoundError(paymentRequestID, "")
	}

	newfilename := time.Now().Format(VersionTimeFormat) + "-" + filename
	uploadFilePath := fmt.Sprintf("/payment-request-uploads/mto-%s/payment-request-%s", paymentRequest.MoveTaskOrderID, paymentRequest.ID)
	uploadFileName := path.Join(uploadFilePath, newfilename)

	return uploadFileName, err
}

func (p *paymentRequestUploadCreator) CreateUpload(file io.ReadCloser, paymentRequestID uuid.UUID, contractorID uuid.UUID, uploadFilename string) (*models.Upload, error) {
	var upload *models.Upload
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		newUploader, err := uploader.NewPrimeUploader(tx, p.logger, p.fileStorer, p.fileSizeLimit)
		if err != nil {
			if err == uploader.ErrFileSizeLimitExceedsMax {
				return services.NewBadDataError(err.Error())
			}
			return err
		}

		fileName, err := p.assembleUploadFilePathName(paymentRequestID, uploadFilename)
		if err != nil {
			return err
		}

		aFile, err := newUploader.PrepareFileForUpload(file, fileName)
		if err != nil {
			return err
		}

		newUploader.SetUploadStorageKey(fileName)

		var paymentRequest models.PaymentRequest
		err = tx.Find(&paymentRequest, paymentRequestID)
		if err != nil {
			return services.NewNotFoundError(paymentRequestID, "")
		}
		// create proof of service doc
		proofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: paymentRequestID,
			PaymentRequest:   paymentRequest,
		}
		verrs, err := tx.ValidateAndCreate(&proofOfServiceDoc)
		if err != nil {
			return fmt.Errorf("failure creating proof of service doc: %w", err) // server err
		}
		if verrs.HasAny() {
			return services.NewInvalidCreateInputError(verrs, "validation error with creating proof of service doc")
		}

		posID := &proofOfServiceDoc.ID
		primeUpload, verrs, err := newUploader.CreatePrimeUploadForDocument(posID, contractorID, uploader.File{File: aFile}, uploader.AllowedTypesPaymentRequest)
		if verrs.HasAny() {
			return services.NewInvalidCreateInputError(verrs, "validation error with creating payment request")
		}
		if err != nil {
			return fmt.Errorf("failure creating payment request primeUpload: %w", err) // server err
		}
		upload = &primeUpload.Upload
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return upload, nil
}
