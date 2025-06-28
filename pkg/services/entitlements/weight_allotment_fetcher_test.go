package entitlements

import (
	"encoding/json"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *EntitlementsServiceSuite) TestGetWeightAllotment() {
	suite.Run("If a weight allotment is fetched by grade, it should be returned", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), string(models.PaygradeE1), internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

		suite.NoError(err)
		suite.NotEmpty(allotment)
	})

	suite.Run("If pay grade does not exist, return an error", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "X-1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
		suite.Empty(allotment)
	})

	suite.Run("Returns an error if GetMaxGunSafeAllowance returns an error in GetWeightAllotment", func() {
		param := models.ApplicationParameters{}
		err := suite.DB().
			Where("parameter_name = ?", "maxGunSafeAllowance").
			First(&param)
		suite.NoError(err)

		err = suite.DB().Destroy(&param)
		suite.NoError(err)

		fetcher := NewWeightAllotmentFetcher()

		_, err = fetcher.GetWeightAllotment(suite.AppContextForTest(), "E-1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching max gun safe allowance")
	})
}

func (suite *EntitlementsServiceSuite) TestGetAllWeightAllotments() {
	suite.Run("Successfully fetch all weight allotments", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotments, err := fetcher.GetAllWeightAllotments(suite.AppContextForTest())
		suite.NoError(err)
		suite.Greater(len(allotments), 0)
	})

	suite.Run("Returns an error if GetMaxGunSafeAllowance returns an error in GetAllWeightAllotments", func() {
		param := models.ApplicationParameters{}
		err := suite.DB().
			Where("parameter_name = ?", "maxGunSafeAllowance").
			First(&param)
		suite.NoError(err)

		err = suite.DB().Destroy(&param)
		suite.NoError(err)

		fetcher := NewWeightAllotmentFetcher()
		_, err = fetcher.GetWeightAllotment(
			suite.AppContextForTest(),
			"E-1",
			internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "error fetching max gun safe allowance")
	})
}

func (suite *EntitlementsServiceSuite) TestGetWeightAllotmentByOrdersType() {
	suite.Run("Successfully fetch student travel allotment from application_parameters", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotment, err := fetcher.GetWeightAllotment(
			suite.AppContextForTest(),
			"E-1",
			internalmessages.OrdersTypeSTUDENTTRAVEL,
		)

		suite.NoError(err)
		suite.Equal(350, allotment.TotalWeightSelf)
		suite.Equal(350, allotment.TotalWeightSelfPlusDependents)
		suite.Equal(0, allotment.ProGearWeight)
		suite.Equal(0, allotment.ProGearWeightSpouse)
	})

	suite.Run("Returns an error if json does not match allotment from db", func() {
		param := models.ApplicationParameters{}
		err := suite.DB().
			Where("parameter_name = ?", "studentTravelHhgAllowance").
			First(&param)
		suite.NoError(err)

		err = suite.DB().Destroy(&param)
		suite.NoError(err)

		// Proper JSON but not proper target struct
		faultyParamJSON := `{
            "TotalWeight": 350,
            "TotalWeightPlusDependents": 350,
            "ProGear": 0,
            "ProGearSpouse": 0,
            "Allowance": 100
        }`
		rawMessage := json.RawMessage(faultyParamJSON)

		parameter := models.ApplicationParameters{
			ParameterName: models.StringPointer("studentTravelHhgAllowance"),
			ParameterJson: &rawMessage,
		}
		suite.MustCreate(&parameter)

		fetcher := NewWeightAllotmentFetcher()
		_, err = fetcher.GetWeightAllotmentByOrdersType(
			suite.AppContextForTest(),
			internalmessages.OrdersTypeSTUDENTTRAVEL,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "failed to parse weight allotment JSON for orders type")

	})

	suite.Run("Returns an error if no application_parameters entry exists for student travel", func() {
		param := models.ApplicationParameters{}
		err := suite.DB().
			Where("parameter_name = ?", "studentTravelHhgAllowance").
			First(&param)
		suite.NoError(err)

		err = suite.DB().Destroy(&param)
		suite.NoError(err)

		fetcher := NewWeightAllotmentFetcher()
		_, err = fetcher.GetWeightAllotment(
			suite.AppContextForTest(),
			"E-1",
			internalmessages.OrdersTypeSTUDENTTRAVEL,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "failed to fetch weight allotment for orders type STUDENT_TRAVEL: sql: no rows in result set")
	})

	suite.Run("Returns an error if the orders type is not in the ordersTypeToAllotmentAppParamName map", func() {
		fetcher := NewWeightAllotmentFetcher()
		_, err := fetcher.GetWeightAllotmentByOrdersType(
			suite.AppContextForTest(),
			internalmessages.OrdersTypeSEPARATION,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "no entitlement found for orders type SEPARATION")
	})
}
