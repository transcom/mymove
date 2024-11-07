package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
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
		inputAddress := &ghcmessages.Address{
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
		inputAddress := &ghcmessages.Address{
			ID:             strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			StreetAddress1: &streetAddress1,
			City:           &city,
			State:          &state,
			PostalCode:     nil,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(streetAddress1, returnedAddress.StreetAddress1)
		suite.Nil(returnedAddress.StreetAddress2)
		suite.Nil(returnedAddress.StreetAddress3)
		suite.Equal(city, returnedAddress.City)
		suite.Equal(state, returnedAddress.State)
		suite.Equal("", returnedAddress.PostalCode)
	})

	suite.Run("Nil input - returns nil", func() {
		returnedAddress := AddressModel(nil)

		suite.Nil(returnedAddress)
	})

	suite.Run("Blank ID and nil StreetAddress1 - returns nil", func() {
		var blankUUID strfmt.UUID
		inputAddress := &ghcmessages.Address{
			ID:             blankUUID,
			StreetAddress1: nil,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.Nil(returnedAddress)
	})

	suite.Run("Blank ID but valid StreetAddress1 - creates model", func() {
		var blankUUID strfmt.UUID
		inputAddress := &ghcmessages.Address{
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
		input := &ghcmessages.CreateMobileHomeShipment{
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
		input := &ghcmessages.CreateMobileHomeShipment{
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
		input := &ghcmessages.CreateMobileHomeShipment{
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

func (suite *PayloadsSuite) TestPPMShipmentModelWithOptionalDestinationStreet1FromCreate() {
	time := time.Now()
	expectedDepartureDate := handlers.FmtDatePtr(&time)

	address := models.Address{
		StreetAddress1: "some address",
		City:           "city",
		State:          "state",
		PostalCode:     "12345",
	}

	var pickupAddress ghcmessages.Address
	var destinationAddress ghcmessages.PPMDestinationAddress

	pickupAddress = ghcmessages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = ghcmessages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := ghcmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model := PPMShipmentModelFromCreate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(models.PPMShipmentStatusSubmitted, model.Status)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)
	suite.NotNil(model)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := ghcmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model2 := PPMShipmentModelFromCreate(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 1
	streetAddress1 := "123 Street"
	destinationAddress.StreetAddress1 = models.StringPointer(streetAddress1)
	ppmShipmentRealDestinatonAddr1 := ghcmessages.CreatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model3 := PPMShipmentModelFromCreate(&ppmShipmentRealDestinatonAddr1)
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

	var pickupAddress ghcmessages.Address
	var destinationAddress ghcmessages.PPMDestinationAddress

	pickupAddress = ghcmessages.Address{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: &address.StreetAddress1,
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}
	destinationAddress = ghcmessages.PPMDestinationAddress{
		City:           &address.City,
		PostalCode:     &address.PostalCode,
		State:          &address.State,
		StreetAddress1: models.StringPointer(""), // empty string
		StreetAddress2: address.StreetAddress2,
		StreetAddress3: address.StreetAddress3,
	}

	ppmShipment := ghcmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model := PPMShipmentModelFromUpdate(&ppmShipment)

	suite.NotNil(model)
	suite.Equal(model.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test when street address 1 contains white spaces
	destinationAddress.StreetAddress1 = models.StringPointer("  ")
	ppmShipmentWhiteSpaces := ghcmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model2 := PPMShipmentModelFromUpdate(&ppmShipmentWhiteSpaces)
	suite.Equal(model2.DestinationAddress.StreetAddress1, models.STREET_ADDRESS_1_NOT_PROVIDED)

	// test with valid street address 1
	streetAddress1 := "123 Street"
	destinationAddress.StreetAddress1 = models.StringPointer(streetAddress1)
	ppmShipmentRealDestinatonAddr1 := ghcmessages.UpdatePPMShipment{
		ExpectedDepartureDate: expectedDepartureDate,
		PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
		DestinationAddress: struct {
			ghcmessages.PPMDestinationAddress
		}{destinationAddress},
	}

	model3 := PPMShipmentModelFromUpdate(&ppmShipmentRealDestinatonAddr1)
	suite.Equal(model3.DestinationAddress.StreetAddress1, streetAddress1)
}
