package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCRateEngineServiceSuite) Test_getPaymentServiceItemParam() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeString, testdatagen.DefaultContractCode),
		setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeTimestamp, time.Now().Format(TimestampParamFormat)),
	}

	suite.T().Run("finding expected param", func(t *testing.T) {
		param := getPaymentServiceItemParam(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.NotNil(param)
		suite.Equal(models.ServiceItemParamNameMTOAvailableToPrimeAt, param.ServiceItemParamKey.Key)
	})

	suite.T().Run("param not found", func(t *testing.T) {
		param := getPaymentServiceItemParam(params, models.ServiceItemParamNameWeightEstimated)
		suite.Nil(param)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamInteger() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameDistanceZip5, models.ServiceItemParamTypeInteger, "1234"),
	}

	suite.T().Run("finding expected param value", func(t *testing.T) {
		value, err := getParamInteger(params, models.ServiceItemParamNameDistanceZip5)
		suite.NoError(err)
		suite.Equal(1234, value)
	})

	suite.T().Run("param not found", func(t *testing.T) {
		_, err := getParamInteger(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
		suite.Equal("could not find param with key WeightEstimated", err.Error())
	})

	suite.T().Run("unexpected type", func(t *testing.T) {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeTimestamp, testdatagen.DefaultContractCode),
		}
		_, err := getParamInteger(badParams, models.ServiceItemParamNameContractCode)
		suite.Error(err)
		suite.Equal("trying to convert ContractCode to an integer, but param is of type TIMESTAMP", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamString() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeString, testdatagen.DefaultContractCode),
	}

	suite.T().Run("finding expected param value", func(t *testing.T) {
		value, err := getParamString(params, models.ServiceItemParamNameContractCode)
		suite.NoError(err)
		suite.Equal(testdatagen.DefaultContractCode, value)
	})

	suite.T().Run("param not found", func(t *testing.T) {
		_, err := getParamString(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
	})

	suite.T().Run("unexpected type", func(t *testing.T) {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeTimestamp, testdatagen.DefaultContractCode),
		}
		_, err := getParamString(badParams, models.ServiceItemParamNameContractCode)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamTime() {
	testDate := time.Date(testdatagen.TestYear, time.June, 11, 5, 2, 10, 123, time.UTC)

	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeTimestamp, testDate.Format(TimestampParamFormat)),
		setupParamConvertParam(models.ServiceItemParamNameRequestedPickupDate, models.ServiceItemParamTypeDate, testDate.Format(DateParamFormat)),
	}

	suite.T().Run("finding expected timestamp param value", func(t *testing.T) {
		value, err := getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testDate.Unix(), value.Unix())
		// Note: The current format of time.RFC3339 does not preserve fractions of a second
	})

	suite.T().Run("finding expected date param value", func(t *testing.T) {
		value, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
		suite.NoError(err)
		suite.Equal(testDate.Year(), value.Year())
		suite.Equal(testDate.Month(), value.Month())
		suite.Equal(testDate.Day(), value.Day())
	})

	suite.T().Run("param not found", func(t *testing.T) {
		_, err := getParamTime(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
	})

	suite.T().Run("unexpected type", func(t *testing.T) {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeString, testDate.Format(TimestampParamFormat)),
		}
		_, err := getParamTime(badParams, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.Error(err)
	})
}

func setupParamConvertParam(key models.ServiceItemParamName, keyType models.ServiceItemParamType, value string) models.PaymentServiceItemParam {
	return models.PaymentServiceItemParam{
		Value: value,
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:  key,
			Type: keyType,
		},
	}
}
