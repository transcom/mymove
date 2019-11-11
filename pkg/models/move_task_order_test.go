package models_test

import (
	"fmt"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGenerateReferenceID() {
	referenceID := models.GenerateReferenceID()
	fmt.Println("test ", referenceID)
	//suite.Equal(reflect.TypeOf(referenceID), "string")
}
