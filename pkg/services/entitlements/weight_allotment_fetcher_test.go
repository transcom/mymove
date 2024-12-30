package entitlements

import (
	"encoding/json"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *EntitlementsServiceSuite) TestGetWeightAllotment() {
	suite.Run("If a weight allotment is fetched by grade, it should be returned", func() {
		fetcher := NewWeightAllotmentFetcher()

		pg := factory.BuildPayGrade(suite.DB(), nil, nil)
		hhgAllowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    pg,
				LinkOnly: true,
			},
		}, nil)

		allotment, err := fetcher.GetWeightAllotment(suite.AppContextForTest(), pg.Grade, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

		suite.NoError(err)
		suite.Equal(hhgAllowance.TotalWeightSelf, allotment.TotalWeightSelf)
		suite.Equal(hhgAllowance.TotalWeightSelfPlusDependents, allotment.TotalWeightSelfPlusDependents)
		suite.Equal(hhgAllowance.ProGearWeight, allotment.ProGearWeight)
		suite.Equal(hhgAllowance.ProGearWeightSpouse, allotment.ProGearWeightSpouse)
		suite.Equal(hhgAllowance.PayGrade.Grade, pg.Grade)
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

		// Build E-5
		e5 := factory.BuildPayGrade(suite.DB(), nil, nil)
		e5Allowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    e5, // Link the pay grade
				LinkOnly: true,
			},
		}, nil)

		// Build E-6
		e6 := factory.BuildPayGrade(suite.DB(), []factory.Customization{
			{
				Model: models.PayGrade{
					Grade: "E-6",
				},
			},
		}, nil)
		e6Allowance := factory.BuildHHGAllowance(suite.DB(), []factory.Customization{
			{
				Model:    e6,
				LinkOnly: true,
			},
		}, nil)

		allotments, err := fetcher.GetAllWeightAllotments(suite.AppContextForTest())
		suite.NoError(err)
		suite.Len(allotments, 2)

		// Check E-5 allotment by its map key
		e5Key := internalmessages.OrderPayGrade(e5.Grade)
		suite.Equal(e5Allowance.TotalWeightSelf, allotments[e5Key].TotalWeightSelf)
		suite.Equal(e5Allowance.TotalWeightSelfPlusDependents, allotments[e5Key].TotalWeightSelfPlusDependents)
		suite.Equal(e5Allowance.ProGearWeight, allotments[e5Key].ProGearWeight)
		suite.Equal(e5Allowance.ProGearWeightSpouse, allotments[e5Key].ProGearWeightSpouse)

		// Check E-6 allotment by its map key
		e6Key := internalmessages.OrderPayGrade(e6.Grade)
		suite.Equal(e6Allowance.TotalWeightSelf, allotments[e6Key].TotalWeightSelf)
		suite.Equal(e6Allowance.TotalWeightSelfPlusDependents, allotments[e6Key].TotalWeightSelfPlusDependents)
		suite.Equal(e6Allowance.ProGearWeight, allotments[e6Key].ProGearWeight)
		suite.Equal(e6Allowance.ProGearWeightSpouse, allotments[e6Key].ProGearWeightSpouse)
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
