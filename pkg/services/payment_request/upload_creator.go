package paymentrequest

import (
	"database/sql"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

const (
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

type paymentRequestUploadCreator struct {
	fileStorer    storage.FileStorer
	fileSizeLimit uploader.ByteSize
}

// NewPaymentRequestUploadCreator returns a new payment request upload creator
func NewPaymentRequestUploadCreator(fileStorer storage.FileStorer) services.PaymentRequestUploadCreator {
	return &paymentRequestUploadCreator{fileStorer, uploader.MaxFileSizeLimit}
}

func (p *paymentRequestUploadCreator) assembleUploadFilePathName(appCtx appcontext.AppContext, paymentRequestID uuid.UUID, filename string) (string, error) {
	var paymentRequest models.PaymentRequest
	err := appCtx.DB().Where("id=$1", paymentRequestID).First(&paymentRequest)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return "", apperror.NewNotFoundError(paymentRequestID, "")
		default:
			return "", apperror.NewQueryError("PaymentRequest", err, "")
		}
	}

	newfilename := time.Now().Format(VersionTimeFormat) + "-" + filename
	uploadFilePath := fmt.Sprintf("/payment-request-uploads/mto-%s/payment-request-%s", paymentRequest.MoveTaskOrderID, paymentRequest.ID)
	uploadFileName := path.Join(uploadFilePath, newfilename)

	return uploadFileName, err
}

func (p *paymentRequestUploadCreator) CreateUpload(appCtx appcontext.AppContext, file io.ReadCloser, paymentRequestID uuid.UUID, contractorID uuid.UUID, uploadFilename string) (*models.Upload, error) {
	var upload *models.Upload
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		newUploader, err := uploader.NewPrimeUploader(p.fileStorer, p.fileSizeLimit)
		if err != nil {
			if err == uploader.ErrFileSizeLimitExceedsMax {
				return apperror.NewBadDataError(err.Error())
			}
			return err
		}

		fileName, err := p.assembleUploadFilePathName(txnAppCtx, paymentRequestID, uploadFilename)
		if err != nil {
			return err
		}

		aFile, err := newUploader.PrepareFileForUpload(txnAppCtx, file, fileName)
		if err != nil {
			return err
		}

		newUploader.SetUploadStorageKey(fileName)

		var paymentRequest models.PaymentRequest
		err = txnAppCtx.DB().Find(&paymentRequest, paymentRequestID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(paymentRequestID, "")
			default:
				return apperror.NewQueryError("PaymentRequest", err, "")
			}
		}
		// create proof of service doc
		proofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: paymentRequestID,
			PaymentRequest:   paymentRequest,
		}
		verrs, err := txnAppCtx.DB().ValidateAndCreate(&proofOfServiceDoc)
		if err != nil {
			return fmt.Errorf("failure creating proof of service doc: %w", err) // server err
		}
		if verrs.HasAny() {
			return apperror.NewInvalidCreateInputError(verrs, "validation error with creating proof of service doc")
		}

		posID := &proofOfServiceDoc.ID
		primeUpload, verrs, err := newUploader.CreatePrimeUploadForDocument(txnAppCtx, posID, contractorID, uploader.File{File: aFile}, uploader.AllowedTypesPaymentRequest)
		if verrs.HasAny() {
			return apperror.NewInvalidCreateInputError(verrs, "validation error with creating payment request")
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
