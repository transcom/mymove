package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestShipmentRecalculator struct {
	paymentRequestRecalculator services.PaymentRequestRecalculator
}

// NewPaymentRequestShipmentRecalculator returns a new payment request recalculator for a shipment
func NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator services.PaymentRequestRecalculator) services.PaymentRequestShipmentRecalculator {
	return &paymentRequestShipmentRecalculator{
		paymentRequestRecalculator: paymentRequestRecalculator,
	}
}

// ShipmentRecalculatePaymentRequest recalculate all PENDING payment requests given a shipment ID
func (p *paymentRequestShipmentRecalculator) ShipmentRecalculatePaymentRequest(appCtx appcontext.AppContext, shipmentID uuid.UUID) error {

	// Given a shipmentID find all of the payment requests in PENDING.
	var paymentRequests []models.PaymentRequest
	/*
		option 1: filter down to shipment ID and PENDING and write code to
		   determine if weight param is present.
		----
			select * from public.payment_requests pr
			left join mto_shipments ms on ms.move_id = pr.move_id
	*/
	/*
		option 2: filter down to shipment ID, PENDING, and if WeightOriginal is
			present in the payment request. Don't need Eager because the query has
		   found the PRs that we need here.
		---
		select * from public.payment_requests pr
		left join mto_shipments ms on ms.move_id = pr.move_id
		left join payment_service_items psi on pr.id = psi.payment_request_id
		left join payment_service_item_params psip on psi.id = psip.payment_service_item_id
		left join service_item_param_keys sipk on psip.service_item_param_key_id = sipk.id
		where pr.status = 'PENDING' and sipk.key = 'WeightOriginal';
	*/
	err := appCtx.DB(). /*EagerPreload(
		"PaymentServiceItems.MTOServiceItem.MTOShipment",
		"PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"ProofOfServiceDocs").*/Q().
		Join("mto_shipments ms", "ms.move_id = payment_requests.move_id").
		Join("payment_service_items psi", "payment_requests.id = psi.payment_request_id").
		Join("payment_service_item_params psip", "psi.id = psip.payment_service_item_id").
		Join("service_item_param_keys sipk", "psip.service_item_param_key_id = sipk.id").
		Where("ms.id = $1", shipmentID).
		Where("payment_requests.status = $2", "PENDING").
		Where("sipk.key = $3", "WeightOriginal").
		All(&paymentRequests)
	if err != nil {
		return err
	}

	// var newPR *models.PaymentRequest
	startNewTx := false
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, pr := range paymentRequests {
			_ /*newPR*/, err = p.paymentRequestRecalculator.RecalculatePaymentRequest(txnAppCtx, pr.ID, startNewTx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if transactionError != nil {
		return transactionError
	}
	return nil
}
