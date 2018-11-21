package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchTariff400ngItems() {
	tariff400ngItem := testdatagen.MakeDefaultTariff400ngItem(suite.db)

	//Do
	accs, err := models.FetchTariff400ngItems(suite.db, false)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(tariff400ngItem.ID, accs[0].ID)
}

func (suite *ModelSuite) TestFetchTariff400ngItemByCode() {
	tariff400ngItem := testdatagen.MakeDefaultTariff400ngItem(suite.db)

	fetchedItem, err := models.FetchTariff400ngItemByCode(suite.db, tariff400ngItem.Code)

	//Test
	suite.NoError(err)
	suite.Equal(tariff400ngItem.ID, fetchedItem.ID)
}
