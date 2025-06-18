package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceRequestDocumentValidation() {
	suite.Run("test valid ServiceRequestDocument", func() {
		validServiceRequestDocument := models.ServiceRequestDocument{
			MTOServiceItemID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceRequestDocument, expErrors, nil)
	})

	suite.Run("test empty ServiceRequestDocument", func() {
		invalidServiceRequestDocument := models.ServiceRequestDocument{}

		expErrors := map[string][]string{
			"mtoservice_item_id": {"MTOServiceItemID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidServiceRequestDocument, expErrors, nil)
	})
}
