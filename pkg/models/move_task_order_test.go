package models_test

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGenerateReferenceID() {
	referenceID := models.GenerateReferenceID()
	firstNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[0])
	secondNum, _ := strconv.Atoi(strings.Split(referenceID, "-")[1])
	suite.Equal(reflect.TypeOf(referenceID).String(), "string")
	suite.Equal(firstNum >= 0 && firstNum <= 9999, true)
	suite.Equal(secondNum >= 0 && secondNum <= 9999, true)
	suite.Equal(string(referenceID[4]), "-")
}
