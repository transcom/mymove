package entitlements

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *EntitlementsServiceSuite) TestWeightRestrictor() {
	suite.Run("Successfully apply a weight restriction within max allowance", func() {
		// Create a blank entitlement db entry, nothing fancy we just want to update columns
		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		// Set a weight restriction within allowance
		restrictor := NewWeightRestrictor()
		updatedEntitlement, err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement, 10000, etag.GenerateEtag(entitlement.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(updatedEntitlement)
		suite.NotNil(updatedEntitlement.WeightRestriction)
		suite.Equal(10000, *updatedEntitlement.WeightRestriction)
	})

	suite.Run("Attempt to apply restriction above max allowance, expect an error", func() {
		// Create a blank entitlement db entry, nothing fancy we just want to update columns
		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		// Set an impossible weight restriction
		restrictor := NewWeightRestrictor()
		updatedEntitlement, err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement, 20000, etag.GenerateEtag(entitlement.UpdatedAt))
		suite.Error(err)
		suite.Nil(updatedEntitlement)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("No maxHhgAllowance parameter found returns error", func() {
		param := models.ApplicationParameters{}
		err := suite.DB().
			Where("parameter_name = ?", "maxHhgAllowance").
			First(&param)
		suite.NoError(err)

		err = suite.DB().Destroy(&param)
		suite.NoError(err)

		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}
		suite.MustCreate(&entitlement)

		restrictor := NewWeightRestrictor()
		updatedEntitlement, err := restrictor.ApplyWeightRestrictionToEntitlement(suite.AppContextForTest(), entitlement, 10000, etag.GenerateEtag(entitlement.UpdatedAt))
		suite.Error(err)
		suite.Nil(updatedEntitlement)
		suite.IsType(apperror.QueryError{}, err)
	})

	suite.Run("Successfully remove a weight restriction", func() {

		// Create an entitlement with a restriction already applied
		weightRestriction := 5000
		entitlement := models.Entitlement{
			ID:                uuid.Must(uuid.NewV4()),
			WeightRestriction: &weightRestriction,
		}
		suite.MustCreate(&entitlement)

		restrictor := NewWeightRestrictor()
		updatedEntitlement, err := restrictor.RemoveWeightRestrictionFromEntitlement(suite.AppContextForTest(), entitlement, etag.GenerateEtag(entitlement.UpdatedAt))
		suite.NoError(err)
		suite.NotNil(updatedEntitlement)
		suite.Nil(updatedEntitlement.WeightRestriction)
	})

	suite.Run("Fails on removing a weight restriction for an entitlement that does not exist", func() {

		entitlement := models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		}

		restrictor := NewWeightRestrictor()
		updatedEntitlement, err := restrictor.RemoveWeightRestrictionFromEntitlement(suite.AppContextForTest(), entitlement, etag.GenerateEtag(entitlement.UpdatedAt))
		suite.Error(err)
		suite.Nil(updatedEntitlement)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
