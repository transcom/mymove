package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_AllFieldsOptionalLoa() {
	loa := &models.LineOfAccounting{
		ID: uuid.Must(uuid.NewV4()),
	}

	err := suite.DB().Save(loa)
	suite.NoError(err)
}
