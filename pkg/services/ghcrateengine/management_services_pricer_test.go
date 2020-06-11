package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	msPriceCents = unit.Cents(12303)
)

var msAvailiableToPrimeAt = time.Date(testdatagen.TestYear, time.June, 3, 12, 57, 33, 123, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceManagementServices() {
	suite.setupManagementServicesData()
	params := suite.setupManagementServicesParams()
	counselingServicesPricer := NewManagementServicesPricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := counselingServicesPricer.PriceUsingParams(params)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := counselingServicesPricer.Price(testdatagen.DefaultContractCode, msAvailiableToPrimeAt)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, err := counselingServicesPricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.T().Run("not finding a rate record", func(t *testing.T) {
		_, err := counselingServicesPricer.Price("BOGUS", msAvailiableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesData() {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	counselingService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})

	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     msPriceCents,
	}
	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupManagementServicesParams() models.PaymentServiceItemParams {
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
			msAvailiableToPrimeAt.Format(TimestampParamFormat),
		},
	}

	for _, param := range paramsToCreate {
		serviceItemParamKey := testdatagen.MakeServiceItemParamKey(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: models.ServiceItemParamKey{
					Key:  param.key,
					Type: param.keyType,
				},
			})

		serviceItemParam := testdatagen.MakePaymentServiceItemParam(suite.DB(),
			testdatagen.Assertions{
				ServiceItemParamKey: serviceItemParamKey,
				PaymentServiceItemParam: models.PaymentServiceItemParam{
					Value: param.value,
				},
			})
		params = append(params, serviceItemParam)
	}

	return params
}
