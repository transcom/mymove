package serviceparamvaluelookups

import (
	"fmt"

	"database/sql"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// PSILinehaulDomLookup does lookup of uuid of payment service item for dom linehaul
type PSILinehaulDomLookup struct {
}

// PSILinehaulDomPriceLookup does lookup of price in cents of payment service item for dom linehaul
type PSILinehaulDomPriceLookup struct {
}

func getPaymentServiceItem(keyData *ServiceItemParamKeyData) (models.PaymentServiceItem, *models.MTOServiceItem, error) {
	db := *keyData.db

	paymentRequestID := keyData.PaymentRequestID
	mtoServiceItemID := keyData.MTOServiceItemID

	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService").Where("id = ?", mtoServiceItemID).First(&mtoServiceItem)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentServiceItem{}, nil, services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return models.PaymentServiceItem{}, nil, err
		}
	}

	// if we're not pricing DLH, we shouldn't attempt to find a PSI for it
	if mtoServiceItem.ReService.Code != models.ReServiceCodeDLH {
		return models.PaymentServiceItem{}, &mtoServiceItem, nil
	}

	// mtoServiceItemID -> mtoShipmentId + mtoId -> find a service where mtoId==, mtoShipmentId==, and reServiceid==DLH
	var paymentServiceItemDLH models.PaymentServiceItem
	err = db.Q().
		Join("mto_service_items msi", "msi.id = payment_service_items.mto_service_item_id").
		Join("re_services rs", "rs.id = msi.re_service_id").
		Where("payment_service_items.status != $1", models.PaymentServiceItemStatusDenied).
		Where("msi.mto_shipment_id = $2", mtoServiceItem.MTOShipmentID).
		Where("rs.code = $3", models.ReServiceCodeDLH).
		Last(&paymentServiceItemDLH) // pop Last orders by created_at

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.PaymentServiceItem{}, &mtoServiceItem, fmt.Errorf("couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", paymentRequestID, mtoServiceItemID)
		default:
			return models.PaymentServiceItem{}, &mtoServiceItem, err
		}
	}

	return paymentServiceItemDLH, &mtoServiceItem, nil
}

func (r PSILinehaulDomLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	paymentServiceItem, _, err := getPaymentServiceItem(keyData)
	if err != nil {
		return "", err
	}

	return paymentServiceItem.ID.String(), nil
}

func (r PSILinehaulDomPriceLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	paymentServiceItem, mtoServiceItem, err := getPaymentServiceItem(keyData)
	if err != nil {
		return "", err
	}

	// if we're not pricing DLH, we shouldn't attempt to find a PSI for it
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDLH {
		if paymentServiceItem.PriceCents == nil {
			return "", fmt.Errorf("found PaymentServiceItem for dom linehaul but it has no price! id found: %s", paymentServiceItem.ID)
		}

		return paymentServiceItem.PriceCents.String(), nil
	}

	return "0", nil
}
