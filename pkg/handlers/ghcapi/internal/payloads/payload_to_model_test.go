package payloads

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestAddressModel() {
	streetAddress1 := "123 Main St"
	streetAddress2 := "Apt 4B"
	streetAddress3 := "Building 5"
	city := "New York"
	state := "NY"
	postalCode := "10001"
	country := "USA"

	expectedAddress := models.Address{
		StreetAddress1: streetAddress1,
		StreetAddress2: &streetAddress2,
		StreetAddress3: &streetAddress3,
		City:           city,
		State:          state,
		PostalCode:     postalCode,
		Country:        &country,
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
			Country:        &country,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(expectedAddress.StreetAddress1, returnedAddress.StreetAddress1)
		suite.Equal(expectedAddress.StreetAddress2, returnedAddress.StreetAddress2)
		suite.Equal(expectedAddress.StreetAddress3, returnedAddress.StreetAddress3)
		suite.Equal(expectedAddress.City, returnedAddress.City)
		suite.Equal(expectedAddress.State, returnedAddress.State)
		suite.Equal(expectedAddress.PostalCode, returnedAddress.PostalCode)
		suite.Equal(expectedAddress.Country, returnedAddress.Country)
	})

	suite.Run("Success - Partial input", func() {
		inputAddress := &ghcmessages.Address{
			ID:             strfmt.UUID(uuid.Must(uuid.NewV4()).String()),
			StreetAddress1: &streetAddress1,
			City:           &city,
			State:          &state,
			PostalCode:     nil,
			Country:        nil,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(streetAddress1, returnedAddress.StreetAddress1)
		suite.Nil(returnedAddress.StreetAddress2)
		suite.Nil(returnedAddress.StreetAddress3)
		suite.Equal(city, returnedAddress.City)
		suite.Equal(state, returnedAddress.State)
		suite.Equal("", returnedAddress.PostalCode)
		suite.Nil(returnedAddress.Country)
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
			Country:        &country,
		}

		returnedAddress := AddressModel(inputAddress)

		suite.IsType(&models.Address{}, returnedAddress)
		suite.Equal(streetAddress1, returnedAddress.StreetAddress1)
		suite.Equal(city, returnedAddress.City)
		suite.Equal(state, returnedAddress.State)
		suite.Equal(postalCode, returnedAddress.PostalCode)
		suite.Equal(country, *returnedAddress.Country)
	})
}

func (suite *PayloadsSuite) TestMobileHomeShipmentModelFromCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    *ghcmessages.CreateMobileHomeShipment
		expected *models.MobileHome
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "Complete input",
			input: &ghcmessages.CreateMobileHomeShipment{
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
			input: &ghcmessages.CreateMobileHomeShipment{
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
			input: &ghcmessages.CreateMobileHomeShipment{
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
