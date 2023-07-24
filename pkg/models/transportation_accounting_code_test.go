package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ValidateTac() {
	tac := models.TransportationAccountingCode{
		ID:  uuid.Must(uuid.NewV4()),
		TAC: "TheTac",
	}

	expErrors := map[string][]string{}

	suite.verifyValidationErrors(&tac, expErrors)
}
