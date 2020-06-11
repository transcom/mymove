package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	tosManagementFee = unit.Cents(12303)
	tosCounselingFee = unit.Cents(8327)
)

var tosAvailableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceTaskOrderServices() {
	suite.setupTaskOrderServicesData()
	params := suite.setupTaskOrderServicesParams()

	suite.T().Run("management success using PaymentServiceItemParams", func(t *testing.T) {
		taskOrderServicesPricer, err := NewTaskOrderServicesPricerFromParams(suite.DB(), models.ReServiceCodeMS, params)
		suite.NoError(err)
		priceCents, err := taskOrderServicesPricer.Price()
		suite.NoError(err)
		suite.Equal(tosManagementFee, priceCents)
	})

	suite.T().Run("management success without PaymentServiceItemParams", func(t *testing.T) {
		taskOrderServicesPricer := NewTaskOrderServicesPricer(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeMS, tosAvailableToPrimeAt)
		priceCents, err := taskOrderServicesPricer.Price()
		suite.NoError(err)
		suite.Equal(tosManagementFee, priceCents)
	})

	suite.T().Run("counseling success using PaymentServiceItemParams", func(t *testing.T) {
		taskOrderServicesPricer, err := NewTaskOrderServicesPricerFromParams(suite.DB(), models.ReServiceCodeCS, params)
		suite.NoError(err)
		priceCents, err := taskOrderServicesPricer.Price()
		suite.NoError(err)
		suite.Equal(tosCounselingFee, priceCents)
	})

	suite.T().Run("counseling success without PaymentServiceItemParams", func(t *testing.T) {
		taskOrderServicesPricer := NewTaskOrderServicesPricer(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeCS, tosAvailableToPrimeAt)
		priceCents, err := taskOrderServicesPricer.Price()
		suite.NoError(err)
		suite.Equal(tosCounselingFee, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		emptyParams := models.PaymentServiceItemParams{}
		_, err := NewTaskOrderServicesPricerFromParams(suite.DB(), models.ReServiceCodeMS, emptyParams)
		suite.Error(err)
	})

	suite.T().Run("sending invalid service code", func(t *testing.T) {
		taskOrderServicesPricer := NewTaskOrderServicesPricer(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDLH, tosAvailableToPrimeAt)
		_, err := taskOrderServicesPricer.Price()
		suite.Error(err)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		taskOrderServicesPricer := NewTaskOrderServicesPricer(suite.DB(), "BOGUS", models.ReServiceCodeMS, tosAvailableToPrimeAt)
		_, err := taskOrderServicesPricer.Price()
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupTaskOrderServicesData() {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	managementService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})

	baseManagementTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      managementService.ID,
		PriceCents:     tosManagementFee,
	}
	suite.MustSave(&baseManagementTaskOrderFee)

	counselingService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeCS,
			},
		})

	baseCounselingTaskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     tosCounselingFee,
	}
	suite.MustSave(&baseCounselingTaskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupTaskOrderServicesParams() models.PaymentServiceItemParams {
	var params models.PaymentServiceItemParams

	paramsToCreate := []struct {
		key     models.ServiceItemParamName
		keyType models.ServiceItemParamType
		value   string
	}{
		{
			models.ServiceItemParamNameContractCode,
			models.ServiceItemParamTypeString,
			testdatagen.DefaultContractCode,
		},
		{
			models.ServiceItemParamNameMTOAvailableToPrimeAt,
			models.ServiceItemParamTypeTimestamp,
			tosAvailableToPrimeAt.Format(TimestampParamFormat),
		},
	}

	for _, param := range paramsToCreate {
		paramKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  param.key,
					Type: param.keyType,
				},
			})

		mtoAvailableToPrimeAtParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: paramKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: param.value,
				},
			})
		params = append(params, mtoAvailableToPrimeAtParam)
	}

	return params
}
