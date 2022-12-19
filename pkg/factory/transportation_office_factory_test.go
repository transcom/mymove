package factory

import (
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *FactorySuite) TestBuildTransportationOffice() {
	suite.Run("Successful creation of default transportation office", func() {
		// Under test:      BuildTransportationOffice
		// Mocked:          None
		// Set up:          Create a transportation office with no customizations or traits
		// Expected outcome:transportationOffice should be created with default values
		defaultName := "JPPSO Testy McTest"

		transportationOffice := BuildTransportationOffice(suite.DB(), nil, nil)
		suite.Equal(defaultName, transportationOffice.Name)
	})

	suite.Run("Successful creation of transportationOffice with customization", func() {
		// Under test:      BuildTransportationOffice
		// Set up:          Create a Transportation Office and pass custom name
		// Expected outcome:transportationOffice should be created with custom name
		customName := "Ft Example"
		customGbloc := "TEST"
		customID := uuid.Must(uuid.NewV4())

		transportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
			{
				Model: models.TransportationOffice{
					Name:  customName,
					Gbloc: customGbloc,
					ID:    customID,
				},
			},
		}, nil)
		suite.Equal(customName, transportationOffice.Name)
		suite.Equal(customGbloc, transportationOffice.Gbloc)
		suite.Equal(customID, transportationOffice.ID)
	})
}

func (suite *FactorySuite) TestBuildTransportationOfficeExtra() {
	suite.Run("Successful creation of transportationOffice with linked address", func() {
		// Under test:       BuildTransportationOffice
		// Set up:           Create a Transportation Office and pass in a precreated address
		// Expected outcome: transportationOffice should link in the precreated address and shouldn't create a new address

		address := testdatagen.MakeDefaultAddress(suite.DB())

		// Count how many addresses we have
		precount, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)

		transportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
			{
				Model:    address,
				LinkOnly: true,
			},
		}, nil)

		// Check that the linked address was used
		suite.Equal(address.ID, transportationOffice.AddressID)
		suite.Equal(address.ID, transportationOffice.Address.ID)
		suite.Equal(address.StreetAddress1, transportationOffice.Address.StreetAddress1)
		suite.Equal(address.PostalCode, transportationOffice.Address.PostalCode)

		// VALIDATION
		// Check that no new address was created
		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	//suite.Run("Successful creation of transportationOffice with forced id address", func() {
	//	// Under test:       BuildTransportationOffice
	//	// Set up:           Create a transportationOffice and pass in an ID for address
	//	// Expected outcome: transportationOffice and address should be created
	//	//                   address should have specified ID
	//
	//	uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
	//
	//	// Check that id cannot be found in DB
	//	foundAddress := models.Address{}
	//	err := suite.DB().Find(&foundAddress, uuid)
	//	suite.Error(err)
	//
	//	transportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
	//		{
	//			Model: models.Address{
	//				ID: uuid,
	//			},
	//		},
	//	}, nil)
	//
	//	// Check that the forced ID was used
	//	suite.Equal(uuid, transportationOffice.AddressID)
	//	suite.Equal(uuid, transportationOffice.Address.ID)
	//
	//	// Check that id can be found in DB
	//	foundAddress = models.Address{}
	//	err = suite.DB().Find(&foundAddress, uuid)
	//	suite.NoError(err)
	//
	//})

	suite.Run("Successful creation of transportationOffice with ShippingOffice", func() {
		// Under test:       BuildTransportationOffice
		// Set up:           Create a Transportation Office and with a linked ShippingOffice
		// Expected outcome: new transportationOffice with a ShippingOffice (total of 2 transportationOffices created)
		//					 ShippingOffice is a pointer to a TransportationOffice

		shippingOffice := BuildTransportationOffice(suite.DB(), []Customization{}, nil)

		// Count how many Transportation Offices we have
		// ShippingOffice is a pointer to a Transportation Office
		precount, err := suite.DB().Count(&models.TransportationOffices{})
		suite.NoError(err)

		transportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
			{
				Model: shippingOffice,
				Type:  &ShippingOffice,
			},
		}, nil)

		// Check that the linked ShippingOffice was used
		suite.Equal(shippingOffice.ID, transportationOffice.AddressID)
		suite.Equal(shippingOffice.ID, transportationOffice.Address.ID)
		suite.Equal(shippingOffice.Name, transportationOffice.ShippingOffice.Name)
		suite.Equal(shippingOffice.Gbloc, transportationOffice.ShippingOffice.Gbloc)

		// VALIDATION
		// Check that 2 new transportationOffices were created
		count, err := suite.DB().Count(&models.TransportationOffice{})
		suite.NoError(err)
		suite.Equal(precount+2, count)
	})

	suite.Run("Successful creation of transportationOffice with linked ShippingOffice", func() {
		// Under test:       BuildTransportationOffice
		// Set up:           Create a Transportation Office and with a linked ShippingOffice
		// Expected outcome: transportationOffice should link in the precreated ShippingOffice
		//                   and shouldn't create a new ShippingOffice

		shippingOffice := BuildTransportationOffice(suite.DB(), []Customization{}, nil)

		// Count how many Transportation Offices we have
		// ShippingOffice is a pointer to a Transportation Office
		precount, err := suite.DB().Count(&models.TransportationOffices{})
		suite.NoError(err)

		transportationOffice := BuildTransportationOffice(suite.DB(), []Customization{
			{
				Model:    shippingOffice,
				LinkOnly: true,
			},
		}, nil)

		// Check that the linked ShippingOffice was used
		suite.Equal(shippingOffice.ID, transportationOffice.AddressID)
		suite.Equal(shippingOffice.ID, transportationOffice.Address.ID)
		suite.Equal(shippingOffice.Name, transportationOffice.ShippingOffice.Name)
		suite.Equal(shippingOffice.Gbloc, transportationOffice.ShippingOffice.Gbloc)

		// VALIDATION
		// Check that no new addresses were created
		count, err := suite.DB().Count(&models.Address{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	//suite.Run("Successful creation of stubbed transportationOffice with forced id address", func() {
	//	// Under test:       BuildTransportationOffice
	//	// Set up:           Create a transportationOffice and pass in a precreated user
	//	// Expected outcome: transportationOffice and User should be created with specified emails
	//	uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
	//	transportationOffice := BuildTransportationOffice(nil, []Customization{
	//		{
	//			Model: models.User{
	//				ID: uuid,
	//			},
	//		},
	//	}, nil)
	//	// Check that the forced ID was used
	//	suite.Equal(uuid, transportationOffice.AddressID)
	//	suite.Equal(uuid, transportationOffice.Address.ID)
	//
	//	// Check that id cannot be found in DB
	//	foundUser := models.User{}
	//	err := suite.DB().Find(&foundUser, uuid)
	//	suite.Error(err)
	//
	//	// Check that email was applied to user
	//	suite.Equal(transportationOffice.Email, transportationOffice.Address.StreetAddress1)
	//})
}
