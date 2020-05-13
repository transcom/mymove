package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentServiceItemParam creates a single PaymentServiceItemParam and associated relationships
func MakePaymentServiceItemParam(db *pop.Connection, assertions Assertions) models.PaymentServiceItemParam {
	paymentServiceItem := assertions.PaymentServiceItem
	if isZeroUUID(paymentServiceItem.ID) {
		paymentServiceItem = MakePaymentServiceItem(db, assertions)
	}

	serviceItemParamKey := assertions.ServiceItemParamKey
	if isZeroUUID(serviceItemParamKey.ID) {
		serviceItemParamKey = FetchOrMakeServiceItemParamKey(db, assertions)
	}

	paymentServiceItemParam := models.PaymentServiceItemParam{
		PaymentServiceItem:    paymentServiceItem,
		PaymentServiceItemID:  paymentServiceItem.ID,
		ServiceItemParamKey:   serviceItemParamKey,
		ServiceItemParamKeyID: serviceItemParamKey.ID,
		Value:                 "123",
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentServiceItemParam, assertions.PaymentServiceItemParam)

	mustCreate(db, &paymentServiceItemParam)

	return paymentServiceItemParam
}

// MakeDefaultPaymentServiceItemParam makes a PaymentServiceItemParam with default values
func MakeDefaultPaymentServiceItemParam(db *pop.Connection) models.PaymentServiceItemParam {
	return MakePaymentServiceItemParam(db, Assertions{})
}
