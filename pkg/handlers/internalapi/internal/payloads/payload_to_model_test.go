package payloads

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestPPMShipmentModelWithOptionalDestinationStreet1FromCreate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}

	var pickupAddress internalmessages.Address
	var destinationAddress internalmessages.PPMDestinationAddress

	pickupAddress = internalmessages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = internalmessages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := internalmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model := PPMShipmentModelFromCreate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := internalmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model2 := PPMShipmentModelFromCreate(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 1
	streetAddress1 := "1234 Street"
	destinationAddress.StreetAddress1 = &streetAddress1
	ppmShipmentValidDestAddress1 := internalmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model3 := PPMShipmentModelFromCreate(&ppmShipmentValidDestAddress1)
	suite.Equal(model3.DestinationAddress.StreetAddress1, streetAddress1)
}

func (suite *PayloadsSuite) TestPPMShipmentModelWithOptionalDestinationStreet1FromUpdate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}

	var pickupAddress internalmessages.Address
	var destinationAddress internalmessages.PPMDestinationAddress

	pickupAddress = internalmessages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = internalmessages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := internalmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model := UpdatePPMShipmentModel(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := internalmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model2 := UpdatePPMShipmentModel(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 1
	streetAddress1 := "1234 Street"
	destinationAddress.StreetAddress1 = &streetAddress1
	ppmShipmentValidDestAddress1 := internalmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         &pickupAddress,
		DestinationAddress:    &destinationAddress,
	}

	model3 := UpdatePPMShipmentModel(&ppmShipmentValidDestAddress1)
	suite.Equal(model3.DestinationAddress.StreetAddress1, streetAddress1)
}
