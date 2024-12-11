package entitlements

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *EntitlementsServiceSuite) TestWeightRestrictor() {

	setupHhgAllowanceParameter := func() {
		parameter := models.ApplicationParameters{
			ParameterName:  models.StringPointer("maxHhgAllowance"),
			ParameterValue: models.StringPointer("18000"),
		}
		suite.MustCreate(&parameter)
	}

	suite.Run("Successfully apply a weight restriction within max allowance", func() {
		setupHhgAllowanceParameter()
		// Create a blank entitlement db entry, nothing fancy we just want to update columns
		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		// Set a weight restriction within allowance
		restrictor := NewWeightRestrictor()
		err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement.ID, 10000)
		suite.NoError(err)

		// Fetch updated entitlement
		var updatedEntitlement models.Entitlement
		err = suite.DB().Find(&updatedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.True(updatedEntitlement.IsWeightRestricted)
		suite.NotNil(updatedEntitlement.WeightRestriction)
		suite.Equal(10000, *updatedEntitlement.WeightRestriction)
	})

	suite.Run("Attempt to apply restriction above max allowance, expect an error", func() {
		setupHhgAllowanceParameter()
		// Create a blank entitlement db entry, nothing fancy we just want to update columns
		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		// Set an impossible weight restriction
		restrictor := NewWeightRestrictor()
		err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement.ID, 20000)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("No maxHhgAllowance parameter found returns error", func() {
		// Ensure no parameter with maxHhgAllowance exists
		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		restrictor := NewWeightRestrictor()
		err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement.ID, 10000)
		suite.Error(err)
		suite.IsType(apperror.QueryError{}, err)
	})

	suite.Run("Successfully remove a weight restriction", func() {
		setupHhgAllowanceParameter()

		// Create an entitlement with a restriction already applied
		weightRestriction := 5000
		entitlement := models.Entitlement{
			ID:                 uuid.Must(uuid.NewV4()),
			IsWeightRestricted: true,
			WeightRestriction:  &weightRestriction,
		}
		suite.MustCreate(&entitlement)

		restrictor := NewWeightRestrictor()

		// Remove the restriction
		err := restrictor.RemoveWeightRestrictionFromEntitlement(suite.AppContextForTest(), entitlement.ID)
		suite.NoError(err)

		// Fetch updated entitlement
		var updatedEntitlement models.Entitlement
		err = suite.DB().Find(&updatedEntitlement, entitlement.ID)
		suite.NoError(err)
		suite.False(updatedEntitlement.IsWeightRestricted)
		suite.Nil(updatedEntitlement.WeightRestriction)
	})

	suite.Run("Fails on removing a weight restriction for an entitlement that does not exist", func() {
		setupHhgAllowanceParameter()

		restrictor := NewWeightRestrictor()

		err := restrictor.RemoveWeightRestrictionFromEntitlement(suite.AppContextForTest(), uuid.Must(uuid.NewV4()))
		suite.Error(err)
	})
}
