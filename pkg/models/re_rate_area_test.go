package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReRateAreaValidation() {
	suite.Run("test valid ReRateArea", func() {
		validReRateArea := models.ReRateArea{
			ContractID: uuid.Must(uuid.NewV4()),
			IsOconus:   true,
			Code:       "123abc",
			Name:       "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReRateArea, expErrors, nil)
	})

	suite.Run("test empty ReRateArea", func() {
		emptyReRateArea := models.ReRateArea{}
		expErrors := map[string][]string{
			"contract_id": {"ContractID can not be blank."},
			"code":        {"Code can not be blank."},
			"name":        {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReRateArea, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchRateAreaID() {
	suite.Run("success - fetching a rate area ID", func() {
		service := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		address := factory.BuildAddress(suite.DB(), nil, nil)
		rateAreaId, err := models.FetchRateAreaID(suite.DB(), address.ID, &service.ID, contract.ID)
		suite.NotNil(rateAreaId)
		suite.NoError(err)
	})

	suite.Run("fail - receive error when not all values are provided", func() {
		var nilUuid uuid.UUID
		nonNilUuid := uuid.Must(uuid.NewV4())
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		rateAreaId, err := models.FetchRateAreaID(suite.DB(), nilUuid, &nonNilUuid, contract.ID)
		suite.Equal(uuid.Nil, rateAreaId)
		suite.Error(err)
	})
}
