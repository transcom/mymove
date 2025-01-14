package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestAuthorizedWeightWhenExistsInDB() {
	aw := 3000

	entitlement := models.Entitlement{DBAuthorizedWeight: &aw}
	entitlement.AdminRestrictedWeightLocation = models.BoolPointer(false)
	err := suite.DB().Create(&entitlement)
	suite.NoError(err)

	suite.Equal(entitlement.DBAuthorizedWeight, entitlement.AuthorizedWeight())
}

func (suite *ModelSuite) TestAuthorizedWeightWhenNotInDBAndHaveWeightAllotment() {
	suite.Run("with no dependents authorized, TotalWeightSelf is AuthorizedWeight", func() {
		entitlement := models.Entitlement{}
		entitlement.SetWeightAllotment("E_1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelf, *entitlement.AuthorizedWeight())
	})

	suite.Run("with dependents authorized, TotalWeightSelfPlusDependents is AuthorizedWeight", func() {
		dependentsAuthorized := true
		entitlement := models.Entitlement{DependentsAuthorized: &dependentsAuthorized}
		entitlement.SetWeightAllotment("E_1", internalmessages.OrdersTypePERMANENTCHANGEOFSTATION)

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelfPlusDependents, *entitlement.AuthorizedWeight())
	})
}

func (suite *ModelSuite) TestProGearAndProGearSpouseWeight() {
	suite.Run("no validation errors for ProGearWeight and ProGearSpouseWeight", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       2000,
			ProGearWeightSpouse: 500,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.False(verrs.HasAny(), "Should not have validation errors")
	})

	suite.Run("validation errors for ProGearWeight and ProGearSpouseWeight over max value", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       2001,
			ProGearWeightSpouse: 501,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("pro_gear_weight"))
		suite.NotNil(verrs.Get("pro_gear_weight_spouse"))
	})

	suite.Run("validation errors for ProGearWeight and ProGearSpouseWeight under min value", func() {
		entitlement := models.Entitlement{
			ProGearWeight:       -1,
			ProGearWeightSpouse: -1,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("pro_gear_weight"))
		suite.NotNil(verrs.Get("pro_gear_weight_spouse"))
	})
}

func (suite *ModelSuite) TestOconusFields() {
	suite.Run("no validation errors for valid DependentsUnderTwelve, DependentsTwelveAndOver, and UBAllowance", func() {
		entitlement := models.Entitlement{
			ProGearWeight:           2000,
			ProGearWeightSpouse:     500,
			DependentsUnderTwelve:   models.IntPointer(1),
			DependentsTwelveAndOver: models.IntPointer(2),
			UBAllowance:             models.IntPointer(100),
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.False(verrs.HasAny())
	})

	suite.Run("validation errors for DependentsUnderTwelve and DependentsTwelveAndOver less than 0", func() {
		entitlement := models.Entitlement{
			DependentsTwelveAndOver: models.IntPointer(-1),
			DependentsUnderTwelve:   models.IntPointer(-1),
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("dependents_under_twelve"))
		suite.NotNil(verrs.Get("dependents_twelve_and_over"))
	})
}

func (suite *ModelSuite) TestTotalDependentsCalculation() {
	suite.Run("calculates total dependents correctly when both fields are set", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:         models.IntPointer(2),
			DependentsTwelveAndOver:       models.IntPointer(3),
			AdminRestrictedWeightLocation: models.BoolPointer(false),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(5, *fetchedEntitlement.TotalDependents) // sum of 2 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsUnderTwelve is nil", func() {
		entitlement := models.Entitlement{
			DependentsTwelveAndOver:       models.IntPointer(3),
			AdminRestrictedWeightLocation: models.BoolPointer(false),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Equal(3, *fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(3, *fetchedEntitlement.TotalDependents) // sum of 0 + 3
	})
	suite.Run("calculates total dependents correctly when DependentsTwelveAndOver is nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:         models.IntPointer(2),
			AdminRestrictedWeightLocation: models.BoolPointer(false),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Equal(2, *fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(fetchedEntitlement.TotalDependents)
		suite.Equal(2, *fetchedEntitlement.TotalDependents) // sum of 2 + 0
	})
	suite.Run("sets total dependents to nil when both fields are nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:         nil,
			DependentsTwelveAndOver:       nil,
			AdminRestrictedWeightLocation: models.BoolPointer(false),
		}
		verrs, err := suite.DB().ValidateAndCreate(&entitlement)
		suite.NoError(err)
		suite.False(verrs.HasAny())
		var fetchedEntitlement models.Entitlement
		err = suite.DB().Find(&fetchedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.Nil(fetchedEntitlement.DependentsUnderTwelve)
		suite.Nil(fetchedEntitlement.DependentsTwelveAndOver)
		suite.Nil(fetchedEntitlement.TotalDependents) // NOT 0, NOT A SUM, nil + nil is NULL
	})
}
