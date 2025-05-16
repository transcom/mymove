package entitlements

import (
	"encoding/json"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
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

func (suite *EntitlementsServiceSuite) TestGetTotalWeightAllotment() {
	suite.Run("returns total weight for office user with dependents authorized", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
		})

		entitlement := models.Entitlement{
			DependentsAuthorized: models.BoolPointer(true),
			GunSafeWeight:        100,
		}

		order := models.Order{
			HasDependents: true,
			Grade:         internalmessages.OrderPayGradeE1.Pointer(),
			OrdersType:    internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		}

		fetcher := NewWeightAllotmentFetcher()
		totalWeight, err := fetcher.GetTotalWeightAllotment(appCtx, order, entitlement)

		suite.NoError(err)
		// E-1 PCS = 8000 with dependents
		expected := 8000 + 100
		suite.Equal(expected, totalWeight)
	})

	suite.Run("returns total weight for mil app user with dependents and dependentsAuthorized", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
		})

		entitlement := models.Entitlement{
			DependentsAuthorized: models.BoolPointer(true),
			GunSafeWeight:        50,
		}

		order := models.Order{
			HasDependents: true,
			Grade:         internalmessages.OrderPayGradeE1.Pointer(),
			OrdersType:    internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		}

		fetcher := NewWeightAllotmentFetcher()
		totalWeight, err := fetcher.GetTotalWeightAllotment(appCtx, order, entitlement)

		suite.NoError(err)
		expected := 8000 + 50
		suite.Equal(expected, totalWeight)
	})

	suite.Run("uses self-only weight if dependents not authorized", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
		})

		entitlement := models.Entitlement{
			DependentsAuthorized: models.BoolPointer(false),
			GunSafeWeight:        20,
		}

		order := models.Order{
			HasDependents: true,
			Grade:         internalmessages.OrderPayGradeE1.Pointer(),
			OrdersType:    internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		}

		fetcher := NewWeightAllotmentFetcher()
		totalWeight, err := fetcher.GetTotalWeightAllotment(appCtx, order, entitlement)

		suite.NoError(err)
		expected := 5000 + 20 // E-1 PCS self-only weight + gun safe
		suite.Equal(expected, totalWeight)
	})

	suite.Run("uses self weight with dependent if dependents authorized", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
		})

		entitlement := models.Entitlement{
			DependentsAuthorized: models.BoolPointer(true),
			GunSafeWeight:        20,
		}

		order := models.Order{
			HasDependents: false,
			Grade:         internalmessages.OrderPayGradeE1.Pointer(),
			OrdersType:    internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		}

		fetcher := NewWeightAllotmentFetcher()
		totalWeight, err := fetcher.GetTotalWeightAllotment(appCtx, order, entitlement)

		suite.NoError(err)
		expected := 5000 + 20 // E-1 PCS self-only weight + gun safe
		suite.Equal(expected, totalWeight)
	})

	suite.Run("returns error if weight allotment fetch fails", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
		})

		entitlement := models.Entitlement{
			DependentsAuthorized: models.BoolPointer(true),
			GunSafeWeight:        0,
		}

		order := models.Order{
			HasDependents: true,
			Grade:         nil,
			OrdersType:    internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		}

		fetcher := NewWeightAllotmentFetcher()
		_, err := fetcher.GetTotalWeightAllotment(appCtx, order, entitlement)

		suite.Error(err)
	})
}
