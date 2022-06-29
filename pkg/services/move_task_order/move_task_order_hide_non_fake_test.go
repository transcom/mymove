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
	// Set up a move with all valid data.
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
	validAddress3 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
		},
	})
	validAddress4 := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "142 E Barrel Hoop Circle #4A",
		},
	})

	serviceMember := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName:              swag.String("Gregory"),
			LastName:               swag.String("Van der Heide"),
			Telephone:              swag.String("999-999-9999"),
			SecondaryTelephone:     swag.String("123-555-9999"),
			PersonalEmail:          swag.String("peyton@example.com"),
			ResidentialAddressID:   &validAddress1.ID,
			ResidentialAddress:     &validAddress1,
			BackupMailingAddressID: &validAddress2.ID,
			BackupMailingAddress:   &validAddress2,
		},
	})

	order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: serviceMember.ID,
			ServiceMember:   serviceMember,
		},
	})

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: order,
	})

	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:                     move,
		PickupAddress:            validAddress1,
		SecondaryPickupAddress:   validAddress2,
		DestinationAddress:       validAddress3,
		SecondaryDeliveryAddress: validAddress4,
	})

	mtoAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOShipment: mtoShipment,
		MTOAgent: models.MTOAgent{
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Peyton"),
			LastName:      swag.String("Wing"),
			Phone:         swag.String("999-999-9999"),
			Email:         swag.String("peyton@example.com"),
		},
	})

	mtoHider := NewMoveTaskOrderHider()

	suite.RunWithPreloadedData("valid MTO, none to hide", func() {
		result, err := mtoHider.Hide(suite.AppContextForTest())
		suite.NoError(err)

		suite.Len(result, 0)
	})

	suite.RunWithPreloadedData("invalid MTO, one to hide", func() {
		// Change an MTO agent name to an invalid name.
		mtoAgent.FirstName = swag.String("Beyonce")
		suite.MustSave(&mtoAgent)

		result, err := mtoHider.Hide(suite.AppContextForTest())
		suite.NoError(err)

		if suite.Len(result, 1) {
			suite.Equal(result[0].MTOID, move.ID)

			// Check the database to make sure the move is truly hidden.
			var savedMove models.Move
			findErr := suite.DB().Find(&savedMove, move.ID)
			suite.NoError(findErr)
			suite.Equal(savedMove.Show, swag.Bool(false))
		}
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeServiceMember() {
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
	invalidReasons := []string{
		"Britney",
		"Spears",
		"415-275-9467",
		"510-607-4545",
		"peyton@gmail.com",
		"24 Main St",
		"123 Not Real Fake Pl",
	}
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
	invalidSm1 := validServiceMemberAssertions
	invalidSm1.ServiceMember.FirstName = swag.String("Britney")
	invalidSm2 := validServiceMemberAssertions
	invalidSm2.ServiceMember.LastName = swag.String("Spears")
	invalidSm3 := validServiceMemberAssertions
	invalidSm3.ServiceMember.Telephone = swag.String("415-275-9467")
	invalidSm4 := validServiceMemberAssertions
	invalidSm4.ServiceMember.SecondaryTelephone = swag.String("510-607-4545")
	invalidSm5 := validServiceMemberAssertions
	invalidSm5.ServiceMember.PersonalEmail = swag.String("peyton@gmail.com")
	invalidSm6 := validServiceMemberAssertions
	invalidSm6.ServiceMember.ResidentialAddress = &invalidAddress1
	invalidSm7 := validServiceMemberAssertions
	invalidSm7.ServiceMember.ResidentialAddress = &invalidAddress2
	invalidFakeData := []testdatagen.Assertions{
		invalidSm1,
		invalidSm2,
		invalidSm3,
		invalidSm4,
		invalidSm5,
		invalidSm6,
		invalidSm7,
	}

	for idx, invalidData := range invalidFakeData {
		suite.RunWithPreloadedData(fmt.Sprintf("invalid fake Service Member data %d", idx), func() {
			sm := testdatagen.MakeServiceMember(suite.DB(), invalidData)
			result, reasons, err := IsValidFakeModelServiceMember(sm)
			suite.NoError(err)
			suite.Equal(false, result)
			toJSONString, _ := json.Marshal(reasons)
			suite.Contains(string(toJSONString), invalidReasons[idx])
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOAgent() {
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

	badFakeData := []testdatagen.Assertions{
		{MTOAgent: models.MTOAgent{FirstName: swag.String("Billy")}},
		{MTOAgent: models.MTOAgent{LastName: swag.String("Smith")}},
		{MTOAgent: models.MTOAgent{Phone: swag.String("111-111-1111")}},
		{MTOAgent: models.MTOAgent{Email: swag.String("billy@move.mil")}},
	}
	for idx, badData := range badFakeData {
		suite.RunWithPreloadedData(fmt.Sprintf("invalid fake MTOAgent data %d", idx), func() {
			agent := testdatagen.MakeMTOAgent(suite.DB(), badData)
			result, err := IsValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelBackupContact() {

	phone := "999-999-9999"
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

	invalidFakeData := []testdatagen.Assertions{
		{BackupContact: models.BackupContact{Name: "Britney"}},
		{BackupContact: models.BackupContact{Email: "Spears"}},
		{BackupContact: models.BackupContact{Phone: swag.String("415-275-9467")}},
	}

	for idx, invalidData := range invalidFakeData {
		suite.RunWithPreloadedData(fmt.Sprintf("invalid fake Backup Contact data %d", idx), func() {

			bc := testdatagen.MakeBackupContact(suite.DB(), invalidData)
			result, err = IsValidFakeModelBackupContact(bc)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelAddress() {
	suite.Run("valid fake address data", func() {
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
			},
		})
		result, err := IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(true, result)
	})

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
	result, reasons, err := IsValidFakeModelMTOShipment(shipment)
	suite.NoError(err)
	suite.Equal(true, result)
	toJSONString, _ := json.Marshal(reasons)
	suite.Equal("{}", string(toJSONString))

	invalidPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "1600 pennsylvania ave",
		},
	})
	invalidSecondaryPickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "20 W 34th St",
		},
	})
	invalidDestinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "86 Pike Pl",
		},
	})
	invalidSecondaryDeliveryAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "4000 Central Florida Blvd",
		},
	})
	var hideReasons = []string{
		"1600 pennsylvania ave",
		"20 W 34th St",
		"86 Pike Pl",
		"4000 Central Florida Blvd",
	}
	validMTOShipmentAssertion := testdatagen.Assertions{
		PickupAddress:            validPickupAddress,
		SecondaryPickupAddress:   validSecondaryPickupAddress,
		DestinationAddress:       validDestinationAddress,
		SecondaryDeliveryAddress: validSecondaryDeliveryAddress,
	}
	invalidMTO1 := validMTOShipmentAssertion
	invalidMTO1.PickupAddress = invalidPickupAddress
	invalidMTO2 := validMTOShipmentAssertion
	invalidMTO2.SecondaryPickupAddress = invalidSecondaryPickupAddress
	invalidMTO3 := validMTOShipmentAssertion
	invalidMTO3.DestinationAddress = invalidDestinationAddress
	invalidMTO4 := validMTOShipmentAssertion
	invalidMTO4.SecondaryDeliveryAddress = invalidSecondaryDeliveryAddress

	invalidFakeData := []testdatagen.Assertions{
		invalidMTO1,
		invalidMTO2,
		invalidMTO3,
		invalidMTO4,
	}
	for idx, invalidData := range invalidFakeData {
		suite.RunWithPreloadedData(fmt.Sprintf("invalid fake MTOShipment data %d", idx), func() {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), invalidData)
			result, reasons, err := IsValidFakeModelMTOShipment(shipment)
			suite.NoError(err)
			suite.Equal(false, result)
			toJSONString, _ := json.Marshal(reasons)
			suite.Contains(string(toJSONString), hideReasons[idx])
		})
	}
}
