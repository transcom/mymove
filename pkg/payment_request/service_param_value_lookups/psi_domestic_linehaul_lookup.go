package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/transcom/mymove/pkg/models"
)

// PSILinehaulDomLookup does lookup of uuid of payment service item for dom linehaul
type PSILinehaulDomLookup struct {
	MTOShipment models.MTOShipment
}

// PSILinehaulDomPriceLookup does lookup of price in cents of payment service item for dom linehaul
type PSILinehaulDomPriceLookup struct {
	MTOShipment models.MTOShipment
}

func getPaymentServiceItem(keyData *ServiceItemParamKeyData, mtoShipment models.MTOShipment) (models.PaymentServiceItem, error) {
	db := *keyData.db

	paymentRequestID := keyData.PaymentRequestID

	// mtoServiceItemID -> mtoShipmentId + mtoId -> find a service where mtoId==, mtoShipmentId==, and reServiceid==DLH
	var paymentServiceItemDLH models.PaymentServiceItem
	err := db.Q().
		Join("mto_service_items msi", "msi.id = payment_service_items.mto_service_item_id").
		Join("re_services rs", "rs.id = msi.re_service_id").
		Where("payment_service_items.status != $1", models.PaymentServiceItemStatusDenied).
		Where("msi.mto_shipment_id = $2", mtoShipment.ID).
		Where("rs.code = $3", models.ReServiceCodeDLH).
		Last(&paymentServiceItemDLH) // pop Last orders by created_at

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentServiceItem{}, fmt.Errorf("couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", paymentRequestID, keyData.MTOServiceItem.ID)
		default:
			return models.PaymentServiceItem{}, err
		}
	}

	return paymentServiceItemDLH, nil
}

func (r PSILinehaulDomLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	paymentServiceItem, err := getPaymentServiceItem(keyData, r.MTOShipment)
	if err != nil {
		return "", err
	}

	return paymentServiceItem.ID.String(), nil
}

func (r PSILinehaulDomPriceLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	paymentServiceItem, err := getPaymentServiceItem(keyData, r.MTOShipment)
	if err != nil {
		return "", err
	}

	if paymentServiceItem.PriceCents == nil {
		return "", fmt.Errorf("found PaymentServiceItem for dom linehaul but it has no price! id found: %s", paymentServiceItem.ID)
	}

	return paymentServiceItem.PriceCents.String(), nil
}
