package paymentrequest

import (
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestRecalculator struct {
	paymentRequestCreator services.PaymentRequestCreator
}

// NewPaymentRequestRecalculator returns a new payment request recalculator
func NewPaymentRequestRecalculator(paymentRequestCreator services.PaymentRequestCreator) services.PaymentRequestRecalculator {
	return &paymentRequestRecalculator{
		paymentRequestCreator: paymentRequestCreator,
	}
}

func (p *paymentRequestRecalculator) RecalculatePaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (*models.PaymentRequest, error) {
	var newPaymentRequest *models.PaymentRequest

	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		newPaymentRequest, err = p.doRecalculate(txnAppCtx, paymentRequestID)
		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return newPaymentRequest, nil
}

// doRecalculate handles the core recalculation process; put in separate method to make it easier to call from the transactional context.
func (p *paymentRequestRecalculator) doRecalculate(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) (*models.PaymentRequest, error) {
	// Fetch the payment request and payment service items from the old request.
	var oldPaymentRequest models.PaymentRequest
	err := appCtx.DB().
		EagerPreload(
			"PaymentServiceItems.MTOServiceItem.ReService",
			"ProofOfServiceDocs.PrimeUploads",
		).
		Find(&oldPaymentRequest, paymentRequestID)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(paymentRequestID, "for PaymentRequest")
		}
		return nil, services.NewQueryError("PaymentRequest", err, fmt.Sprintf("unexpected error while querying for payment request ID %s", paymentRequestID))
	}

	// Only pending payment requests can be recalculated.
	if oldPaymentRequest.Status != models.PaymentRequestStatusPending {
		return nil, services.NewConflictError(paymentRequestID, fmt.Sprintf("only pending payment requests can be recalculated, but this payment request has status of %s", oldPaymentRequest.Status))
	}

	// Re-create the payment request arg including service items, then call the create service (which should
	// price it with current inputs).
	inputPaymentRequest := buildPaymentRequestForRecalcuating(oldPaymentRequest)
	newPaymentRequest, err := p.paymentRequestCreator.CreatePaymentRequest(appCtx, &inputPaymentRequest)
	if err != nil {
		return nil, err // Just pass the error type from the PaymentRequestCreator.
	}

	// Set the (now) old payment request's status.
	// TODO: We need a better status for this -- something like "DEPRECATED".
	err = updateOldPaymentRequestStatus(appCtx, &oldPaymentRequest)
	if err != nil {
		return nil, err
	}

	// Duplicate the proof-of-service upload associations to the new payment request.
	err = associateProofOfServiceDocs(appCtx, oldPaymentRequest.ProofOfServiceDocs, newPaymentRequest)
	if err != nil {
		return nil, err
	}

	// Link the new payment request to the old one.
	err = linkNewToOldPaymentRequest(appCtx, newPaymentRequest, &oldPaymentRequest)
	if err != nil {
		return nil, err
	}

	return newPaymentRequest, nil
}

// buildPaymentRequestForRecalcuating builds up the expected payment request data based upon the old payment request.
func buildPaymentRequestForRecalcuating(oldPaymentRequest models.PaymentRequest) models.PaymentRequest {
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

func updateOldPaymentRequestStatus(appCtx appcontext.AppContext, oldPaymentRequest *models.PaymentRequest) error {
	newStatus := models.PaymentRequestStatusReviewedAllRejected
	oldPaymentRequest.Status = newStatus
	verrs, err := appCtx.DB().ValidateAndUpdate(oldPaymentRequest)
	if err != nil {
		return services.NewQueryError("PaymentRequest", err, fmt.Sprintf("failed to set old payment request status to %s for ID %s", newStatus, oldPaymentRequest.ID))
	}
	if verrs.HasAny() {
		return services.NewInvalidInputError(oldPaymentRequest.ID, err, verrs, fmt.Sprintf("failed to validate old payment request when setting status to %s", newStatus))
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
			return services.NewQueryError("ProofOfServiceDoc", err, fmt.Sprintf("failed to create proof of service doc for new payment request ID %s", newPaymentRequest.ID))
		}
		if verrs.HasAny() {
			return services.NewInvalidInputError(newPaymentRequest.ID, err, verrs, "failed to validate proof of service doc")
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
				return services.NewQueryError("PrimeUpload", err, fmt.Sprintf("failed to create prime upload for new payment request ID %s and new proof of service doc ID %s", newPaymentRequest.ID, newProofOfServiceDoc.ID))
			}
			if verrs.HasAny() {
				return services.NewInvalidInputError(newProofOfServiceDoc.ID, err, verrs, "failed to validate prime upload")
			}

			newProofOfServiceDoc.PrimeUploads = append(newProofOfServiceDoc.PrimeUploads, newPrimeUpload)
		}

		newPaymentRequest.ProofOfServiceDocs = append(newPaymentRequest.ProofOfServiceDocs, newProofOfServiceDoc)
	}

	return nil
}

// linkNewToOldPaymentRequest links a new payment request to the old payment request that was recalculated.
func linkNewToOldPaymentRequest(appCtx appcontext.AppContext, newPaymentRequest *models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
	newPaymentRequest.RepricedPaymentRequestID = &oldPaymentRequest.ID
	verrs, err := appCtx.DB().ValidateAndUpdate(newPaymentRequest)
	if err != nil {
		return services.NewQueryError("PaymentRequest", err, fmt.Sprintf("failed to link new payment request to old payment request ID %s", oldPaymentRequest.ID))
	}
	if verrs.HasAny() {
		return services.NewInvalidInputError(newPaymentRequest.ID, err, verrs, fmt.Sprintf("failed to validate new payment request when linking to old payment request ID %s", oldPaymentRequest.ID))
	}

	return nil
}
