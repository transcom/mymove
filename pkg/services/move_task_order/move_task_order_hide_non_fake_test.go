package movetaskorder_test

import (
	"encoding/json"
	"fmt"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/services/move_task_order"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_Hide() {

	setupTestData := func() models.ServiceMember {
		validAddress1 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
				},
			},
		}, nil)
		validAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName:          models.StringPointer("Gregory"),
					LastName:           models.StringPointer("Van der Heide"),
					Telephone:          models.StringPointer("999-999-9999"),
					SecondaryTelephone: models.StringPointer("123-555-9999"),
					PersonalEmail:      models.StringPointer("peyton@example.com"),
				},
			},
			{
				Model:    validAddress1,
				Type:     &factory.Addresses.ResidentialAddress,
				LinkOnly: true,
			},
			{
				Model:    validAddress2,
				Type:     &factory.Addresses.BackupMailingAddress,
				LinkOnly: true,
			},
		}, nil)
		return serviceMember
	}

	mtoHider := m.NewMoveTaskOrderHider()

	suite.Run("valid MTO, none to hide", func() {
		// Under test:       Hide function hides moves that aren't using fake data
		//                   Returns a list with hidden move IDs and reasons
		// Mocked:           None
		// Set up:           Create a move where all the data is valid
		// Expected outcome: No error, no hidden moves returned

		// The basic move created here should use valid fake data
		// Unfortunately testdatagen uses an invalid address so we supply a service member
		// with valid fake address
		serviceMember := setupTestData()
		factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)

		result, err := mtoHider.Hide(suite.AppContextForTest())
		// Expect no error, no hidden moves
		suite.NoError(err)
		suite.Len(result, 0)
	})

	suite.Run("invalid MTO, one to hide", func() {
		// Under test:       Hide function hides moves that aren't using fake data
		//                   Returns a list with hidden move IDs and reasons
		// Set up:           Create a move where one data item is invalid
		// Expected outcome: No error, 1 hidden move returned, show field is disabled in db

		// Make a whole move with an invalid MTO agent name
		serviceMember := setupTestData()
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOAgent{
					FirstName: models.StringPointer("Beyonce"),
				},
			},
		}, nil)
		result, err := mtoHider.Hide(suite.AppContextForTest())
		suite.NoError(err)

		// Expect 1 hidden move
		suite.Len(result, 1)
		suite.Equal(result[0].MTOID, mtoShipment.MoveTaskOrder.ID)

		// Check the database to make sure the move is truly hidden.
		var savedMove models.Move
		findErr := suite.DB().Find(&savedMove, mtoShipment.MoveTaskOrder.ID)
		suite.NoError(findErr)
		suite.Equal(savedMove.Show, models.BoolPointer(false))
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeServiceMember() {

	suite.Run("valid servicemember data", func() {
		// Under test:       IsValidFakeModelServiceMember function
		//                   Returns true/false, the reasons, and err
		// Set up:           Create a servicemember with valid data
		// Expected outcome: Returns true, no reasons
		address1 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
				},
			},
		}, nil)
		address2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		sm := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName:          models.StringPointer("Peyton"),
					LastName:           models.StringPointer("Wing"),
					Telephone:          models.StringPointer("999-999-9999"),
					SecondaryTelephone: models.StringPointer("999-999-9999"),
					PersonalEmail:      models.StringPointer("peyton@example.com"),
				},
			},
			{
				Model:    address1,
				Type:     &factory.Addresses.ResidentialAddress,
				LinkOnly: true,
			},
			{
				Model:    address2,
				Type:     &factory.Addresses.BackupMailingAddress,
				LinkOnly: true,
			}}, nil)
		result, reasons, err := m.IsValidFakeModelServiceMember(sm)
		suite.NoError(err)
		suite.Equal(true, result)
		toJSONString, _ := json.Marshal(reasons)
		suite.Equal("{}", string(toJSONString))
	})

	// These are the expected failures in the following invalid data test
	invalidFields := [][]string{
		{"fullname", "Britney"},
		{"fullname", "Spears"},
		{"phone", "415-275-9467"},
		{"phone2", "510-607-4545"},
		{"email", "peyton@gmail.com"},
		{"residentialaddress", "24 Main St"},
		{"backupmailingaddress", "123 Not Real Fake Pl"},
	}

	setupInvalidTestData := func(index int) (models.ServiceMember, []string) {
		// Create a valid service member, then replace each field with an invalid string and
		// ensure it is caught by the function.
		address1 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
				},
			},
		}, nil)
		address2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		invalidAddress1 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "24 Main St",
				},
			},
		}, nil)
		invalidAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Not Real Fake Pl",
				},
			},
		}, nil)
		validServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Peyton"),
			LastName:           models.StringPointer("Wing"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("999-999-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
		}

		invalidData := validServiceMember
		invalidResAddress := address1
		invalidBackupAddress := address2

		if index == 0 {
			invalidData.FirstName = models.StringPointer("Britney")
		} else if index == 1 {
			invalidData.LastName = models.StringPointer("Spears")
		} else if index == 2 {
			invalidData.Telephone = models.StringPointer("415-275-9467")
		} else if index == 3 {
			invalidData.SecondaryTelephone = models.StringPointer("510-607-4545")
		} else if index == 4 {
			invalidData.PersonalEmail = models.StringPointer("peyton@gmail.com")
		} else if index == 5 {
			invalidResAddress = invalidAddress1
		} else if index == 6 {
			invalidBackupAddress = invalidAddress2
		}

		invalidSm := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{Model: invalidData},
			{
				Model:    invalidResAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.ResidentialAddress,
			},
			{
				Model:    invalidBackupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.BackupMailingAddress,
			},
		}, nil)
		return invalidSm, invalidFields[index]
	}

	// Under test:       IsValidFakeModelServiceMember function
	//                   Returns true/false, the reasons, and err
	// Set up:           Create a set of valid data. Then for each field, create invalid data.
	//                   One at a time change a field from valid to invalid
	// Expected outcome: Returns false, with the invalid field and reason supplied
	for idx := 0; idx < 7; idx++ {
		suite.Run(fmt.Sprintf("invalid fake Service Member %s", invalidFields[idx][0]), func() {
			// Get the invalid service member and expected error
			sm, expectedReason := setupInvalidTestData(idx)
			result, reasons, err := m.IsValidFakeModelServiceMember(sm)

			// Expect no error, false result and the expected reason to match.
			suite.NoError(err)
			suite.Equal(false, result)
			toJSONString, _ := json.Marshal(reasons)
			suite.Contains(string(toJSONString), expectedReason[0])
			suite.Contains(string(toJSONString), expectedReason[1])
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOAgent() {

	badFakeData := []factory.Customization{
		{Model: models.MTOAgent{FirstName: models.StringPointer("Billy")}},
		{Model: models.MTOAgent{LastName: models.StringPointer("Smith")}},
		{Model: models.MTOAgent{Phone: models.StringPointer("111-111-1111")}},
		{Model: models.MTOAgent{Email: models.StringPointer("billy@move.mil")}},
	}

	suite.Run("valid MTOAgent data", func() {
		// Under test:       IsValidFakeModelMTOAgent function checks if mtoagent is valid
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create an agent with valid data
		// Expected outcome: Returns true, no error
		agent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{
			{
				Model: models.MTOAgent{
					FirstName: models.StringPointer("Peyton"),
					LastName:  models.StringPointer("Wing"),
					Phone:     models.StringPointer("999-999-9999"),
					Email:     models.StringPointer("peyton@example.com"),
				},
			},
		}, nil)
		result, err := m.IsValidFakeModelMTOAgent(agent)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	// Under test:       IsValidFakeModelMTOAgent function checks if mtoagent is valid
	//                   Returns true/false, and err if there was a failure
	// Set up:           For each field, create invalid data.
	//                   One at a time, create an agent with one field changed to invalid
	// Expected outcome: Returns false
	for idx, badData := range badFakeData {
		suite.Run(fmt.Sprintf("invalid fake MTOAgent data %d", idx), func() {
			agent := factory.BuildMTOAgent(suite.DB(), []factory.Customization{badData}, nil)
			result, err := m.IsValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelBackupContact() {

	phone := "999-999-9999"

	invalidFakeData := []models.BackupContact{
		{FirstName: "Britney"},
		{LastName: "Blonde"},
		{Email: "Spears"},
		{Phone: "415-275-9467"},
	}

	suite.Run("valid backup contact", func() {
		// Under test:       IsValidFakeModelBackupContact function
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create a contact with valid data
		// Expected outcome: Returns true, no error

		validBackupContact := factory.BuildBackupContact(suite.DB(), []factory.Customization{
			{
				Model: models.BackupContact{
					FirstName: "Robin",
					LastName:  "Fenstermacher",
					Email:     "robin@example.com",
					Phone:     phone,
				},
			},
		}, nil)
		result, err := m.IsValidFakeModelBackupContact(validBackupContact)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	// Under test:       IsValidFakeModelBackupContact function
	//                   Returns true/false, and err if there was a failure
	// Set up:           For each field, create invalid data.
	//                   One at a time, create a contact with one field changed to invalid
	// Expected outcome: Returns false

	for idx, invalidData := range invalidFakeData {
		suite.Run(fmt.Sprintf("invalid fake Backup Contact data %d", idx), func() {
			bc := factory.BuildBackupContact(nil, []factory.Customization{
				{
					Model: invalidData,
				},
			}, nil)
			result, err := m.IsValidFakeModelBackupContact(bc)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelAddress() {
	suite.Run("valid fake address data", func() {
		// Under test:       IsValidFakeModelAddress function
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create an address with valid data
		// Expected outcome: Returns true, no error

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
				},
			},
		}, nil)
		result, err := m.IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	// Under test:       IsValidFakeModelAddress function
	//                   Returns true/false, and err if there was a failure
	// Set up:           For each field, create invalid data.
	//                   One at a time, create an address with one field changed to invalid
	// Expected outcome: Returns false

	suite.Run("invalid fake address data", func() {
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1600 pennsylvania ave",
				},
			},
		}, nil)
		result, err := m.IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(false, result)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipment() {

	setupTestData := func() models.MTOShipment {

		validPickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "7 Q St",
				},
			},
		}, nil)
		validSecondaryPickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		validTertiaryPickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		validDestinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
				},
			},
		}, nil)
		validSecondaryDeliveryAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "142 E Barrel Hoop Circle #4A",
				},
			},
		}, nil)
		validTertiaryDeliveryAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "448 Washington Blvd NE",
				},
			},
		}, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
				},
			},
			{
				Model:    validPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    validSecondaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model:    validTertiaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.TertiaryPickupAddress,
			},
			{
				Model:    validDestinationAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
			{
				Model:    validSecondaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
			{
				Model:    validTertiaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.TertiaryDeliveryAddress,
			},
		}, nil)
		return shipment
	}
	suite.Run("valid shipment data", func() {
		// Under test:       IsValidFakeModelMTOShipment function
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create a shipment with valid data
		// Expected outcome: Returns true, no error

		validShipment := setupTestData()
		result, reasons, err := m.IsValidFakeModelMTOShipment(validShipment)
		suite.NoError(err)
		suite.Equal(true, result)
		toJSONString, _ := json.Marshal(reasons)
		suite.Equal("{}", string(toJSONString))
	})

	var hideReasons = [][]string{
		{"pickupaddress", "1600 pennsylvania ave"},
		{"secondarypickupaddress", "20 W 34th St"},
		{"destinationaddress", "86 Pike Pl"},
		{"secondarydeliveryaddress", "4000 Central Florida Blvd"},
	}

	setupInvalidTestData := func(index int) (models.MTOShipment, []string) {
		// Create valid shipment only to have valid addresses
		validShipment := setupTestData()

		// Create a valid set of customizations
		validCustomizations := []factory.Customization{
			{
				Model:    *validShipment.PickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    *validShipment.SecondaryPickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			},
			{
				Model:    *validShipment.DestinationAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			},
			{
				Model:    *validShipment.SecondaryDeliveryAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			},
		}

		// Based on test index, swap out a valid customization with an invalid one
		var shipment models.MTOShipment
		invalidCustomization := validCustomizations
		if index == 0 {
			// Copy the valid customizations then overwrite the pickup address
			invalidCustomization[0] = factory.Customization{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							StreetAddress1: "1600 pennsylvania ave",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			}
			shipment = factory.BuildMTOShipment(suite.DB(), invalidCustomization, nil)

		} else if index == 1 {
			invalidCustomization[1] = factory.Customization{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							StreetAddress1: "20 W 34th St",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryPickupAddress,
			}
			shipment = factory.BuildMTOShipment(suite.DB(), invalidCustomization, nil)

		} else if index == 2 {
			invalidCustomization[2] = factory.Customization{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							StreetAddress1: "86 Pike Pl",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.DeliveryAddress,
			}
			shipment = factory.BuildMTOShipment(suite.DB(), invalidCustomization, nil)

		} else if index == 3 {
			invalidCustomization[3] = factory.Customization{
				Model: factory.BuildAddress(suite.DB(), []factory.Customization{
					{
						Model: models.Address{
							StreetAddress1: "4000 Central Florida Blvd",
						},
					},
				}, nil),
				LinkOnly: true,
				Type:     &factory.Addresses.SecondaryDeliveryAddress,
			}
			shipment = factory.BuildMTOShipment(suite.DB(), invalidCustomization, nil)
		}
		return shipment, hideReasons[index]
	}
	// Under test:       IsValidFakeModelMTOShipment function
	//                   Returns true/false, reasons, and err if there was a failure
	// Set up:           For each field, create invalid data.
	//                   One at a time, create a shipment with one field changed to invalid
	// Expected outcome: Returns false, with correct reason

	for idx := 0; idx < 4; idx++ {
		suite.Run(fmt.Sprintf("invalid fake MTOShipment %s", hideReasons[idx][0]), func() {
			shipment, expectedReason := setupInvalidTestData(idx)
			result, reasons, err := m.IsValidFakeModelMTOShipment(shipment)
			suite.NoError(err)
			suite.Equal(false, result)
			toJSONString, _ := json.Marshal(reasons)
			suite.Contains(string(toJSONString), expectedReason[0])
			suite.Contains(string(toJSONString), expectedReason[1])
		})
	}
}
