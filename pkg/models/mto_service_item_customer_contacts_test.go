package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOServiceItemCustomerContactValidation() {
	suite.Run("test valid MTOServiceItemCustomerContact", func() {
		validMTOServiceItemDimension := models.MTOServiceItemCustomerContact{
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              time.Now().Add(time.Hour * 24),
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors, nil)
	})

	suite.Run("test invalid MTOServiceItemCustomerContact", func() {
		validMTOServiceItemDimension := models.MTOServiceItemCustomerContact{
			Type:                       "NOT VALID",
			DateOfContact:              time.Time{},
			TimeMilitary:               "",
			FirstAvailableDeliveryDate: time.Time{},
		}
		expErrors := map[string][]string{
			"type":                          {"Type is not in the list [FIRST, SECOND]."},
			"date_of_contact":               {"DateOfContact can not be blank."},
			"time_military":                 {"TimeMilitary can not be blank."},
			"first_available_delivery_date": {"FirstAvailableDeliveryDate can not be blank."},
		}
		suite.verifyValidationErrors(&validMTOServiceItemDimension, expErrors, nil)
	})
}
