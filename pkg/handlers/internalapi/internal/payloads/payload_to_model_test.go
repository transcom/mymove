package payloads

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestPPMShipmentModelFromUpdate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)
	estimatedWeight := int64(5000)
	proGearWeight := int64(500)
	spouseProGearWeight := int64(50)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
		Country:        &models.Country{Country: "US"},
	}
	address2 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "11111",
	}
	address3 := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "54321",
	}

	var pickupAddress internalmessages.Address
	var secondaryPickupAddress internalmessages.Address
	var tertiaryPickupAddress internalmessages.Address
	var destinationAddress internalmessages.PPMDestinationAddress
	var secondaryDestinationAddress internalmessages.Address
	var tertiaryDestinationAddress internalmessages.Address

	pickupAddress = internalmessages.Address{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = internalmessages.PPMDestinationAddress{
		City:           &address.City,
		Country:        &address.Country.Country,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	secondaryPickupAddress = internalmessages.Address{
		City:           &address2.City,
		Country:        &address.Country.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	secondaryDestinationAddress = internalmessages.Address{
		City:           &address2.City,
		Country:        &address.Country.Country,
		PostalCode:     &address2.PostalCode,
		State:          &address2.State,
		StreetAddress1: &address2.StreetAddress1,
		StreetAddress2: address2.StreetAddress2,
		StreetAddress3: address2.StreetAddress3,
	}
	tertiaryPickupAddress = internalmessages.Address{
		City:           &address3.City,
		Country:        &address.Country.Country,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}
	tertiaryDestinationAddress = internalmessages.Address{
		City:           &address3.City,
		Country:        &address.Country.Country,
		PostalCode:     &address3.PostalCode,
		State:          &address3.State,
		StreetAddress1: &address3.StreetAddress1,
		StreetAddress2: address3.StreetAddress2,
		StreetAddress3: address3.StreetAddress3,
	}

	ppmShipment := internalmessages.UpdatePPMShipment{
		ExpectedDepartureDate:        expectedDepartureDate,
		PickupAddress:                &pickupAddress,
		SecondaryPickupAddress:       &secondaryPickupAddress,
		TertiaryPickupAddress:        &tertiaryPickupAddress,
		DestinationAddress:           &destinationAddress,
		SecondaryDestinationAddress:  &secondaryDestinationAddress,
		TertiaryDestinationAddress:   &tertiaryDestinationAddress,
		SitExpected:                  models.BoolPointer(true),
		EstimatedWeight:              &estimatedWeight,
		HasProGear:                   models.BoolPointer(true),
		ProGearWeight:                &proGearWeight,
		SpouseProGearWeight:          &spouseProGearWeight,
		IsActualExpenseReimbursement: models.BoolPointer(true),
	}

	model := UpdatePPMShipmentModel(&ppmShipment)

	suite.NotNil(model)
	suite.True(*model.SITExpected)
	suite.Equal(unit.Pound(estimatedWeight), *model.EstimatedWeight)
	suite.True(*model.HasProGear)
	suite.Equal(unit.Pound(proGearWeight), *model.ProGearWeight)
	suite.Equal(unit.Pound(spouseProGearWeight), *model.SpouseProGearWeight)
	suite.Nil(model.HasSecondaryPickupAddress)
	suite.Nil(model.HasSecondaryDestinationAddress)
	suite.Nil(model.HasTertiaryPickupAddress)
	suite.Nil(model.HasTertiaryDestinationAddress)
	suite.True(*model.IsActualExpenseReimbursement)
	suite.NotNil(model)
}

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
