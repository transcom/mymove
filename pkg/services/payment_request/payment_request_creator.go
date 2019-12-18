package paymentrequest

import (
	"fmt"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type paymentRequestCreator struct {
	db *pop.Connection
	logger storage.Logger
	fileStorer storage.FileStorer
	fileSizeLimit uploader.ByteSize
}

func NewPaymentRequestCreator(db *pop.Connection, logger storage.Logger, fileStorer storage.FileStorer, fileSizeLimit uploader.ByteSize) services.PaymentRequestCreator {
	return &paymentRequestCreator{db, logger, fileStorer, uploader.MaxFileSizeLimit}
}

func (p *paymentRequestCreator) CreatePaymentRequest(paymentRequest *models.PaymentRequest) (*models.PaymentRequest, error) {
	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		now := time.Now()

		// Verify that the MTO ID exists
		var moveTaskOrder models.MoveTaskOrder
		err := tx.Find(&moveTaskOrder, paymentRequest.MoveTaskOrderID)
		if err != nil {
			return fmt.Errorf("could not find MoveTaskOrderID [%s]: %w", paymentRequest.MoveTaskOrderID, err)
		}
		paymentRequest.MoveTaskOrder = moveTaskOrder

		paymentRequest.Status = models.PaymentRequestStatusPending
		paymentRequest.RequestedAt = now

		// Create the payment request first
		verrs, err := tx.ValidateAndCreate(paymentRequest)
		if err != nil {
			return fmt.Errorf("failure creating payment request: %w", err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("validation error creating payment request: %w", verrs)
		}

		// Create each payment service item for the payment request
		var newPaymentServiceItems models.PaymentServiceItems
		for _, paymentServiceItem := range paymentRequest.PaymentServiceItems {
			// Verify that the service item ID exists
			var mtoServiceItem models.MTOServiceItem
			err := tx.Find(&mtoServiceItem, paymentServiceItem.ServiceItemID)
			if err != nil {
				return fmt.Errorf("could not find ServiceItemID [%s]: %w", paymentServiceItem.ServiceItemID, err)
			}
			paymentServiceItem.ServiceItem = mtoServiceItem

			paymentServiceItem.PaymentRequestID = paymentRequest.ID
			paymentServiceItem.PaymentRequest = *paymentRequest
			paymentServiceItem.Status = models.PaymentServiceItemStatusRequested
			paymentServiceItem.PriceCents = unit.Cents(0) // TODO: Placeholder until we have pricing ready.
			paymentServiceItem.RequestedAt = now

			verrs, err := tx.ValidateAndCreate(&paymentServiceItem)
			if err != nil {
				return fmt.Errorf("failure creating payment service item: %w", err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("validation error creating payment service item: %w", verrs)
			}

			// Create each payment service item parameter for the payment service item
			var newPaymentServiceItemParams models.PaymentServiceItemParams
			for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
				// If the ServiceItemParamKeyID is provided, verify it exists; otherwise, lookup
				// via the IncomingKey field
				var serviceItemParamKey models.ServiceItemParamKey
				if paymentServiceItemParam.ServiceItemParamKeyID != uuid.Nil {
					err := tx.Find(&serviceItemParamKey, paymentServiceItemParam.ServiceItemParamKeyID)
					if err != nil {
						return fmt.Errorf("could not find ServiceItemParamKeyID [%s]: %w", paymentServiceItemParam.ServiceItemParamKeyID, err)
					}
				} else {
					err := tx.Where("key = ?", paymentServiceItemParam.IncomingKey).First(&serviceItemParamKey)
					if err != nil {
						return fmt.Errorf("could not find param key [%s]: %w", paymentServiceItemParam.IncomingKey, err)
					}
				}
				paymentServiceItemParam.ServiceItemParamKeyID = serviceItemParamKey.ID
				paymentServiceItemParam.ServiceItemParamKey = serviceItemParamKey

				paymentServiceItemParam.PaymentServiceItemID = paymentServiceItem.ID
				paymentServiceItemParam.PaymentServiceItem = paymentServiceItem

				verrs, err := tx.ValidateAndCreate(&paymentServiceItemParam)
				if err != nil {
					return fmt.Errorf("failure creating payment service item param: %w", err)
				}
				if verrs.HasAny() {
					return fmt.Errorf("validation error creating payment service item param: %w", verrs)
				}

				newPaymentServiceItemParams = append(newPaymentServiceItemParams, paymentServiceItemParam)
			}
			paymentServiceItem.PaymentServiceItemParams = newPaymentServiceItemParams

			newPaymentServiceItems = append(newPaymentServiceItems, paymentServiceItem)
		}
		paymentRequest.PaymentServiceItems = newPaymentServiceItems

		//newUploader, err := uploader.NewUploader(p.db, p.logger, p.fileStorer, p.fileSizeLimit)
		//if err != nil {
		//	return fmt.Errorf("cannot create uploader in paymentRequestCreator: %w", err)
		//}
		//
		//stubbedUserID, err := uuid.FromString("d7b09d8d-541b-4f30-bd59-c5f5ba87ac2e")
		//if err != nil {
		//	return fmt.Errorf("cannot create uuid form string in paymentRequestCreator %w", err)
		//}
		//
		//// loop through proof of service docs and create uploads and proof of service docs
		//for _, doc := range paymentRequest.ProofOfServiceDocs {
		//	// call CreateUpload (pass in proofOfServiceID)
		//	upload := newUploader.CreateUpload(stubbedUserID, doc.Upload, uploader.AllowedTypesServiceMember) // figure out what file types we accept
		//	// create proof of service doc
		//	proofOfServiceDoc := models.ProofOfServiceDoc{
		//		PaymentRequestID: paymentRequest.ID,
		//		UploadID: upload.ID,
		//	}
		//	verrs, err := tx.ValidateAndCreate(&proofOfServiceDoc)
		//	if err != nil {
		//		return fmt.Errorf("failure creating proof of service doc: %w", err)
		//	}
		//	if verrs.HasAny() {
		//		return fmt.Errorf("validation error creating proof of service doc: %w", verrs)
		//	}
		//}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return paymentRequest, nil
}
