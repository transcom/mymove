package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_AllFieldsOptionalCanSave() {
	loa := factory.BuildDefaultLineOfAccounting(suite.DB())

	suite.MustCreate(&loa)
}

func (suite *ModelSuite) Test_AllFieldsPresentCanSave() {
	loa := factory.BuildFullLineOfAccounting(suite.DB())

	suite.MustCreate(loa)
}

func (suite *ModelSuite) Test_CanSaveAndFetch() {
	// Can save
	loa := models.LineOfAccounting{LoaSysID: models.IntPointer(1234)}

	suite.MustCreate(&loa)

	// Can fetch
	var fetchedLoa models.LineOfAccounting
	err := suite.DB().Where("loa_sys_id = $1", *loa.LoaSysID).First(&fetchedLoa)

	suite.NoError(err)
	suite.Equal(*loa.LoaSysID, *fetchedLoa.LoaSysID)
}
