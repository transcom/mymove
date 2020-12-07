package movetaskorder_test

import (
	"fmt"
	"testing"

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

	mtoHider := NewMoveTaskOrderHider(suite.DB())

	suite.T().Run("valid MTO, none to hide", func(t *testing.T) {
		result, err := mtoHider.Hide()
		suite.NoError(err)

		suite.Len(result, 0)
	})

	suite.T().Run("invalid MTO, one to hide", func(t *testing.T) {
		// Change an MTO agent name to an invalid name.
		mtoAgent.FirstName = swag.String("Beyonce")
		suite.MustSave(&mtoAgent)

		result, err := mtoHider.Hide()
		suite.NoError(err)

		if suite.Len(result, 1) {
			suite.Equal(result[0].ID, move.ID)
			suite.Equal(result[0].Show, swag.Bool(false))

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
	result, err := IsValidFakeModelServiceMember(sm)
	suite.NoError(err)
	suite.Equal(true, result)

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
	invalidFakeData := []testdatagen.Assertions{
		{ServiceMember: models.ServiceMember{FirstName: swag.String("Britney")}},
		{ServiceMember: models.ServiceMember{LastName: swag.String("Spears")}},
		{ServiceMember: models.ServiceMember{Telephone: swag.String("415-275-9467")}},
		{ServiceMember: models.ServiceMember{SecondaryTelephone: swag.String("510-607-4545")}},
		{ServiceMember: models.ServiceMember{PersonalEmail: swag.String("peyton@gmail.com")}},
		{ServiceMember: models.ServiceMember{ResidentialAddress: &invalidAddress1}},
		{ServiceMember: models.ServiceMember{BackupMailingAddress: &invalidAddress2}},
	}
	for idx, invalidData := range invalidFakeData {
		suite.T().Run(fmt.Sprintf("invalid fake Service Member data %d", idx), func(t *testing.T) {
			sm := testdatagen.MakeServiceMember(suite.DB(), invalidData)
			result, err := IsValidFakeModelServiceMember(sm)
			suite.NoError(err)
			suite.Equal(false, result)
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
		suite.T().Run(fmt.Sprintf("invalid fake MTOAgent data %d", idx), func(t *testing.T) {
			agent := testdatagen.MakeMTOAgent(suite.DB(), badData)
			result, err := IsValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelAddress() {
	suite.T().Run("valid fake address data", func(t *testing.T) {
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "3373 NW Martin Luther King Jr Blvd",
			},
		})
		result, err := IsValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	suite.T().Run("invalid fake address data", func(t *testing.T) {
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
			PickupAddress:            &validPickupAddress,
			SecondaryPickupAddress:   &validSecondaryPickupAddress,
			DestinationAddress:       &validDestinationAddress,
			SecondaryDeliveryAddress: &validSecondaryDeliveryAddress,
		},
	})
	result, err := IsValidFakeModelMTOShipment(shipment)
	suite.NoError(err)
	suite.Equal(true, result)

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
	invalidFakeData := []testdatagen.Assertions{
		{MTOShipment: models.MTOShipment{PickupAddress: &invalidPickupAddress}},
		{MTOShipment: models.MTOShipment{SecondaryPickupAddress: &invalidSecondaryPickupAddress}},
		{MTOShipment: models.MTOShipment{DestinationAddress: &invalidDestinationAddress}},
		{MTOShipment: models.MTOShipment{SecondaryDeliveryAddress: &invalidSecondaryDeliveryAddress}},
	}
	for idx, invalidData := range invalidFakeData {
		suite.T().Run(fmt.Sprintf("invalid fake MTOShipment data %d", idx), func(t *testing.T) {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), invalidData)
			result, err := IsValidFakeModelMTOShipment(shipment)
			suite.NoError(err)
			suite.Equal(false, result)
		})
	}
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOShipments() {
	suite.T().Run("valid fake MTOShipments data", func(t *testing.T) {
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

		var shipments models.MTOShipments

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PickupAddress:            &validPickupAddress,
				SecondaryPickupAddress:   &validSecondaryPickupAddress,
				DestinationAddress:       &validDestinationAddress,
				SecondaryDeliveryAddress: &validSecondaryDeliveryAddress,
			},
		})
		shipments = append(shipments, shipment)
		result, err := IsValidFakeModelMTOShipments(shipments)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	suite.T().Run("invalid fake MTOShipments data", func(t *testing.T) {
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
		invalidFakeData := []testdatagen.Assertions{
			{MTOShipment: models.MTOShipment{PickupAddress: &invalidPickupAddress}},
			{MTOShipment: models.MTOShipment{SecondaryPickupAddress: &invalidSecondaryPickupAddress}},
			{MTOShipment: models.MTOShipment{DestinationAddress: &invalidDestinationAddress}},
			{MTOShipment: models.MTOShipment{SecondaryDeliveryAddress: &invalidSecondaryDeliveryAddress}},
		}
		var shipments models.MTOShipments
		for _, invalidData := range invalidFakeData {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), invalidData)
			shipments = append(shipments, shipment)
		}
		result, err := IsValidFakeModelMTOShipments(shipments)
		suite.NoError(err)
		suite.Equal(false, result)
	})
}
