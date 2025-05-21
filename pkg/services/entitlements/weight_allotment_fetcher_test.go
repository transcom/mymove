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

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), "E_1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

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
}

func (suite *EntitlementsServiceSuite) TestGetAllWeightAllotments() {
	suite.Run("Successfully fetch all weight allotments", func() {
		fetcher := NewWeightAllotmentFetcher()

		allotments, err := fetcher.GetAllWeightAllotments(suite.AppContextForTest())
		suite.NoError(err)
		suite.Greater(len(allotments), 0)
	})
}

func (suite *EntitlementsServiceSuite) TestGetWeightAllotmentByOrdersType() {
	setupHhgStudentAllowanceParameter := func() {
		paramJSON := `{
            "TotalWeightSelf": 350,
            "TotalWeightSelfPlusDependents": 350,
            "ProGearWeight": 0,
            "ProGearWeightSpouse": 0,
            "UnaccompaniedBaggageAllowance": 100
        }`
		rawMessage := json.RawMessage(paramJSON)

		parameter := models.ApplicationParameters{
			ParameterName: models.StringPointer("studentTravelHhgAllowance"),
			ParameterJson: &rawMessage,
		}
		suite.MustCreate(&parameter)
	}

	suite.Run("Successfully fetch student travel allotment from application_parameters", func() {
		setupHhgStudentAllowanceParameter()
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
		suite.Equal(100, allotment.UnaccompaniedBaggageAllowance)
	})

	suite.Run("Returns an error if json does not match allotment from db", func() {
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
		_, err := fetcher.GetWeightAllotmentByOrdersType(
			suite.AppContextForTest(),
			internalmessages.OrdersTypeSTUDENTTRAVEL,
		)

		suite.Error(err)
		suite.Contains(err.Error(), "failed to parse weight allotment JSON for orders type")

	})

	suite.Run("Returns an error if no application_parameters entry exists for student travel", func() {
		// Don’t create the parameter this time for the student travel
		fetcher := NewWeightAllotmentFetcher()
		_, err := fetcher.GetWeightAllotment(
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
