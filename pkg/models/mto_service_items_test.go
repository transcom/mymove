package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMtoServiceItemValidation() {
	suite.T().Run("test valid MtoServiceItem", func(t *testing.T) {
		validMtoServiceItem := models.MtoServiceItem{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			MtoShipmentID:   uuid.Must(uuid.NewV4()),
			ReServiceID:     uuid.Must(uuid.NewV4()),
			MetaID:          uuid.Must(uuid.NewV4()),
			MetaType:        "unknown",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMtoServiceItem, expErrors)
	})
}
