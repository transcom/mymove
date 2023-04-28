package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// CreatePaymentServiceItemParams helper struct to facilitate creation of multiple payment service item params
type CreatePaymentServiceItemParams struct {
	Key     models.ServiceItemParamName
	KeyType models.ServiceItemParamType
	Value   string
}

func BuildPaymentServiceItemParam(db *pop.Connection, customs []Customization, traits []Trait) models.PaymentServiceItemParam {
	customs = setupCustomizations(customs, traits)

	// Find PaymentServiceItemParam customization and convert to models.PaymentServiceItemParam
	var cPaymentServiceItemParam models.PaymentServiceItemParam
	if result := findValidCustomization(customs, PaymentServiceItemParam); result != nil {
		cPaymentServiceItemParam = result.Model.(models.PaymentServiceItemParam)
		if result.LinkOnly {
			return cPaymentServiceItemParam
		}
	}

	paymentServiceItem := BuildPaymentServiceItem(db, customs, traits)

	serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, customs, traits)

	// Create PaymentServiceItemParam
	paymentServiceItemParam := models.PaymentServiceItemParam{
		PaymentServiceItem:    paymentServiceItem,
		PaymentServiceItemID:  paymentServiceItem.ID,
		ServiceItemParamKey:   serviceItemParamKey,
		ServiceItemParamKeyID: serviceItemParamKey.ID,
		Value:                 "123",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&paymentServiceItemParam, cPaymentServiceItemParam)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &paymentServiceItemParam)
	}
	return paymentServiceItemParam

}

func BuildPaymentServiceItemWithParams(db *pop.Connection, serviceCode models.ReServiceCode, paramsToCreate []CreatePaymentServiceItemParams, customs []Customization, traits []Trait) models.PaymentServiceItem {
	var params models.PaymentServiceItemParams

	// Make customizations for PaymentServiceItem
	paymentServiceItemCustoms := customs
	paymentServiceItemCustoms = append(paymentServiceItemCustoms, Customization{
		Model: models.ReService{
			Code: serviceCode,
		},
	})
	paymentServiceItem := BuildPaymentServiceItem(db, paymentServiceItemCustoms, traits)

	for _, param := range paramsToCreate {
		// Make customizations for ServiceItemParamKey
		serviceItemCustoms := customs
		serviceItemCustoms = append(serviceItemCustoms, Customization{
			Model: models.ServiceItemParamKey{
				Key:  param.Key,
				Type: param.KeyType,
			},
		})
		serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, serviceItemCustoms, traits)

		// Make customizations for PaymentServiceItemParam
		newPaymentServiceItemParamCustoms := []Customization{
			{
				Model:    serviceItemParamKey,
				LinkOnly: true,
			},
			{
				Model: models.PaymentServiceItemParam{
					Value: param.Value,
				},
			},
			{
				Model:    paymentServiceItem,
				LinkOnly: true,
			},
		}

		paymentServiceItemParam := BuildPaymentServiceItemParam(db, newPaymentServiceItemParamCustoms, traits)
		params = append(params, paymentServiceItemParam)
	}

	paymentServiceItem.PaymentServiceItemParams = params

	return paymentServiceItem
}
