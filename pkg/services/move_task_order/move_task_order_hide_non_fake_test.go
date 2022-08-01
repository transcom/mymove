package movetaskorder_test

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_Hide() {

	setupTestData := func() models.ServiceMember {
		validAddress1 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "7 Q St",
			},
		})
		validAddress2 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "448 Washington Blvd NE",
			},
		})
		serviceMember := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				FirstName:          swag.String("Gregory"),
				LastName:           swag.String("Van der Heide"),
				Telephone:          swag.String("999-999-9999"),
				SecondaryTelephone: swag.String("123-555-9999"),
				PersonalEmail:      swag.String("peyton@example.com"),

				ResidentialAddressID:   &validAddress1.ID,
				ResidentialAddress:     &validAddress1,
				BackupMailingAddressID: &validAddress2.ID,
				BackupMailingAddress:   &validAddress2,
			},
		})
		return serviceMember
	}

	mtoHider := NewMoveTaskOrderHider()

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
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			ServiceMember: serviceMember,
			Order: models.Order{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})

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
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			ServiceMember: serviceMember,
			Order: models.Order{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})
		testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
			MTOAgent: models.MTOAgent{
				FirstName: swag.String("Beyonce"),
			},
		})

		result, err := mtoHider.Hide(suite.AppContextForTest())
		suite.NoError(err)

		// Expect 1 hidden move
		suite.Len(result, 1)
		suite.Equal(result[0].MTOID, mtoShipment.MoveTaskOrder.ID)

		// Check the database to make sure the move is truly hidden.
		var savedMove models.Move
		findErr := suite.DB().Find(&savedMove, mtoShipment.MoveTaskOrder.ID)
		suite.NoError(findErr)
		suite.Equal(savedMove.Show, swag.Bool(false))
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeServiceMember() {

	suite.Run("valid servicemember data", func() {
		// Under test:       IsValidFakeModelServiceMember function
		//                   Returns true/false, the reasons, and err
		// Set up:           Create a servicemember with valid data
		// Expected outcome: Returns true, no reasons
		address1 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "7 Q St",
			},
		})
		address2 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "448 Washington Blvd NE",
			},
		})
		sm := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				FirstName:            swag.String("Peyton"),
				LastName:             swag.String("Wing"),
				Telephone:            swag.String("999-999-9999"),
				SecondaryTelephone:   swag.String("999-999-9999"),
				PersonalEmail:        swag.String("peyton@example.com"),
				ResidentialAddress:   &address1,
				BackupMailingAddress: &address2,
			}},
		)
		result, reasons, err := IsValidFakeModelServiceMember(sm)
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
		address1 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "7 Q St",
			},
		})
		address2 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "448 Washington Blvd NE",
			},
		})
		invalidAddress1 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "24 Main St",
			},
		})
		invalidAddress2 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "123 Not Real Fake Pl",
			},
		})
		validServiceMemberAssertions := testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				FirstName:            swag.String("Peyton"),
				LastName:             swag.String("Wing"),
				Telephone:            swag.String("999-999-9999"),
				SecondaryTelephone:   swag.String("999-999-9999"),
				PersonalEmail:        swag.String("peyton@example.com"),
				ResidentialAddress:   &address1,
				BackupMailingAddress: &address2,
			}}
		invalidData := validServiceMemberAssertions
		if index == 0 {
			invalidData.ServiceMember.FirstName = swag.String("Britney")
		} else if index == 1 {
			invalidData.ServiceMember.LastName = swag.String("Spears")
		} else if index == 2 {
			invalidData.ServiceMember.Telephone = swag.String("415-275-9467")
		} else if index == 3 {
			invalidData.ServiceMember.SecondaryTelephone = swag.String("510-607-4545")
		} else if index == 4 {
			invalidData.ServiceMember.PersonalEmail = swag.String("peyton@gmail.com")
		} else if index == 5 {
			invalidData.ServiceMember.ResidentialAddress = &invalidAddress1
		} else if index == 6 {
			invalidData.ServiceMember.BackupMailingAddress = &invalidAddress2
		}

		invalidSm := testdatagen.MakeServiceMember(suite.DB(), invalidData)
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
			result, reasons, err := IsValidFakeModelServiceMember(sm)

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

	badFakeData := []testdatagen.Assertions{
		{MTOAgent: models.MTOAgent{FirstName: swag.String("Billy")}},
		{MTOAgent: models.MTOAgent{LastName: swag.String("Smith")}},
		{MTOAgent: models.MTOAgent{Phone: swag.String("111-111-1111")}},
		{MTOAgent: models.MTOAgent{Email: swag.String("billy@move.mil")}},
	}

	suite.Run("valid MTOAgent data", func() {
		// Under test:       IsValidFakeModelMTOAgent function checks if mtoagent is valid
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create an agent with valid data
		// Expected outcome: Returns true, no error
		agent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				FirstName: swag.String("Peyton"),
				LastName:  swag.String("Wing"),
				Phone:     swag.String("999-999-9999"),
				Email:     swag.String("peyton@example.com"),
			},
		})
		result, err := IsValidFakeModelMTOAgent(agent)
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
			agent := testdatagen.MakeMTOAgent(suite.DB(), badData)
			result, err := IsValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelBackupContact() {

	phone := "999-999-9999"

	invalidFakeData := []testdatagen.Assertions{
		{BackupContact: models.BackupContact{Name: "Britney"}},
		{BackupContact: models.BackupContact{Email: "Spears"}},
		{BackupContact: models.BackupContact{Phone: swag.String("415-275-9467")}},
	}

	suite.Run("valid backup contact", func() {
		// Under test:       IsValidFakeModelBackupContact function
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create a contact with valid data
		// Expected outcome: Returns true, no error

		validBackupContact := testdatagen.MakeBackupContact(suite.DB(), testdatagen.Assertions{
			BackupContact: models.BackupContact{
				Name:  "Robin Fenstermacher",
				Email: "robin@example.com",
				Phone: &phone,
			},
		})
		result, err := IsValidFakeModelBackupContact(validBackupContact)
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
			bc := testdatagen.MakeBackupContact(suite.DB(), invalidData)
			result, err := IsValidFakeModelBackupContact(bc)
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

		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
			},
		})
		result, err := IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	// Under test:       IsValidFakeModelAddress function
	//                   Returns true/false, and err if there was a failure
	// Set up:           For each field, create invalid data.
	//                   One at a time, create an address with one field changed to invalid
	// Expected outcome: Returns false

	suite.Run("invalid fake address data", func() {
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "1600 pennsylvania ave",
			},
		})
		result, err := IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(false, result)
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipment() {

	setupTestData := func() models.MTOShipment {

		validPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "7 Q St",
			},
		})
		validSecondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "448 Washington Blvd NE",
			},
		})
		validDestinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
			},
		})
		validSecondaryDeliveryAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "142 E Barrel Hoop Circle #4A",
			},
		})
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
			PickupAddress:            validPickupAddress,
			SecondaryPickupAddress:   validSecondaryPickupAddress,
			DestinationAddress:       validDestinationAddress,
			SecondaryDeliveryAddress: validSecondaryDeliveryAddress,
		})
		return shipment
	}
	suite.Run("valid shipment data", func() {
		// Under test:       IsValidFakeModelMTOShipment function
		//                   Returns true/false, and err if there was a failure
		// Set up:           Create a shipment with valid data
		// Expected outcome: Returns true, no error

		validShipment := setupTestData()
		result, reasons, err := IsValidFakeModelMTOShipment(validShipment)
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

		// Create a valid set of assertions
		validMTOShipmentAssertion := testdatagen.Assertions{
			PickupAddress:            *validShipment.PickupAddress,
			SecondaryPickupAddress:   *validShipment.SecondaryPickupAddress,
			DestinationAddress:       *validShipment.DestinationAddress,
			SecondaryDeliveryAddress: *validShipment.SecondaryDeliveryAddress,
		}

		// Based on test index, swap out a valid assertion with an invalid one
		var shipment models.MTOShipment
		invalidAssertion := validMTOShipmentAssertion
		if index == 0 {
			invalidPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
				Address: models.Address{
					StreetAddress1: "1600 pennsylvania ave",
				},
			})
			// Copy the valid assertions then overwrite the pickup address
			invalidAssertion.PickupAddress = invalidPickupAddress
			shipment = testdatagen.MakeMTOShipment(suite.DB(), invalidAssertion)

		} else if index == 1 {
			invalidSecondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
				Address: models.Address{
					StreetAddress1: "20 W 34th St",
				},
			})
			invalidAssertion.SecondaryPickupAddress = invalidSecondaryPickupAddress
			shipment = testdatagen.MakeMTOShipment(suite.DB(), invalidAssertion)

		} else if index == 2 {

			invalidDestinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
				Address: models.Address{
					StreetAddress1: "86 Pike Pl",
				},
			})
			invalidAssertion.DestinationAddress = invalidDestinationAddress
			shipment = testdatagen.MakeMTOShipment(suite.DB(), invalidAssertion)

		} else if index == 3 {

			invalidSecondaryDeliveryAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
				Address: models.Address{
					StreetAddress1: "4000 Central Florida Blvd",
				},
			})
			invalidAssertion.SecondaryDeliveryAddress = invalidSecondaryDeliveryAddress
			shipment = testdatagen.MakeMTOShipment(suite.DB(), invalidAssertion)
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
			result, reasons, err := IsValidFakeModelMTOShipment(shipment)
			suite.NoError(err)
			suite.Equal(false, result)
			toJSONString, _ := json.Marshal(reasons)
			suite.Contains(string(toJSONString), expectedReason[0])
			suite.Contains(string(toJSONString), expectedReason[1])
		})
	}
}
