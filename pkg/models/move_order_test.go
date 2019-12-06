package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMoveOrderValidation() {
	suite.T().Run("test valid MoveOrder", func(t *testing.T) {
		validMoveOrder := models.MoveOrder{
			CustomerID:               uuid.Must(uuid.NewV4()),
			EntitlementID:            uuid.Must(uuid.NewV4()),
			DestinationDutyStationID: uuid.Must(uuid.NewV4()),
			OriginDutyStationID:      uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMoveOrder, expErrors)
	})
}
