package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMTOServiceItemValidation() {
	suite.Run("test valid MTOServiceItem", func() {
		moveTaskOrderID := uuid.Must(uuid.NewV4())
		mtoShipmentID := uuid.Must(uuid.NewV4())
		reServiceID := uuid.Must(uuid.NewV4())

		validMTOServiceItem := models.MTOServiceItem{
			MoveTaskOrderID: moveTaskOrderID,
			MTOShipmentID:   &mtoShipmentID,
			ReServiceID:     reServiceID,
			Status:          models.MTOServiceItemStatusSubmitted,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOServiceItem, expErrors)
	})
}
