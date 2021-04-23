package models_test

import (
	"testing"

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
	suite.T().Run("with no dependents authorized, TotalWeightSelf is AuthorizedWeight", func(t *testing.T) {
		entitlement := models.Entitlement{}
		entitlement.SetWeightAllotment("E_1")

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelf, *entitlement.AuthorizedWeight())
	})

	suite.T().Run("with dependents authorized, TotalWeightSelfPlusDependents is AuthorizedWeight", func(t *testing.T) {
		dependentsAuthorized := true
		entitlement := models.Entitlement{DependentsAuthorized: &dependentsAuthorized}
		entitlement.SetWeightAllotment("E_1")

		suite.Equal(entitlement.WeightAllotment().TotalWeightSelfPlusDependents, *entitlement.AuthorizedWeight())
	})
}

func (suite *ModelSuite) TestProGearAndProGearSpouseWeight() {
	suite.T().Run("no validation errors for ProGearWeight and ProGearSpouseWeight", func(t *testing.T) {
		entitlement := models.Entitlement{
			ProGearWeight:       2000,
			ProGearWeightSpouse: 500,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.False(verrs.HasAny(), "Should not have validation errors")
	})

	suite.T().Run("validation errors for ProGearWeight and ProGearSpouseWeight over max value", func(t *testing.T) {
		entitlement := models.Entitlement{
			ProGearWeight:       2001,
			ProGearWeightSpouse: 501,
		}
		verrs, _ := entitlement.Validate(suite.DB())
		suite.True(verrs.HasAny())
		suite.NotNil(verrs.Get("pro_gear_weight"))
		suite.NotNil(verrs.Get("pro_gear_weight_spouse"))
	})

	suite.T().Run("validation errors for ProGearWeight and ProGearSpouseWeight under min value", func(t *testing.T) {
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
