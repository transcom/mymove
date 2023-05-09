package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildSITAddressUpdate() {
	suite.Run("Successful creation of default SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create an SITAddressUpdate with no customizations or traits
		// Expected outcome:SITAddressUpdate should be created with default values

		// SETUP
		defaultOldAddress := BuildAddress(suite.DB(), nil, nil)
		defaultNewAddress := BuildAddress(suite.DB(), nil, []Trait{GetTraitAddress2})
		// Create a default SITAddressUpdate to compare values
		defaultSIT := models.SITAddressUpdate{
			Distance: 40,
			Status:   models.SITAddressStatusRequested,
		}

		// FUNCTION UNDER TEST
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultSIT.Distance, sitAddressUpdate.Distance)
		suite.Equal(defaultSIT.Status, sitAddressUpdate.Status)
		suite.Nil(sitAddressUpdate.OfficeRemarks)
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.False(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.False(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.False(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.False(sitAddressUpdate.OldAddressID.IsNil())
		suite.Equal(defaultOldAddress.PostalCode, sitAddressUpdate.OldAddress.PostalCode)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.False(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.False(sitAddressUpdate.NewAddressID.IsNil())
		suite.Equal(defaultNewAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
	})

	suite.Run("Successful creation of customized SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create SITAddressUpdate with customization
		// Expected outcome:SITAddressUpdate should be created with customized values

		// SETUP
		customUpdate := models.SITAddressUpdate{
			ID:                uuid.Must(uuid.NewV4()),
			ContractorRemarks: "custom contractor remarks",
			OfficeRemarks:     models.StringPointer("office remarks"),
			Distance:          40,
			Reason:            "new reason",
			Status:            models.SITAddressStatusRejected,
		}

		customServiceItem := models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusRejected,
		}

		customOldAddress := models.Address{
			ID:         uuid.Must(uuid.NewV4()),
			PostalCode: "77083",
		}

		customNewAddress := models.Address{
			ID:         uuid.Must(uuid.NewV4()),
			PostalCode: "90210",
		}

		// FUNCTION UNDER TEST
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), []Customization{
			{Model: customUpdate},
			{Model: customServiceItem},
			{
				Model: customOldAddress,
				Type:  &Addresses.SITAddressUpdateOldAddress,
			},
			{
				Model: customNewAddress,
				Type:  &Addresses.SITAddressUpdateNewAddress,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customUpdate.ID, sitAddressUpdate.ID)
		suite.Equal(customUpdate.ContractorRemarks, sitAddressUpdate.ContractorRemarks)
		suite.Equal(*customUpdate.OfficeRemarks, *sitAddressUpdate.OfficeRemarks)
		suite.Equal(customUpdate.Distance, sitAddressUpdate.Distance)
		suite.Equal(customUpdate.Reason, sitAddressUpdate.Reason)
		suite.Equal(customUpdate.Status, sitAddressUpdate.Status)
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.False(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.False(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.Equal(customServiceItem.ID, sitAddressUpdate.MTOServiceItem.ID)
		suite.Equal(customServiceItem.Status, sitAddressUpdate.MTOServiceItem.Status)
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.False(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.False(sitAddressUpdate.OldAddressID.IsNil())
		suite.Equal(customOldAddress.ID, sitAddressUpdate.OldAddress.ID)
		suite.Equal(customOldAddress.PostalCode, sitAddressUpdate.OldAddress.PostalCode)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.False(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.False(sitAddressUpdate.NewAddressID.IsNil())
		suite.Equal(customNewAddress.ID, sitAddressUpdate.NewAddress.ID)
		suite.Equal(customNewAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
	})

	suite.Run("Successful return of linkOnly SITAddressUpdate", func() {
		// Under test:       BuildSITAddressUpdate
		// Set up:           Pass in a linkOnly SITAddressUpdate
		// Expected outcome: No new SITAddressUpdate should be created.

		// Check num SITAddressUpdates
		precount, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), []Customization{
			{
				Model: models.SITAddressUpdate{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(id, sitAddressUpdate.ID)

		// Count how many notification are in the DB, no new
		// SITAddressUpdate should have been created
		count, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of stubbed SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Set up:          Create a stubbed SITAddressUpdate
		// Expected outcome:No new SITAddressUpdate should be created

		// Check num SITAddressUpdates
		precount, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)

		sitAddressUpdate := BuildSITAddressUpdate(nil, nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.True(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.True(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.True(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.True(sitAddressUpdate.OldAddressID.IsNil())
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.True(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.True(sitAddressUpdate.NewAddressID.IsNil())

		// Count how many notification are in the DB, no new
		// SITAddressUpdate should have been created
		count, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}
