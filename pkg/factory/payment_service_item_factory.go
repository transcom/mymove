package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func BuildPaymentServiceItem(db *pop.Connection, customs []Customization, traits []Trait) models.PaymentServiceItem {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cPaymentServiceItem models.PaymentServiceItem
	if result := findValidCustomization(customs, PaymentServiceItem); result != nil {
		cPaymentServiceItem = result.Model.(models.PaymentServiceItem)
		if result.LinkOnly {
			return cPaymentServiceItem
		}
	}

	paymentRequest := BuildPaymentRequest(db, customs, traits)
	mtoServiceItem := BuildMTOServiceItem(db, customs, traits)

	// Create default PaymentServiceItem
	var cents = unit.Cents(888)
	paymentServiceItem := models.PaymentServiceItem{
		PaymentRequest:   paymentRequest,
		PaymentRequestID: paymentRequest.ID,
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
		PriceCents:       &cents,
		Status:           models.PaymentServiceItemStatusRequested,
		RequestedAt:      time.Now(),
	}
	// Overwrite values with those from customizations
	testdatagen.MergeModels(&paymentServiceItem, cPaymentServiceItem)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &paymentServiceItem)
	}

	return paymentServiceItem
}
