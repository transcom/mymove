package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func newUUIDPointer() *uuid.UUID {
	u := uuid.Must(uuid.NewV4())
	return &u
}

func (suite *ModelSuite) TestMoveOrderValidation() {
	suite.T().Run("test valid MoveOrder", func(t *testing.T) {
		validMoveOrder := models.MoveOrder{
			CustomerID:               newUUIDPointer(),
			EntitlementID:            newUUIDPointer(),
			DestinationDutyStationID: newUUIDPointer(),
			OriginDutyStationID:      newUUIDPointer(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMoveOrder, expErrors)
	})
}
