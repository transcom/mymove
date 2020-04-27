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

	suite.T().Run("test invalid MoveTaskOrder", func(t *testing.T) {
		invalidMoveTaskOrder := models.MoveTaskOrder{}
		expErrors := map[string][]string{
			"move_order_id": {"MoveOrderID can not be blank."},
			"reference_id":  {"ReferenceID can not be blank."},
			"contractor_id": {"ContractorID can not be blank."},
		}
		suite.verifyValidationErrors(&invalidMoveTaskOrder, expErrors)
	})
}
