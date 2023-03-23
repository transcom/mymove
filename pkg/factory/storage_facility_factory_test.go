package factory

import (
	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildStorageFacility() {
	suite.Run("Successful creation of default StorageFacility", func() {
		// Under test:      BuildStorageFacility
		// Mocked:          None
		// Set up:          Create a transportation office with no customizations or traits
		// Expected outcome:storageFacility should be created with default values

		// SETUP
		defaultOffice := models.StorageFacility{
			FacilityName: "Storage R Us",
			LotNumber:    models.StringPointer("1234"),
			Phone:        models.StringPointer("555-555-5555"),
			Email:        models.StringPointer("storage@email.com"),
		}
		defaultAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		// CALL FUNCTION UNDER TEST
		storageFacility := BuildStorageFacility(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultOffice.FacilityName, storageFacility.FacilityName)
		suite.Equal(defaultOffice.LotNumber, storageFacility.LotNumber)
		suite.Equal(defaultOffice.Phone, storageFacility.Phone)
		suite.Equal(defaultOffice.Email, storageFacility.Email)

		// Check that address was hooked in
		suite.Equal(defaultAddress.StreetAddress1, storageFacility.Address.StreetAddress1)

	})

	suite.Run("Successful creation of customized StorageFacility", func() {
		// Under test:      BuildStorageFacility
		// Set up:          Create a Storage Facility and pass custom fields
		// Expected outcome:storageFacility should be created with custom fields
		// SETUP
		customStorageFacility := models.StorageFacility{
			ID:           uuid.Must(uuid.NewV4()),
			FacilityName: "Storage R Us",
			LotNumber:    swag.String("1234"),
			Phone:        swag.String("555-555-5555"),
			Email:        swag.String("storage@email.com"),
		}
		customAddress := models.Address{
			StreetAddress1: "987 Another Street",
		}

		// CALL FUNCTION UNDER TEST
		storageFacility := BuildStorageFacility(suite.DB(), []Customization{
			{Model: customStorageFacility},
			{Model: customAddress},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customStorageFacility.ID, storageFacility.ID)
		suite.Equal(customStorageFacility.FacilityName, storageFacility.FacilityName)
		suite.Equal(customStorageFacility.LotNumber, storageFacility.LotNumber)
		suite.Equal(customStorageFacility.Phone, storageFacility.Phone)
		suite.Equal(customStorageFacility.Email, storageFacility.Email)

		// Check that the address was customized
		suite.Equal(customAddress.StreetAddress1, storageFacility.Address.StreetAddress1)
	})

	suite.Run("Successful creation of StorageFacility with trait", func() {
		// Under test:      BuildStorageFacility
		// Mocked:          None
		// Set up:          Create a Storage Facility but pass in a trait that sets
		//                  the address zip to somewhere in the KKFA GBLOC
		// Expected outcome:StorageFacility should have the a zip in KKFA

		storageFacility := BuildStorageFacility(suite.DB(), nil, []Trait{
			GetTraitStorageFacilityKKFA,
		})
		suite.Equal(storageFacility.Address.PostalCode, "85004")
	})

	suite.Run("Successful return of linkOnly StorageFacility", func() {
		// Under test:       BuildStorageFacility
		// Set up:           Pass in a linkOnly storageFacility
		// Expected outcome: No new StorageFacility should be created.

		// Check num StorageFacility records
		precount, err := suite.DB().Count(&models.StorageFacility{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		storageFacility := BuildStorageFacility(suite.DB(), []Customization{
			{
				Model: models.StorageFacility{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.StorageFacility{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, storageFacility.ID)

	})
	suite.Run("Successful return of stubbed StorageFacility", func() {
		// Under test:       BuildStorageFacility
		// Set up:           Pass in a linkOnly storageFacility
		// Expected outcome: No new StorageFacility should be created.

		// Check num StorageFacility records
		precount, err := suite.DB().Count(&models.StorageFacility{})
		suite.NoError(err)

		// Nil passed in as db
		storageFacility := BuildStorageFacility(nil, []Customization{
			{
				Model: models.StorageFacility{
					FacilityName: "A Different Storage Facility",
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.StorageFacility{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal("A Different Storage Facility", storageFacility.FacilityName)

	})
}
