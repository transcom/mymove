package paymentrequest

import (
	"fmt"

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

func NewUploadCreator(db *pop.Connection, logger storage.Logger, fileStorer storage.FileStorer) services.PaymentRequestUploadCreator {
	return &paymentRequestUploadCreator{db, logger, fileStorer, uploader.MaxFileSizeLimit}
}

func (p *paymentRequestUploadCreator) CreateUpload(file uploader.File, paymentRequestID uuid.UUID, userID uuid.UUID) (*models.Upload, error) {
	upload := models.Upload{}
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		newUploader, err := uploader.NewUploader(tx, p.logger, p.fileStorer, p.fileSizeLimit)
		if err != nil {
			return fmt.Errorf("cannot create uploader in payment request uploadCreator: %w", err)
		}

		upload, verrs, err := newUploader.CreateUpload(userID, file, uploader.AllowedTypesPaymentRequest)
		if err != nil {
			return fmt.Errorf("failure creating payment request upload: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating payment request upload: %w", verrs)
		}

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

	return &upload, nil
}
