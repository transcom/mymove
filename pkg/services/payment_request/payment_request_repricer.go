package paymentrequest

import (
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestRepricer struct {
	paymentRequestCreator services.PaymentRequestCreator
}

// NewPaymentRequestRepricer returns a new payment request repricer
func NewPaymentRequestRepricer(paymentRequestCreator services.PaymentRequestCreator) services.PaymentRequestRepricer {
	return &paymentRequestRepricer{
		paymentRequestCreator: paymentRequestCreator,
	}
}

func (p *paymentRequestRepricer) RepricePaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (*models.PaymentRequest, error) {
	var newPaymentRequest *models.PaymentRequest

	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// Fetch the payment request and payment service items from the old request.
		var oldPaymentRequest models.PaymentRequest
		err := txnAppCtx.DB().
			EagerPreload(
				"PaymentServiceItems.MTOServiceItem.ReService",
				"ProofOfServiceDocs.PrimeUploads",
			).
			Find(&oldPaymentRequest, paymentRequestID)
		if err != nil {
			return err
		}

		// Re-create the payment request arg including service items, then call the create service (which should
		// price it with current inputs).
		inputPaymentRequest := buildPaymentRequestForRepricing(oldPaymentRequest)
		newPaymentRequest, err = p.paymentRequestCreator.CreatePaymentRequest(txnAppCtx, &inputPaymentRequest)
		if err != nil {
			return err
		}

		// Set the (now) old payment request's status.
		// TODO: We need a better status for this -- something like "REPRICED".
		err = updateOldPaymentRequest(appCtx, &oldPaymentRequest)
		if err != nil {
			return err
		}

		// Duplicate the proof-of-service upload associations to the new payment request.
		err = associateProofOfServiceDocs(appCtx, oldPaymentRequest.ProofOfServiceDocs, newPaymentRequest)
		if err != nil {
			return err
		}

		// Link the new payment request to the old one.
		err = linkNewToOldPaymentRequest(appCtx, newPaymentRequest, &oldPaymentRequest)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return newPaymentRequest, nil
}

// buildPaymentRequestForRepricing builds up the expected payment request data based upon the old payment request.
func buildPaymentRequestForRepricing(oldPaymentRequest models.PaymentRequest) models.PaymentRequest {
	newPaymentRequest := models.PaymentRequest{
		IsFinal:         oldPaymentRequest.IsFinal,
		MoveTaskOrderID: oldPaymentRequest.MoveTaskOrderID,
	}

	var newPaymentServiceItems models.PaymentServiceItems
	for _, oldPaymentServiceItem := range oldPaymentRequest.PaymentServiceItems {
		newPaymentServiceItem := models.PaymentServiceItem{
			MTOServiceItemID: oldPaymentServiceItem.MTOServiceItemID,
			MTOServiceItem:   oldPaymentServiceItem.MTOServiceItem,
		}

		newPaymentServiceItems = append(newPaymentServiceItems, newPaymentServiceItem)
	}

	sort.SliceStable(newPaymentServiceItems, func(i, j int) bool {
		return newPaymentServiceItems[i].MTOServiceItem.ReService.Priority < newPaymentServiceItems[j].MTOServiceItem.ReService.Priority
	})

	newPaymentRequest.PaymentServiceItems = newPaymentServiceItems

	return newPaymentRequest
}

func updateOldPaymentRequest(appCtx appcontext.AppContext, oldPaymentRequest *models.PaymentRequest) error {
	newStatus := models.PaymentRequestStatusReviewedAllRejected
	oldPaymentRequest.Status = newStatus
	verrs, err := appCtx.DB().ValidateAndUpdate(oldPaymentRequest)
	if err != nil {
		return fmt.Errorf("failed to set old payment request status to %s: %w", newStatus, err)
	}
	if verrs.HasAny() {
		return fmt.Errorf("failed to validate old payment request when setting status to %s: %w", newStatus, verrs)
	}

	return nil
}

// associateProofOfServiceDocs duplicates/associates a set of proof of service doc relationships to a payment request.
func associateProofOfServiceDocs(appCtx appcontext.AppContext, proofOfServiceDocs models.ProofOfServiceDocs, newPaymentRequest *models.PaymentRequest) error {
	for _, proofOfServiceDoc := range proofOfServiceDocs {
		newProofOfServiceDoc := models.ProofOfServiceDoc{
			PaymentRequestID: newPaymentRequest.ID,
		}
		verrs, err := appCtx.DB().ValidateAndCreate(&newProofOfServiceDoc)
		if err != nil {
			return fmt.Errorf("failed to create proof of service doc for new payment request ID %s: %w", newPaymentRequest.ID, err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate proof of service doc for new payment request ID %s: %w", newPaymentRequest.ID, verrs)
		}

		for _, primeUpload := range proofOfServiceDoc.PrimeUploads {
			newPrimeUpload := models.PrimeUpload{
				ProofOfServiceDocID: newProofOfServiceDoc.ID,
				ContractorID:        primeUpload.ContractorID,
				UploadID:            primeUpload.UploadID,
				DeletedAt:           primeUpload.DeletedAt,
			}
			verrs, err := appCtx.DB().ValidateAndCreate(&newPrimeUpload)
			if err != nil {
				return fmt.Errorf("failed to create prime upload for new payment request ID %s and new proof of service doc ID %s: %w", newPaymentRequest.ID, newProofOfServiceDoc.ID, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate proof of service doc for new payment request ID %s and new proof of service doc ID %s: %w", newPaymentRequest.ID, newProofOfServiceDoc.ID, verrs)
			}

			newProofOfServiceDoc.PrimeUploads = append(newProofOfServiceDoc.PrimeUploads, newPrimeUpload)
		}

		newPaymentRequest.ProofOfServiceDocs = append(newPaymentRequest.ProofOfServiceDocs, newProofOfServiceDoc)
	}

	return nil
}

// linkNewToOldPaymentRequest links a new payment request to the old payment request that was repriced.
func linkNewToOldPaymentRequest(appCtx appcontext.AppContext, newPaymentRequest *models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
	newPaymentRequest.RepricedPaymentRequestID = &oldPaymentRequest.ID
	verrs, err := appCtx.DB().ValidateAndUpdate(newPaymentRequest)
	if err != nil {
		return fmt.Errorf("failed to set new payment request to old payment request ID %s: %w", oldPaymentRequest.ID, err)
	}
	if verrs.HasAny() {
		return fmt.Errorf("failed to validate old payment request when setting status to %s: %w", oldPaymentRequest.ID, verrs)
	}

	return nil
}
