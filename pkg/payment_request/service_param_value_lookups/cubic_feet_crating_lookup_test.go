package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// todo test for multiple crating dimensions?

func (suite *ServiceParamValueLookupsSuite) TestCubicFeetCratingLookup() {
	key := models.ServiceItemParamNameCubicFeetCrating

	suite.T().Run("golden path", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		cratingDimension := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem.ID,
				Type:             models.DimensionTypeCrate,
				Length:           24000,
				Height:           24000,
				Width:            24000,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Time{},
			},
		})
		itemDimension := testdatagen.MakeMTOServiceItemDimension(suite.DB(), testdatagen.Assertions{
			MTOServiceItemDimension: models.MTOServiceItemDimension{
				MTOServiceItemID: mtoServiceItem.ID,
				Type:             models.DimensionTypeItem,
				Length:           12000,
				Height:           12000,
				Width:            12000,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Time{},
			},
		})
		mtoServiceItem.Dimensions = []models.MTOServiceItemDimension{itemDimension, cratingDimension}
		suite.MustSave(&mtoServiceItem)
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		stringValue, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)

		suite.Equal("8", stringValue)
	})

	suite.T().Run("missing dimension should error", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)

		suite.Error(err)
		suite.Contains(err.Error(), "missing crate dimensions")
	})
}
