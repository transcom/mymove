package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildMTOServiceItemDimension() {
	suite.Run("Successful creation of default MTOServiceItemDimension", func() {
		// Under test:      BuildMTOServiceItemDimension
		// Mocked:          None
		// Set up:          Create a service item dimension with no customizations or traits
		// Expected outcome:mtoServiceItemDimension should be created with default values

		// SETUP
		// CALL FUNCTION UNDER TEST
		mTOServiceItemDimension := BuildMTOServiceItemDimension(suite.DB(), nil, nil)

		suite.NotNil(mTOServiceItemDimension.MTOServiceItem)
		suite.NotNil(mTOServiceItemDimension.MTOServiceItemID)
		suite.Equal(models.DimensionTypeItem, mTOServiceItemDimension.Type)
		suite.Equal(unit.ThousandthInches(12000), mTOServiceItemDimension.Length)
		suite.Equal(unit.ThousandthInches(12000), mTOServiceItemDimension.Height)
		suite.Equal(unit.ThousandthInches(12000), mTOServiceItemDimension.Width)
	})

	suite.Run("Successful creation of customized MTOServiceItemDimension", func() {
		// Under test:      BuildMTOServiceItemDimension
		// Mocked:          None
		// Set up:          Create a service item dimension and pass custom fields
		// Expected outcome:mtoServiceItemDimension should be created with custom values

		// SETUP
		customLength := unit.ThousandthInches(16*12*1000 + 1000)
		customHeight := unit.ThousandthInches(8 * 12 * 1000)
		customWidth := unit.ThousandthInches(9 * 12 * 1000)

		customServiceItemDimension := models.MTOServiceItemDimension{
			Type:   models.DimensionTypeCrate,
			Length: customLength,
			Height: customHeight,
			Width:  customWidth,
		}

		// CALL FUNCTION UNDER TEST
		serviceItemDimension := BuildMTOServiceItemDimension(suite.DB(), []Customization{
			{
				Model: customServiceItemDimension,
			},
		}, nil)

		suite.Equal(models.DimensionTypeCrate, serviceItemDimension.Type)
		suite.Equal(customLength, serviceItemDimension.Length)
		suite.Equal(customWidth, serviceItemDimension.Width)
		suite.Equal(customHeight, serviceItemDimension.Height)
	})

	suite.Run("Successful return of linkOnly MTOServiceItemDimension", func() {
		// Under test:       BuildMTOServiceItemDimension
		// Set up:           Pass in a linkOnly mtoServiceItemDimension
		// Expected outcome: No new MTOServiceItemDimension should be created.

		// Check num MTOServiceItemDimension records
		precount, err := suite.DB().Count(&models.MTOServiceItemDimension{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceItemDimension := BuildMTOServiceItemDimension(suite.DB(), []Customization{
			{
				Model: models.MTOServiceItemDimension{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.MTOServiceItemDimension{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceItemDimension.ID)
	})
}
