package movetaskorder

import (
	"fmt"
	"testing"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_Hide() {
	mtoHider := NewMoveTaskOrderHider(suite.DB())
	suite.T().Run("valid MTO, none to hide", func(t *testing.T) {
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

		mtoAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
			MTOAgent: models.MTOAgent{
				FirstName: swag.String("Peyton"),
				LastName:  swag.String("Wing"),
				Phone:     swag.String("999-999-9999"),
				Email:     swag.String("peyton@example.com"),
			},
		})

		serviceMember := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				FirstName:            swag.String("Gregory"),
				LastName:             swag.String("Van der Heide"),
				Telephone:            swag.String("999-999-9999"),
				SecondaryTelephone:   swag.String("555-123-9999"),
				PersonalEmail:        swag.String("peyton@example.com"),
				ResidentialAddress:   &validAddress1,
				BackupMailingAddress: &validAddress2,
			},
		})

		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				ServiceMember: serviceMember,
			},
		})

		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Order: order,
		})

		testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				PickupAddress:            &validAddress1,
				SecondaryPickupAddress:   &validAddress2,
				DestinationAddress:       &validAddress3,
				SecondaryDeliveryAddress: &validAddress4,
			},
			MTOAgent: mtoAgent,
		})

		result, err := mtoHider.Hide()
		suite.NoError(err)
		suite.Equal(0, len(result))
	})
}

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_isValidFakeModelMTOAgent() {
	validFakeData := []testdatagen.Assertions{
		{MTOAgent: models.MTOAgent{FirstName: swag.String("Peyton")}},
		{MTOAgent: models.MTOAgent{LastName: swag.String("Wing")}},
		{MTOAgent: models.MTOAgent{Phone: swag.String("999-999-9999")}},
		{MTOAgent: models.MTOAgent{Email: swag.String("peyton@example.com")}},
	}
	for idx, validData := range validFakeData {
		suite.T().Run(fmt.Sprintf("valid fake MTOAgent data %d", idx), func(t *testing.T) {
			agent := testdatagen.MakeMTOAgent(suite.DB(), validData)
			result, err := isValidFakeModelMTOAgent(agent)
			suite.NoError(err)
			suite.Equal(true, result)
		})
	}

	badFakeData := []testdatagen.Assertions{
		{MTOAgent: models.MTOAgent{FirstName: swag.String("Billy")}},
		{MTOAgent: models.MTOAgent{LastName: swag.String("Smith")}},
		{MTOAgent: models.MTOAgent{Phone: swag.String("111-111-1111")}},
		{MTOAgent: models.MTOAgent{Email: swag.String("billy@move.mil")}},
	}
	for idx, badData := range badFakeData {
		suite.T().Run(fmt.Sprintf("invalid fake MTOAgent data %d", idx), func(t *testing.T) {
			agent := testdatagen.MakeMTOAgent(suite.DB(), badData)
			result, err := isValidFakeModelMTOAgent(agent)
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
		result, err := isValidFakeModelAddress(&address)
		suite.NoError(err)
		suite.Equal(true, result)
	})

	suite.T().Run("invalid fake address data", func(t *testing.T) {
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				StreetAddress1: "1600 pennsylvania ave",
			},
		})
		result, err := isValidFakeModelAddress(&address)
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
	validFakeData := []testdatagen.Assertions{
		{MTOShipment: models.MTOShipment{PickupAddress: &validPickupAddress}},
		{MTOShipment: models.MTOShipment{SecondaryPickupAddress: &validSecondaryPickupAddress}},
		{MTOShipment: models.MTOShipment{DestinationAddress: &validDestinationAddress}},
		{MTOShipment: models.MTOShipment{SecondaryDeliveryAddress: &validSecondaryDeliveryAddress}},
	}
	for idx, validData := range validFakeData {
		suite.T().Run(fmt.Sprintf("valid fake MTOShipment data %d", idx), func(t *testing.T) {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), validData)
			result, err := isValidFakeModelMTOShipment(shipment)
			suite.NoError(err)
			suite.Equal(true, result)
		})
	}

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
			result, err := isValidFakeModelMTOShipment(shipment)
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
		validFakeData := []testdatagen.Assertions{
			{MTOShipment: models.MTOShipment{PickupAddress: &validPickupAddress}},
			{MTOShipment: models.MTOShipment{SecondaryPickupAddress: &validSecondaryPickupAddress}},
			{MTOShipment: models.MTOShipment{DestinationAddress: &validDestinationAddress}},
			{MTOShipment: models.MTOShipment{SecondaryDeliveryAddress: &validSecondaryDeliveryAddress}},
		}
		var shipments models.MTOShipments
		for _, validData := range validFakeData {
			shipment := testdatagen.MakeMTOShipment(suite.DB(), validData)
			shipments = append(shipments, shipment)
		}

		result, err := isValidFakeModelMTOShipments(shipments)
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
		result, err := isValidFakeModelMTOShipments(shipments)
		suite.NoError(err)
		suite.Equal(false, result)
	})
}
