package models_test

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestFetchAccessorials() {
	accessorial := testdatagen.MakeDummyAccessorial(suite.db)

	//Do
	accs, err := models.FetchAccessorials(suite.db)

	//Test
	suite.NoError(err)
	suite.Equal(1, len(accs))
	suite.Equal(accessorial.ID, accs[0].ID)
}
