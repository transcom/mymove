package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOServiceItemCustomerContactValidation() {
	suite.Run("test valid MTOServiceItemCustomerContact", func() {
		mtoServiceItemDimensionID := uuid.Must(uuid.NewV4())

		validMTOServiceItemDimension := models.MTOServiceItemCustomerContact{
			MTOServiceItemID:           mtoServiceItemDimensionID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors)
	})

	suite.Run("test invalid MTOServiceItemCustomerContact", func() {
		validMTOServiceItemDimension := models.MTOServiceItemCustomerContact{
			MTOServiceItemID:           uuid.Nil,
			Type:                       "NOT VALID",
			TimeMilitary:               "",
			FirstAvailableDeliveryDate: time.Time{},
		}
		expErrors := map[string][]string{
			"mtoservice_item_id":            {"MTOServiceItemID can not be blank."},
			"type":                          {"Type is not in the list [FIRST, SECOND]."},
			"time_military":                 {"TimeMilitary can not be blank."},
			"first_available_delivery_date": {"FirstAvailableDeliveryDate can not be blank."},
		}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors)
	})
}
