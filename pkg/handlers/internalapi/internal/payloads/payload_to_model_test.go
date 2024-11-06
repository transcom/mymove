package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PayloadsSuite) TestAddressModel() {
	streetAddress1 := "123 Main St"
	streetAddress2 := "Apt 4B"
	streetAddress3 := "Building 5"
	city := "New York"
	state := "NY"
	postalCode := "10001"

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalCode,
	}

	suite.Run("Success - Complete input", func() {
		inputAddress := &internalmessages.Address{
			ID:             strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			StreetAddress1: &streetAddress1,
			StreetAddress2: &streetAddress2,
			StreetAddress3: &streetAddress3,
			City:           &city,
			State:          &state,
			PostalCode:     &postalCode,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(expectedAddress.StreetAddress1, returnedAddress.StreetAddress1)
		suite.Equal(expectedAddress.StreetAddress2, returnedAddress.StreetAddress2)
		suite.Equal(expectedAddress.StreetAddress3, returnedAddress.StreetAddress3)
		suite.Equal(expectedAddress.City, returnedAddress.City)
		suite.Equal(expectedAddress.State, returnedAddress.State)
		suite.Equal(expectedAddress.PostalCode, returnedAddress.PostalCode)
	})

	suite.Run("Success - Partial input", func() {
		inputAddress := &internalmessages.Address{
			ID:             strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			StreetAddress1: &streetAddress1,
			City:           &city,
			State:          &state,
			PostalCode:     &postalCode,
			Country:        nil,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(streetAddress1, returnedAddress.StreetAddress1)
		suite.Nil(returnedAddress.StreetAddress2)
		suite.Nil(returnedAddress.StreetAddress3)
		suite.Equal(city, returnedAddress.City)
		suite.Equal(state, returnedAddress.State)
		suite.Equal(postalCode, returnedAddress.PostalCode)
		suite.Nil(returnedAddress.Country)
	})

	suite.Run("Nil input - returns nil", func() {
		returnedAddress := AddressModel(nil)

		suite.Nil(returnedAddress)
	})

	suite.Run("Blank ID and nil StreetAddress1 - returns nil", func() {
		var blankUUID strfmt.UUID
		inputAddress := &internalmessages.Address{
			ID:             blankUUID,
			StreetAddress1: nil,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.Nil(returnedAddress)
	})

	suite.Run("Blank ID but valid StreetAddress1 - creates model", func() {
		var blankUUID strfmt.UUID
		inputAddress := &internalmessages.Address{
			ID:             blankUUID,
			StreetAddress1: &streetAddress1,
			City:           &city,
			State:          &state,
			PostalCode:     &postalCode,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(streetAddress1, returnedAddress.StreetAddress1)
		suite.Equal(city, returnedAddress.City)
		suite.Equal(state, returnedAddress.State)
		suite.Equal(postalCode, returnedAddress.PostalCode)
	})
}

func (suite *PayloadsSuite) TestMobileHomeShipmentModelFromCreate() {
	make := "BrandA"
	model := "ModelX"
	year := int64(2024)
	lengthInInches := int64(60)
	heightInInches := int64(13)
	widthInInches := int64(10)

	expectedMobileHome := models.MobileHome{
		Make:           models.StringPointer(make),
		Model:          models.StringPointer(model),
		Year:           models.IntPointer(int(year)),
		LengthInInches: models.IntPointer(int(lengthInInches)),
		HeightInInches: models.IntPointer(int(heightInInches)),
		WidthInInches:  models.IntPointer(int(widthInInches)),
	}

	suite.Run("Success - Complete input", func() {
		input := &internalmessages.CreateMobileHomeShipment{
			Make:           models.StringPointer(make),
			Model:          models.StringPointer(model),
			Year:           &year,
			LengthInInches: &lengthInInches,
			HeightInInches: &heightInInches,
			WidthInInches:  &widthInInches,
		}

		returnedMobileHome := MobileHomeShipmentModelFromCreate(input)

		suite.IsType(&models.MobileHome{}, returnedMobileHome)
		suite.Equal(expectedMobileHome.Make, returnedMobileHome.Make)
		suite.Equal(expectedMobileHome.Model, returnedMobileHome.Model)
		suite.Equal(expectedMobileHome.Year, returnedMobileHome.Year)
		suite.Equal(expectedMobileHome.LengthInInches, returnedMobileHome.LengthInInches)
		suite.Equal(expectedMobileHome.HeightInInches, returnedMobileHome.HeightInInches)
		suite.Equal(expectedMobileHome.WidthInInches, returnedMobileHome.WidthInInches)
	})

	suite.Run("Success - Partial input", func() {
		input := &internalmessages.CreateMobileHomeShipment{
			Make:           models.StringPointer(make),
			Model:          models.StringPointer(model),
			Year:           nil,
			LengthInInches: &lengthInInches,
			HeightInInches: nil,
			WidthInInches:  &widthInInches,
		}

		returnedMobileHome := MobileHomeShipmentModelFromCreate(input)

		suite.IsType(&models.MobileHome{}, returnedMobileHome)
		suite.Equal(make, *returnedMobileHome.Make)
		suite.Equal(model, *returnedMobileHome.Model)
		suite.Nil(returnedMobileHome.Year)
		suite.Equal(int(lengthInInches), *returnedMobileHome.LengthInInches)
		suite.Nil(returnedMobileHome.HeightInInches)
		suite.Equal(int(widthInInches), *returnedMobileHome.WidthInInches)
	})

	suite.Run("Nil input - returns nil", func() {
		returnedMobileHome := MobileHomeShipmentModelFromCreate(nil)

		suite.Nil(returnedMobileHome)
	})

	suite.Run("All fields are nil - returns empty MobileHome", func() {
		input := &internalmessages.CreateMobileHomeShipment{
			Make:           models.StringPointer(""),
			Model:          models.StringPointer(""),
			Year:           nil,
			LengthInInches: nil,
			HeightInInches: nil,
			WidthInInches:  nil,
		}

		returnedMobileHome := MobileHomeShipmentModelFromCreate(input)

		suite.IsType(&models.MobileHome{}, returnedMobileHome)
		suite.Equal("", *returnedMobileHome.Make)
		suite.Equal("", *returnedMobileHome.Model)
		suite.Nil(returnedMobileHome.Year)
		suite.Nil(returnedMobileHome.LengthInInches)
		suite.Nil(returnedMobileHome.HeightInInches)
		suite.Nil(returnedMobileHome.WidthInInches)
	})
}

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
