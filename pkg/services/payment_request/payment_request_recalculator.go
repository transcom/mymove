package paymentrequest

import (
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestRecalculator struct {
	paymentRequestCreator       services.PaymentRequestCreator
	paymentRequestStatusUpdater services.PaymentRequestStatusUpdater
}

// NewPaymentRequestRecalculator returns a new payment request recalculator
func NewPaymentRequestRecalculator(paymentRequestCreator services.PaymentRequestCreator, paymentRequestStatusUpdater services.PaymentRequestStatusUpdater) services.PaymentRequestRecalculator {
	return &paymentRequestRecalculator{
		paymentRequestCreator:       paymentRequestCreator,
		paymentRequestStatusUpdater: paymentRequestStatusUpdater,
	}
}

func (p *paymentRequestRecalculator) RecalculatePaymentRequest(appCtx appcontext.AppContext, paymentRequestID uuid.UUID, startNewDBTx bool) (*models.PaymentRequest, error) {
	var newPaymentRequest *models.PaymentRequest

	if startNewDBTx {
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
	} else {
		var err error
		newPaymentRequest, err = p.doRecalculate(appCtx, paymentRequestID)
		return nil, err
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
	oldPaymentRequestEtag := etag.GenerateEtag(oldPaymentRequest.UpdatedAt)

	// Only pending payment requests can be recalculated.
	if oldPaymentRequest.Status != models.PaymentRequestStatusPending {
		return nil, services.NewConflictError(paymentRequestID, fmt.Sprintf("only pending payment requests can be recalculated, but this payment request has status of %s", oldPaymentRequest.Status))
	}

	// Set the (now) old payment request's status.  Doing this before we recalculate in case the create
	// payment request service needs to know that this will be deprecated if recalculation is successful.
	oldPaymentRequest.Status = models.PaymentRequestStatusDeprecated
	_, err = p.paymentRequestStatusUpdater.UpdatePaymentRequestStatus(appCtx, &oldPaymentRequest, oldPaymentRequestEtag)
	if err != nil {
		return nil, err
	}

	// Re-create the payment request arg including service items, then call the create service (which should
	// price it with current inputs).
	inputPaymentRequest := buildPaymentRequestForRecalculating(oldPaymentRequest)
	newPaymentRequest, err := p.paymentRequestCreator.CreatePaymentRequest(appCtx, &inputPaymentRequest)
	if err != nil {
		return nil, err // Just pass the error type from the PaymentRequestCreator.
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

// buildPaymentRequestForRecalculating builds up the expected payment request data based upon the old payment request.
func buildPaymentRequestForRecalculating(oldPaymentRequest models.PaymentRequest) models.PaymentRequest {
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
	newPaymentRequest.RecalculationOfPaymentRequestID = &oldPaymentRequest.ID
	verrs, err := appCtx.DB().ValidateAndUpdate(newPaymentRequest)
	if err != nil {
		return services.NewQueryError("PaymentRequest", err, fmt.Sprintf("failed to link new payment request to old payment request ID %s", oldPaymentRequest.ID))
	}
	if verrs.HasAny() {
		return services.NewInvalidInputError(newPaymentRequest.ID, err, verrs, fmt.Sprintf("failed to validate new payment request when linking to old payment request ID %s", oldPaymentRequest.ID))
	}

	return nil
}
