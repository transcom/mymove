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
	gunSafeWeight := int64(321)

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
		PpmType:                      internalmessages.PPMType(models.PPMTypeActualExpense),
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
		HasGunSafe:                   models.BoolPointer(true),
		GunSafeWeight:                &gunSafeWeight,
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
	suite.Equal(model.PPMType, models.PPMTypeActualExpense)
	suite.True(*model.HasGunSafe)
	suite.Equal(unit.Pound(gunSafeWeight), *model.GunSafeWeight)
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

func (suite *PayloadsSuite) TestVLocationModel() {
	city := "LOS ANGELES"
	state := "CA"
	postalCode := "90210"
	county := "LOS ANGELES"
	usPostRegionCityId := uuid.Must(uuid.NewV4())

	vLocation := &internalmessages.VLocation{
		City:                 city,
		State:                state,
		PostalCode:           postalCode,
		County:               &county,
		UsPostRegionCitiesID: strfmt.UUID(usPostRegionCityId.String()),
	}

	payload := VLocationModel(vLocation)

	suite.IsType(payload, &models.VLocation{})
	suite.Equal(usPostRegionCityId.String(), payload.UsPostRegionCitiesID.String(), "Expected UsPostRegionCitiesID to match")
	suite.Equal(city, payload.CityName, "Expected City to match")
	suite.Equal(state, payload.StateName, "Expected State to match")
	suite.Equal(postalCode, payload.UsprZipID, "Expected PostalCode to match")
	suite.Equal(county, payload.UsprcCountyNm, "Expected County to match")
}

func (suite *PayloadsSuite) TestMovingExpenseModelFromUpdate() {
	suite.Run("Success - Complete input", func() {
		description := "Test moving expense"
		trackingNumber := "TRACK123"
		amount := int64(1000)
		weightStored := int64(1500)
		weightShipped := int64(1200)
		sitReimburseableAmount := int64(2500)
		sitStartDate := strfmt.Date(time.Now())
		sitEndDate := strfmt.Date(time.Now().Add(24 * time.Hour))
		isProGear := true
		proGearBelongsToSelf := false
		proGearDescription := "Pro gear details"

		expenseType := internalmessages.MovingExpenseTypeSMALLPACKAGE
		sitLocation := internalmessages.SITLocationTypeORIGIN

		updateMovingExpense := &internalmessages.UpdateMovingExpense{
			MovingExpenseType:      &expenseType,
			Description:            &description,
			SitLocation:            &sitLocation,
			Amount:                 &amount,
			SitStartDate:           sitStartDate,
			SitEndDate:             sitEndDate,
			WeightStored:           weightStored,
			SitReimburseableAmount: &sitReimburseableAmount,
			TrackingNumber:         &trackingNumber,
			WeightShipped:          weightShipped,
			IsProGear:              &isProGear,
			ProGearBelongsToSelf:   &proGearBelongsToSelf,
			ProGearDescription:     proGearDescription,
		}

		result := MovingExpenseModelFromUpdate(updateMovingExpense)
		suite.IsType(&models.MovingExpense{}, result)

		suite.Equal(models.MovingExpenseReceiptTypeSmallPackage, *result.MovingExpenseType, "MovingExpenseType should match")
		suite.Equal(&description, result.Description, "Description should match")
		suite.Equal(models.SITLocationTypeOrigin, *result.SITLocation, "SITLocation should match")
		suite.Equal(handlers.FmtInt64PtrToPopPtr(&amount), result.Amount, "Amount should match")
		suite.Equal(handlers.FmtDatePtrToPopPtr(&sitStartDate), result.SITStartDate, "SITStartDate should match")
		suite.Equal(handlers.FmtDatePtrToPopPtr(&sitEndDate), result.SITEndDate, "SITEndDate should match")
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&weightStored), result.WeightStored, "WeightStored should match")
		suite.Equal(handlers.FmtInt64PtrToPopPtr(&sitReimburseableAmount), result.SITReimburseableAmount, "SITReimburseableAmount should match")
		suite.Equal(handlers.FmtStringPtr(&trackingNumber), result.TrackingNumber, "TrackingNumber should match")
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&weightShipped), result.WeightShipped, "WeightShipped should match")
		suite.Equal(handlers.FmtBoolPtr(&isProGear), result.IsProGear, "IsProGear should match")
		suite.Equal(handlers.FmtBoolPtr(&proGearBelongsToSelf), result.ProGearBelongsToSelf, "ProGearBelongsToSelf should match")
		suite.Equal(handlers.FmtStringPtr(&proGearDescription), result.ProGearDescription, "ProGearDescription should match")
	})
}

func (suite *PayloadsSuite) TestVIntlLocationModel() {
	city := "LONDON"
	principalDivision := "CARDIFF"
	intlCityCountriesId := uuid.Must(uuid.NewV4())

	vIntlLocation := &internalmessages.VIntlLocation{
		City:                city,
		PrincipalDivision:   principalDivision,
		IntlCityCountriesID: strfmt.UUID(intlCityCountriesId.String()),
	}

	payload := VIntlLocationModel(vIntlLocation)

	suite.IsType(payload, &models.VIntlLocation{})
	suite.Equal(intlCityCountriesId.String(), payload.IntlCityCountriesID.String(), "Expected IntlCityCountriesID to match")
	suite.Equal(city, *payload.CityName, "Expected City to match")
	suite.Equal(principalDivision, *payload.CountryPrnDivName, "Expected Principal Division to match")
}
