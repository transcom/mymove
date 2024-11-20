package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestAuthorizedWeightWhenExistsInDB() {
	aw := 3000
	entitlement := models.Entitlement{DBAuthorizedWeight: &aw}
	err := suite.DB().Create(&entitlement)
	suite.NoError(err)

	suite.Equal(entitlement.DBAuthorizedWeight, entitlement.AuthorizedWeight())
}

func (suite *ModelSuite) TestAuthorizedWeightWhenNotInDBAndHaveWeightAllotment() {
	suite.Run("with no dependents authorized, TotalWeightSelf is AuthorizedWeight", func() {
		entitlement := models.Entitlement{}
		entitlement.SetWeightAllotment("E_1")

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelf, *entitlement.AuthorizedWeight())
	})

	suite.Run("with dependents authorized, TotalWeightSelfPlusDependents is AuthorizedWeight", func() {
		dependentsAuthorized := true
		entitlement := models.Entitlement{DependentsAuthorized: &dependentsAuthorized}
		entitlement.SetWeightAllotment("E_1")

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

func (suite *ModelSuite) TestTotalDependentsValidation() {
	suite.Run("sets sum of total dependents when nil", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: models.IntPointer(3),
			TotalDependents:         nil,
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.NotNil(entitlement.TotalDependents)
		suite.Equal(5, *entitlement.TotalDependents)
	})

	suite.Run("nothing breaks if sum of total dependents is already correct", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: models.IntPointer(3),
			TotalDependents:         models.IntPointer(5),
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.NotNil(entitlement.TotalDependents)
		suite.Equal(5, *entitlement.TotalDependents)
	})

	suite.Run("fixes sum of total dependents if incorrect", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: models.IntPointer(3),
			TotalDependents:         models.IntPointer(6),
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.NotNil(entitlement.TotalDependents)
		suite.Equal(5, *entitlement.TotalDependents)
	})

	suite.Run("handles nil DependentsUnderTwelve", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   nil,
			DependentsTwelveAndOver: models.IntPointer(3),
			TotalDependents:         nil,
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.NotNil(entitlement.TotalDependents)
		suite.Equal(3, *entitlement.TotalDependents)
	})

	suite.Run("handles nil DependentsTwelveAndOver", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   models.IntPointer(2),
			DependentsTwelveAndOver: nil,
			TotalDependents:         nil,
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.NotNil(entitlement.TotalDependents)
		suite.Equal(2, *entitlement.TotalDependents)
	})

	suite.Run("handles nil DependentsUnderTwelve and DependentsTwelveAndOver", func() {
		entitlement := models.Entitlement{
			DependentsUnderTwelve:   nil,
			DependentsTwelveAndOver: nil,
			TotalDependents:         nil,
		}
		verrs, err := entitlement.Validate(suite.DB())
		suite.NoError(err)
		suite.False(verrs.HasAny())
		suite.Nil(entitlement.TotalDependents)
	})
}
