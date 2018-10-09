package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchAccessorials() {
	accessorial := testdatagen.MakeDefaultTariff400ngItem(suite.db)

	//Do
	accs, err := models.FetchTariff400ngItems(suite.db, false)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(accessorial.ID, accs[0].ID)
}
