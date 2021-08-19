package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCRateEngineServiceSuite) Test_getPaymentServiceItemParam() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeString, testdatagen.DefaultContractCode),
		setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeTimestamp, time.Now().Format(TimestampParamFormat)),
	}

	suite.Run("finding expected param", func() {
		param := getPaymentServiceItemParam(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.NotNil(param)
		suite.Equal(models.ServiceItemParamNameMTOAvailableToPrimeAt, param.ServiceItemParamKey.Key)
	})

	suite.Run("param not found", func() {
		param := getPaymentServiceItemParam(params, models.ServiceItemParamNameWeightEstimated)
		suite.Nil(param)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamInt() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameDistanceZip5, models.ServiceItemParamTypeInteger, "1234"),
	}

	suite.Run("finding expected param value", func() {
		value, err := getParamInt(params, models.ServiceItemParamNameDistanceZip5)
		suite.NoError(err)
		suite.Equal(1234, value)
	})

	suite.Run("param not found", func() {
		_, err := getParamInt(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
		suite.Equal("could not find param with key WeightEstimated", err.Error())
	})

	suite.Run("unexpected type", func() {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeTimestamp, testdatagen.DefaultContractCode),
		}
		_, err := getParamInt(badParams, models.ServiceItemParamNameContractCode)
		suite.Error(err)
		suite.Equal("trying to convert ContractCode to an int, but param is of type TIMESTAMP", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamFloat() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier, models.ServiceItemParamTypeDecimal, "0.0006255"),
	}

	suite.Run("finding expected param value", func() {
		value, err := getParamFloat(params, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
		suite.NoError(err)
		suite.Equal(0.0006255, value)
	})

	suite.Run("param not found", func() {
		_, err := getParamFloat(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
		suite.Equal("could not find param with key WeightEstimated", err.Error())
	})

	suite.Run("unexpected type", func() {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeTimestamp, testdatagen.DefaultContractCode),
		}
		_, err := getParamFloat(badParams, models.ServiceItemParamNameContractCode)
		suite.Error(err)
		suite.Equal("trying to convert ContractCode to an float, but param is of type TIMESTAMP", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamString() {
	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeString, testdatagen.DefaultContractCode),
	}

	suite.Run("finding expected param value", func() {
		value, err := getParamString(params, models.ServiceItemParamNameContractCode)
		suite.NoError(err)
		suite.Equal(testdatagen.DefaultContractCode, value)
	})

	suite.Run("param not found", func() {
		_, err := getParamString(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
		suite.Equal("could not find param with key WeightEstimated", err.Error())
	})

	suite.Run("unexpected type", func() {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameContractCode, models.ServiceItemParamTypeTimestamp, testdatagen.DefaultContractCode),
		}
		_, err := getParamString(badParams, models.ServiceItemParamNameContractCode)
		suite.Error(err)
		suite.Equal("trying to convert ContractCode to a string, but param is of type TIMESTAMP", err.Error())
	})
}

func (suite *GHCRateEngineServiceSuite) Test_getParamTime() {
	testDate := time.Date(testdatagen.TestYear, time.June, 11, 5, 2, 10, 123, time.UTC)

	params := models.PaymentServiceItemParams{
		setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeTimestamp, testDate.Format(TimestampParamFormat)),
		setupParamConvertParam(models.ServiceItemParamNameRequestedPickupDate, models.ServiceItemParamTypeDate, testDate.Format(DateParamFormat)),
	}

	suite.Run("finding expected timestamp param value", func() {
		value, err := getParamTime(params, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testDate.Unix(), value.Unix())
		// Note: The current format of time.RFC3339 does not preserve fractions of a second
	})

	suite.Run("finding expected date param value", func() {
		value, err := getParamTime(params, models.ServiceItemParamNameRequestedPickupDate)
		suite.NoError(err)
		suite.Equal(testDate.Year(), value.Year())
		suite.Equal(testDate.Month(), value.Month())
		suite.Equal(testDate.Day(), value.Day())
	})

	suite.Run("param not found", func() {
		_, err := getParamTime(params, models.ServiceItemParamNameWeightEstimated)
		suite.Error(err)
		suite.Equal("could not find param with key WeightEstimated", err.Error())
	})

	suite.Run("unexpected type", func() {
		badParams := models.PaymentServiceItemParams{
			setupParamConvertParam(models.ServiceItemParamNameMTOAvailableToPrimeAt, models.ServiceItemParamTypeString, testDate.Format(TimestampParamFormat)),
		}
		_, err := getParamTime(badParams, models.ServiceItemParamNameMTOAvailableToPrimeAt)
		suite.Error(err)
		suite.Equal("trying to convert MTOAvailableToPrimeAt to a time, but param is of type STRING", err.Error())
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
