package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOServiceItemValidation() {
	suite.T().Run("test valid MTOServiceItem", func(t *testing.T) {
		moveTaskOrderID := uuid.Must(uuid.NewV4())
		mtoShipmentID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())

		validMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &mtoShipmentID,
			ReServiceID:     reServiceID,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItem, expErrors)
	})
}

func (suite *ModelSuite) TestMTOServiceItemGetDimension() {
	mtoServiceItem := models.MTOServiceItem{
		Dimensions: models.MTOServiceItemDimensions{
			models.MTOServiceItemDimension{
				Type: models.DimensionTypeItem,
			},
			models.MTOServiceItemDimension{
				Type: models.DimensionTypeCrate,
			},
		},
	}

	suite.T().Run("test valid ITEM dimension exists", func(t *testing.T) {
		dimension := mtoServiceItem.GetItemDimension()
		suite.IsType(models.DimensionTypeItem, dimension.Type)
	})

	suite.T().Run("test valid CRATE dimension exists", func(t *testing.T) {
		dimension := mtoServiceItem.GetCrateDimension()
		suite.IsType(models.DimensionTypeCrate, dimension.Type)
	})

	suite.T().Run("test return nil if list is empty", func(t *testing.T) {
		mtoServiceItem = models.MTOServiceItem{}
		suite.Nil(mtoServiceItem.GetItemDimension())
		suite.Nil(mtoServiceItem.GetCrateDimension())
	})
}

func (suite *ModelSuite) TestMTOServiceItemGetCustomerContact() {
	mtoServiceItem := models.MTOServiceItem{
		CustomerContacts: models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type: models.CustomerContactTypeFirst,
			},
			models.MTOServiceItemCustomerContact{
				Type: models.CustomerContactTypeSecond,
			},
		},
	}

	suite.T().Run("test valid first customer contact exists", func(t *testing.T) {
		customerContact := mtoServiceItem.GetFirstCustomerContact()
		suite.IsType(models.CustomerContactTypeFirst, customerContact.Type)
	})

	suite.T().Run("test valid second customer contact exists", func(t *testing.T) {
		customerContact := mtoServiceItem.GetSecondCustomerContact()
		suite.IsType(models.CustomerContactTypeSecond, customerContact.Type)
	})

	suite.T().Run("test return nil if list is empty", func(t *testing.T) {
		mtoServiceItem = models.MTOServiceItem{}
		suite.Nil(mtoServiceItem.GetFirstCustomerContact())
		suite.Nil(mtoServiceItem.GetSecondCustomerContact())
	})
}
