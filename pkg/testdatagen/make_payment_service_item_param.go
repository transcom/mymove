package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// CreatePaymentServiceItemParams helper struct to facilitate creation of multiple payment service item params
type CreatePaymentServiceItemParams struct {
	Key     models.ServiceItemParamName
	KeyType models.ServiceItemParamType
	Value   string
}

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

// MakeMultiplePaymentServiceItemParams creates more than one payment service item param at a time
func MakeMultiplePaymentServiceItemParams(db *pop.Connection, serviceCode models.ReServiceCode, paramsToCreate []CreatePaymentServiceItemParams) models.PaymentServiceItem {
	var params models.PaymentServiceItemParams

	paymentServiceItem := MakePaymentServiceItem(db, Assertions{
		ReService: models.ReService{
			Code: serviceCode,
		},
	})

	for _, param := range paramsToCreate {
		serviceItemParamKey := FetchOrMakeServiceItemParamKey(db,
			Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  param.Key,
					Type: param.KeyType,
				},
			})

		serviceItemParam := MakePaymentServiceItemParam(db,
			Assertions{
				PaymentServiceItem:  paymentServiceItem,
				ServiceItemParamKey: serviceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: param.Value,
				},
			})
		params = append(params, serviceItemParam)
	}

	paymentServiceItem.PaymentServiceItemParams = params

	return paymentServiceItem
}

// MakeDefaultPaymentServiceItemParam makes a PaymentServiceItemParam with default values
func MakeDefaultPaymentServiceItemParam(db *pop.Connection) models.PaymentServiceItemParam {
	return MakePaymentServiceItemParam(db, Assertions{})
}
