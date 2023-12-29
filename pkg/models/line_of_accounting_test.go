package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_LineOfAccountingAllFieldsOptionalCanSave() {
	loa := factory.BuildDefaultLineOfAccounting(suite.DB())

	suite.MustSave(&loa)
}

func (suite *ModelSuite) Test_LineOfAccountingAllFieldsPresentCanSave() {
	loa := factory.BuildFullLineOfAccounting(suite.DB(), nil, nil)

	suite.MustSave(&loa)
}

func (suite *ModelSuite) Test_LineOfAccountingCanSaveAndFetch() {
	// Can save
	loa := models.LineOfAccounting{LoaSysID: models.StringPointer("1234")}

	suite.MustCreate(&loa)

	// Can fetch
	var fetchedLoa models.LineOfAccounting
	err := suite.DB().Where("loa_sys_id = $1", *loa.LoaSysID).First(&fetchedLoa)

	suite.NoError(err)
	suite.Equal(*loa.LoaSysID, *fetchedLoa.LoaSysID)
}
