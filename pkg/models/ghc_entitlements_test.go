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
