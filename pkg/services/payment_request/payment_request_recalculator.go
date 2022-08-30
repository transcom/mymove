package paymentrequest

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
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
			"ProofOfServiceDocs",
		).
		Find(&oldPaymentRequest, paymentRequestID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(paymentRequestID, "for PaymentRequest")
		default:
			return nil, apperror.NewQueryError("PaymentRequest", err, fmt.Sprintf("unexpected error while querying for payment request ID %s", paymentRequestID))
		}
	}
	oldPaymentRequestEtag := etag.GenerateEtag(oldPaymentRequest.UpdatedAt)

	// Only pending payment requests can be recalculated.
	if oldPaymentRequest.Status != models.PaymentRequestStatusPending {
		return nil, apperror.NewConflictError(paymentRequestID, fmt.Sprintf("only PENDING payment requests can be recalculated, but this payment request has status of %s", oldPaymentRequest.Status))
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
	inputPaymentRequest, err := buildPaymentRequestForRecalculating(appCtx, oldPaymentRequest)
	if err != nil {
		return nil, err
	}
	newPaymentRequest, err := p.paymentRequestCreator.CreatePaymentRequestCheck(appCtx, &inputPaymentRequest)
	if err != nil {
		return nil, err // Just pass the error type from the PaymentRequestCreator.
	}

	// Remap the proof-of-service upload associations to the new payment request.
	err = remapProofOfServiceDocs(appCtx, oldPaymentRequest.ProofOfServiceDocs, newPaymentRequest)
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
func buildPaymentRequestForRecalculating(appCtx appcontext.AppContext, oldPaymentRequest models.PaymentRequest) (models.PaymentRequest, error) {
	newPaymentRequest := models.PaymentRequest{
		IsFinal:         oldPaymentRequest.IsFinal,
		MoveTaskOrderID: oldPaymentRequest.MoveTaskOrderID,
	}

	// Get list of all service IDs that have payment-request based parameters (not many yet).
	paymentRequestServiceIDs, err := getServiceIDsWithPaymentRequestOrigin(appCtx)
	if err != nil {
		return models.PaymentRequest{}, err
	}

	var newPaymentServiceItems models.PaymentServiceItems
	for _, oldPaymentServiceItem := range oldPaymentRequest.PaymentServiceItems {
		newPaymentServiceItem := models.PaymentServiceItem{
			MTOServiceItemID: oldPaymentServiceItem.MTOServiceItemID,
			MTOServiceItem:   oldPaymentServiceItem.MTOServiceItem,
		}

		// If this service item is in the list of payment request service IDs, then we need to add the
		// params that came in via the payment request to the template payment request we're building.
		if uuidInSlice(paymentRequestServiceIDs, oldPaymentServiceItem.MTOServiceItem.ReServiceID) {
			newPaymentServiceItem.PaymentServiceItemParams, err = buildPaymentServiceItemParams(appCtx, oldPaymentServiceItem)
			if err != nil {
				return models.PaymentRequest{}, err
			}
		}

		newPaymentServiceItems = append(newPaymentServiceItems, newPaymentServiceItem)
	}

	sort.SliceStable(newPaymentServiceItems, func(i, j int) bool {
		return newPaymentServiceItems[i].MTOServiceItem.ReService.Priority < newPaymentServiceItems[j].MTOServiceItem.ReService.Priority
	})

	newPaymentRequest.PaymentServiceItems = newPaymentServiceItems

	return newPaymentRequest, nil
}

// getServiceIDsWithPaymentRequestOrigin returns all UUIDs of services that have at least one param with
// a payment request origin.
func getServiceIDsWithPaymentRequestOrigin(appCtx appcontext.AppContext) ([]uuid.UUID, error) {
	// Find all the service IDs that have params that originate from the payment request itself.
	var uuids []uuid.UUID
	query := `SELECT DISTINCT sp.service_id
		FROM service_params sp
		INNER JOIN service_item_param_keys sipk ON sp.service_item_param_key_id = sipk.id
		WHERE origin = ?`
	err := appCtx.DB().RawQuery(query, models.ServiceItemParamOriginPaymentRequest).All(&uuids)
	if err != nil {
		return nil, apperror.NewQueryError("ServiceParams", err, fmt.Sprintf("unexpected error while querying for service_params with params of origin %s", models.ServiceItemParamOriginPaymentRequest))
	}

	return uuids, nil
}

// uuidInSlice returns true if the given target UUID is in the slice of UUIDs; false otherwise.
func uuidInSlice(uuids []uuid.UUID, targetUUID uuid.UUID) bool {
	for _, uuid := range uuids {
		if uuid == targetUUID {
			return true
		}
	}

	return false
}

// buildPaymentServiceItemParams builds up any payment service items given as input to the original payment request.
func buildPaymentServiceItemParams(appCtx appcontext.AppContext, oldPaymentServiceItem models.PaymentServiceItem) (models.PaymentServiceItemParams, error) {
	// Get the payment service item param values for the specific service code.
	var paymentServiceItemParams models.PaymentServiceItemParams
	err := appCtx.DB().
		EagerPreload("ServiceItemParamKey").
		Join("service_item_param_keys sipk", "payment_service_item_params.service_item_param_key_id = sipk.id").
		Where("payment_service_item_id = ?", oldPaymentServiceItem.ID).
		Where("sipk.origin = ?", models.ServiceItemParamOriginPaymentRequest).
		All(&paymentServiceItemParams)
	if err != nil {
		return nil, apperror.NewQueryError("PaymentServiceItemParams", err, fmt.Sprintf("unexpected error while querying for payment service item params for payment service ID %s", oldPaymentServiceItem.ID))
	}

	// Create the incoming payment service item param for the new payment request we're going to be creating.
	var newPaymentServiceItemParams models.PaymentServiceItemParams
	for _, paymentServiceItemParam := range paymentServiceItemParams {
		newPaymentServiceItemParam := models.PaymentServiceItemParam{
			IncomingKey: paymentServiceItemParam.ServiceItemParamKey.Key.String(),
			Value:       paymentServiceItemParam.Value,
		}

		newPaymentServiceItemParams = append(newPaymentServiceItemParams, newPaymentServiceItemParam)
	}

	return newPaymentServiceItemParams, nil
}

// remapProofOfServiceDocs remaps a set of proof of service doc relationships to a new payment request.
func remapProofOfServiceDocs(appCtx appcontext.AppContext, proofOfServiceDocs models.ProofOfServiceDocs, newPaymentRequest *models.PaymentRequest) error {
	for _, proofOfServiceDoc := range proofOfServiceDocs {
		copyOfProofOfServiceDoc := proofOfServiceDoc // Make copy to avoid implicit memory aliasing of items from a range statement.
		copyOfProofOfServiceDoc.PaymentRequestID = newPaymentRequest.ID
		verrs, err := appCtx.DB().ValidateAndUpdate(&copyOfProofOfServiceDoc)
		if err != nil {
			return apperror.NewQueryError("ProofOfServiceDoc", err, fmt.Sprintf("failed to update proof of service doc for new payment request ID %s", newPaymentRequest.ID))
		}
		if verrs.HasAny() {
			return apperror.NewInvalidInputError(newPaymentRequest.ID, err, verrs, "failed to validate proof of service doc")
		}

		newPaymentRequest.ProofOfServiceDocs = append(newPaymentRequest.ProofOfServiceDocs, copyOfProofOfServiceDoc)
	}

	return nil
}

// linkNewToOldPaymentRequest links a new payment request to the old payment request that was recalculated.
func linkNewToOldPaymentRequest(appCtx appcontext.AppContext, newPaymentRequest *models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
	newPaymentRequest.RecalculationOfPaymentRequestID = &oldPaymentRequest.ID
	verrs, err := appCtx.DB().ValidateAndUpdate(newPaymentRequest)
	if err != nil {
		return apperror.NewQueryError("PaymentRequest", err, fmt.Sprintf("failed to link new payment request to old payment request ID %s", oldPaymentRequest.ID))
	}
	if verrs.HasAny() {
		return apperror.NewInvalidInputError(newPaymentRequest.ID, err, verrs, fmt.Sprintf("failed to validate new payment request when linking to old payment request ID %s", oldPaymentRequest.ID))
	}

	return nil
}
