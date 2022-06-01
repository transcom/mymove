package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

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

	mustCreate(db, &paymentServiceItemParam, assertions.Stub)

	return paymentServiceItemParam
}

// MakeDefaultPaymentServiceItemParam makes a PaymentServiceItemParam with default values
func MakeDefaultPaymentServiceItemParam(db *pop.Connection) models.PaymentServiceItemParam {
	return MakePaymentServiceItemParam(db, Assertions{})
}

// MakePaymentServiceItemWithParams creates more than one payment service item param at a time
func MakePaymentServiceItemWithParams(db *pop.Connection, serviceCode models.ReServiceCode, paramsToCreate []CreatePaymentServiceItemParams, assertions Assertions) models.PaymentServiceItem {
	var params models.PaymentServiceItemParams

	assertions.ReService = models.ReService{
		Code: serviceCode,
	}
	paymentServiceItem := MakePaymentServiceItem(db, assertions)

	for _, param := range paramsToCreate {
		assertions.ServiceItemParamKey = models.ServiceItemParamKey{
			Key:  param.Key,
			Type: param.KeyType,
		}
		serviceItemParamKey := FetchOrMakeServiceItemParamKey(db, assertions)

		assertions.PaymentServiceItem = paymentServiceItem
		assertions.ServiceItemParamKey = serviceItemParamKey
		assertions.PaymentServiceItemParam = models.PaymentServiceItemParam{
			Value: param.Value,
		}
		serviceItemParam := MakePaymentServiceItemParam(db, assertions)
		params = append(params, serviceItemParam)
	}

	paymentServiceItem.PaymentServiceItemParams = params

	return paymentServiceItem
}

// MakeDefaultPaymentServiceItemWithParams creates more than one payment service item param at a time with default values
func MakeDefaultPaymentServiceItemWithParams(db *pop.Connection, serviceCode models.ReServiceCode, paramsToCreate []CreatePaymentServiceItemParams) models.PaymentServiceItem {
	return MakePaymentServiceItemWithParams(db, serviceCode, paramsToCreate, Assertions{})
}
