package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMoveTaskOrderValidation() {
	suite.T().Run("test valid MoveTaskOrder", func(t *testing.T) {
		validMoveTaskOrder := models.MoveTaskOrder{
			MoveOrderID:  uuid.Must(uuid.NewV4()),
			ReferenceID:  "Testing",
			ContractorID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMoveTaskOrder, expErrors)
	})
}
