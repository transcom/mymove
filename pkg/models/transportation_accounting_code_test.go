package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ValidTac() {
	tac := models.TransportationAccountingCode{
		ID:  uuid.Must(uuid.NewV4()),
		TAC: "TheTac",
	}

	expErrors := map[string][]string{}

	suite.verifyValidationErrors(&tac, expErrors)
}

func (suite *ModelSuite) Test_InvalidTac() {
	tac := models.TransportationAccountingCode{
		ID: uuid.Must(uuid.NewV4()),
	}

	expErrors := map[string][]string{
		"tac": {"TAC can not be blank."},
	}

	suite.verifyValidationErrors(&tac, expErrors)
}
