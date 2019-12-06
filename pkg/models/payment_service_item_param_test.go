package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPaymentServiceItemParamValidation() {
	suite.T().Run("test valid PaymentServiceItemParam", func(t *testing.T) {
		validPaymentServiceItemParam := models.PaymentServiceItemParam{
			PaymentServiceItemID:  uuid.Must(uuid.NewV4()),
			ServiceItemParamKeyID: uuid.Must(uuid.NewV4()),
			Value:                 "Value",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPaymentServiceItemParam, expErrors)
	})

	suite.T().Run("test empty PaymentServiceItemParam", func(t *testing.T) {
		invalidPaymentServiceItemParam := models.PaymentServiceItemParam{}

		expErrors := map[string][]string{
			"payment_service_item_id":   {"PaymentServiceItemID can not be blank."},
			"service_item_param_key_id": {"ServiceItemParamKeyID can not be blank."},
			"value":                     {"Value can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPaymentServiceItemParam, expErrors)
	})
}