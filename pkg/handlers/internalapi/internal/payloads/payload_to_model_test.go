package payloads

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
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

func TestMobileHomeShipmentModelFromCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    *internalmessages.CreateMobileHomeShipment
		expected *models.MobileHome
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "Complete input",
			input: &internalmessages.CreateMobileHomeShipment{
				Make:           models.StringPointer("BrandA"),
				Model:          models.StringPointer("ModelX"),
				Year:           models.Int64Pointer(2024),
				LengthInInches: models.Int64Pointer(60),
				HeightInInches: models.Int64Pointer(13),
				WidthInInches:  models.Int64Pointer(10),
			},
			expected: &models.MobileHome{
				Make:           models.StringPointer("BrandA"),
				Model:          models.StringPointer("ModelX"),
				Year:           models.IntPointer(2024),
				LengthInInches: models.IntPointer(60),
				HeightInInches: models.IntPointer(13),
				WidthInInches:  models.IntPointer(10),
			},
		},
		{
			name: "Partial input with nil values",
			input: &internalmessages.CreateMobileHomeShipment{
				Make:           models.StringPointer("BrandA"),
				Model:          models.StringPointer("ModelX"),
				Year:           nil,
				LengthInInches: models.Int64Pointer(60),
				HeightInInches: nil,
				WidthInInches:  models.Int64Pointer(10),
			},
			expected: &models.MobileHome{
				Make:           models.StringPointer("BrandA"),
				Model:          models.StringPointer("ModelX"),
				Year:           nil,
				LengthInInches: models.IntPointer(60),
				HeightInInches: nil,
				WidthInInches:  models.IntPointer(10),
			},
		},
		{
			name: "All fields are nil",
			input: &internalmessages.CreateMobileHomeShipment{
				Make:           models.StringPointer(""),
				Model:          models.StringPointer(""),
				Year:           nil,
				LengthInInches: nil,
				HeightInInches: nil,
				WidthInInches:  nil,
			},
			expected: &models.MobileHome{
				Make:           models.StringPointer(""),
				Model:          models.StringPointer(""),
				Year:           nil,
				LengthInInches: nil,
				HeightInInches: nil,
				WidthInInches:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MobileHomeShipmentModelFromCreate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
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
