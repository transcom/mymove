package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
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

// ShipmentRecalculatePaymentRequest recalculate all PENDING payment requests for shipmentID if
// the payment request has any service items that allow for reweigh or weight adjustments
func (p *paymentRequestShipmentRecalculator) ShipmentRecalculatePaymentRequest(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.PaymentRequests, error) {
	var recalculatedPaymentRequests models.PaymentRequests
	// Find all applicable payment request for the shipment
	paymentRequestIDs, err := findPendingPaymentRequestsForShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	// Note: Payment requests can have MTO Service Items from different shipments. RecalculatePaymentRequest
	// will reprice the whole payment request if the payment request has a service item that needs to be recalculated.

	// Recalculate the payment request
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		for _, pr := range paymentRequestIDs {
			var newPR *models.PaymentRequest
			newPR, err = p.paymentRequestRecalculator.RecalculatePaymentRequest(txnAppCtx, pr)
			if err != nil {
				return err
			}
			if newPR != nil {
				recalculatedPaymentRequests = append(recalculatedPaymentRequests, *newPR)
			}
		}
		return nil
	})
	if transactionError != nil {
		return nil, transactionError
	}

	return &recalculatedPaymentRequests, nil
}

func findPendingPaymentRequestsForShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) ([]uuid.UUID, error) {

	// Given a shipmentID find all of the payment requests in PENDING
	// that have service items that can have WeightReweigh or WeightAdjusted params.
	var uuids []uuid.UUID
	query := `SELECT DISTINCT pr.id
		FROM payment_requests pr
		INNER JOIN mto_shipments ms ON ms.move_id = pr.move_id
		INNER JOIN payment_service_items psi ON pr.id = psi.payment_request_id
		INNER JOIN mto_service_items msi ON ms.id = msi.mto_shipment_id AND psi.mto_service_item_id = msi.id
		INNER JOIN re_services rs ON msi.re_service_id = rs.id
		INNER JOIN service_params sp ON rs.id = sp.service_id
		INNER JOIN service_item_param_keys sipk ON sp.service_item_param_key_id = sipk.id
		WHERE pr.status = $1 AND ms.id = $2 AND sipk.key in ($3, $4);`
	err := appCtx.DB().RawQuery(query,
		models.PaymentRequestStatusPending,
		shipmentID,
		models.ServiceItemParamNameWeightReweigh, models.ServiceItemParamNameWeightAdjusted,
	).All(&uuids)

	return uuids, err
}
