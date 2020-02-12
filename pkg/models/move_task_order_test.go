package models_test

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestMoveTaskOrderValidation() {
	suite.T().Run("test valid MoveTaskOrder", func(t *testing.T) {
		validMoveTaskOrder := models.MoveTaskOrder{
			MoveOrderID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{
			"referenceID": {"ReferenceID cannot be blank."},
		}
		suite.verifyValidationErrors(&validMoveTaskOrder, expErrors)
	})
}

func (suite *ModelSuite) TestGenerateReferenceID() {
	r, err := models.GenerateReferenceID(suite.DB())
	suite.NotNil(r)
	referenceID := r
	suite.NoError(err)
	firstNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[0])
	secondNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[1])
	suite.Equal(reflect.TypeOf(referenceID).String(), "string")
	suite.Equal(firstNum >= 0 && firstNum <= 9999, true)
	suite.Equal(secondNum >= 0 && secondNum <= 9999, true)
	suite.Equal(string(referenceID[4]), "-")
}
