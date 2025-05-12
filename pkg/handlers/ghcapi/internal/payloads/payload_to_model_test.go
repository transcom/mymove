package payloads

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
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
		PpmType:               ghcmessages.PPMType(models.PPMTypeIncentiveBased),
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
	suite.Equal(model.PPMType, models.PPMTypeIncentiveBased)
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

func (suite *PayloadsSuite) TestVLocationModel() {
	city := "LOS ANGELES"
	state := "CA"
	postalCode := "90210"
	county := "LOS ANGELES"
	usPostRegionCityId := uuid.Must(uuid.NewV4())

	vLocation := &ghcmessages.VLocation{
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

func (suite *PayloadsSuite) TestWeightTicketModelFromUpdate() {
	suite.Run("Success - Complete input", func() {
		emptyWeight := int64(5000)
		fullWeight := int64(8000)
		ownsTrailer := true
		trailerMeetsCriteria := false
		status := ghcmessages.PPMDocumentStatusAPPROVED
		reason := "Valid reason"
		adjustedNetWeight := int64(2900)
		netWeightRemarks := "Adjusted for fuel weight"
		vehicleDescription := "Ford F-150"
		missingEmptyWeightTicket := true
		missingFullWeightTicket := false

		input := &ghcmessages.UpdateWeightTicket{
			EmptyWeight:              &emptyWeight,
			FullWeight:               &fullWeight,
			OwnsTrailer:              ownsTrailer,
			TrailerMeetsCriteria:     trailerMeetsCriteria,
			Status:                   status,
			Reason:                   reason,
			AdjustedNetWeight:        &adjustedNetWeight,
			NetWeightRemarks:         netWeightRemarks,
			VehicleDescription:       &vehicleDescription,
			MissingEmptyWeightTicket: &missingEmptyWeightTicket,
			MissingFullWeightTicket:  &missingFullWeightTicket,
		}

		result := WeightTicketModelFromUpdate(input)

		suite.IsType(&models.WeightTicket{}, result)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&emptyWeight), result.EmptyWeight)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&fullWeight), result.FullWeight)
		suite.Equal(handlers.FmtBool(ownsTrailer), result.OwnsTrailer)
		suite.Equal(handlers.FmtBool(trailerMeetsCriteria), result.TrailerMeetsCriteria)
		suite.Equal(handlers.FmtString(reason), result.Reason)
		suite.Equal((*models.PPMDocumentStatus)(handlers.FmtString(string(status))), result.Status)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&adjustedNetWeight), result.AdjustedNetWeight)
		suite.Equal(handlers.FmtString(netWeightRemarks), result.NetWeightRemarks)
		suite.Equal(handlers.FmtString(vehicleDescription), result.VehicleDescription)
		suite.Equal(handlers.FmtBool(missingEmptyWeightTicket), result.MissingEmptyWeightTicket)
		suite.Equal(handlers.FmtBool(missingFullWeightTicket), result.MissingFullWeightTicket)
	})

	suite.Run("Success - Missing optional fields", func() {
		emptyWeight := int64(5000)
		fullWeight := int64(8000)
		ownsTrailer := true
		trailerMeetsCriteria := false
		status := ghcmessages.PPMDocumentStatusAPPROVED
		reason := "Valid reason"
		adjustedNetWeight := int64(2900)
		netWeightRemarks := "Adjusted for fuel weight"

		input := &ghcmessages.UpdateWeightTicket{
			EmptyWeight:          &emptyWeight,
			FullWeight:           &fullWeight,
			OwnsTrailer:          ownsTrailer,
			TrailerMeetsCriteria: trailerMeetsCriteria,
			Status:               status,
			Reason:               reason,
			AdjustedNetWeight:    &adjustedNetWeight,
			NetWeightRemarks:     netWeightRemarks,
		}

		result := WeightTicketModelFromUpdate(input)

		suite.IsType(&models.WeightTicket{}, result)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&emptyWeight), result.EmptyWeight)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&fullWeight), result.FullWeight)
		suite.Equal(handlers.FmtBool(ownsTrailer), result.OwnsTrailer)
		suite.Equal(handlers.FmtBool(trailerMeetsCriteria), result.TrailerMeetsCriteria)
		suite.Equal(handlers.FmtString(reason), result.Reason)
		suite.Equal((*models.PPMDocumentStatus)(handlers.FmtString(string(status))), result.Status)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&adjustedNetWeight), result.AdjustedNetWeight)
		suite.Equal(handlers.FmtString(netWeightRemarks), result.NetWeightRemarks)

		suite.Nil(result.VehicleDescription)
		suite.Nil(result.MissingEmptyWeightTicket)
		suite.Nil(result.MissingFullWeightTicket)
	})
}

func (suite *PayloadsSuite) TestMovingExpenseModelFromUpdate() {
	suite.Run("Success - Complete input", func() {
		description := "Test moving expense"
		reason := "Just testing"
		trackingNumber := "TRACK123"
		movingExpenseStatus := ghcmessages.PPMDocumentStatusAPPROVED
		amount := int64(1000)
		weightStored := int64(1500)
		weightShipped := int64(1200)
		sitEstimatedCost := int64(2000)
		sitReimburseableAmount := int64(2500)
		sitStartDate := strfmt.Date(time.Now())
		sitEndDate := strfmt.Date(time.Now().Add(24 * time.Hour))
		isProGear := true
		proGearBelongsToSelf := false
		proGearDescription := "Pro gear details"
		paidWithGTCC := true
		missingReceipt := false

		expenseType := ghcmessages.OmittableMovingExpenseTypeSMALLPACKAGE
		sitLocation := ghcmessages.SITLocationTypeORIGIN

		updateMovingExpense := &ghcmessages.UpdateMovingExpense{
			MovingExpenseType:      &expenseType,
			Description:            &description,
			SitLocation:            &sitLocation,
			Amount:                 amount,
			SitStartDate:           sitStartDate,
			SitEndDate:             sitEndDate,
			Status:                 movingExpenseStatus,
			Reason:                 reason,
			WeightStored:           weightStored,
			SitEstimatedCost:       &sitEstimatedCost,
			SitReimburseableAmount: &sitReimburseableAmount,
			TrackingNumber:         &trackingNumber,
			WeightShipped:          &weightShipped,
			IsProGear:              &isProGear,
			ProGearBelongsToSelf:   &proGearBelongsToSelf,
			ProGearDescription:     &proGearDescription,
			PaidWithGTCC:           &paidWithGTCC,
			MissingReceipt:         &missingReceipt,
		}

		result := MovingExpenseModelFromUpdate(updateMovingExpense)
		suite.IsType(&models.MovingExpense{}, result)

		suite.Equal(models.MovingExpenseReceiptTypeSmallPackage, *result.MovingExpenseType, "MovingExpenseType should match")
		suite.Equal(&description, result.Description, "Description should match")
		suite.Equal(models.SITLocationTypeOrigin, *result.SITLocation, "SITLocation should match")
		suite.Equal(handlers.FmtInt64PtrToPopPtr(&amount), result.Amount, "Amount should match")
		suite.Equal(handlers.FmtDatePtrToPopPtr(&sitStartDate), result.SITStartDate, "SITStartDate should match")
		suite.Equal(handlers.FmtDatePtrToPopPtr(&sitEndDate), result.SITEndDate, "SITEndDate should match")
		suite.Equal(models.PPMDocumentStatusApproved, *result.Status, "Status should match")
		suite.Equal(handlers.FmtString(reason), result.Reason, "Reason should match")
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&weightStored), result.WeightStored, "WeightStored should match")
		suite.Equal(handlers.FmtInt64PtrToPopPtr(&sitEstimatedCost), result.SITEstimatedCost, "SITEstimatedCost should match")
		suite.Equal(handlers.FmtInt64PtrToPopPtr(&sitReimburseableAmount), result.SITReimburseableAmount, "SITReimburseableAmount should match")
		suite.Equal(handlers.FmtStringPtr(&trackingNumber), result.TrackingNumber, "TrackingNumber should match")
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&weightShipped), result.WeightShipped, "WeightShipped should match")
		suite.Equal(handlers.FmtBoolPtr(&isProGear), result.IsProGear, "IsProGear should match")
		suite.Equal(handlers.FmtBoolPtr(&proGearBelongsToSelf), result.ProGearBelongsToSelf, "ProGearBelongsToSelf should match")
		suite.Equal(handlers.FmtStringPtr(&proGearDescription), result.ProGearDescription, "ProGearDescription should match")
		suite.Equal(handlers.FmtBoolPtr(&paidWithGTCC), result.PaidWithGTCC, "PaidWithGTCC should match")
		suite.Equal(handlers.FmtBoolPtr(&missingReceipt), result.MissingReceipt, "MissingReceipt should match")
	})
}

func (suite *PayloadsSuite) TestOfficeUserModelFromUpdate() {
	suite.Run("success - complete input", func() {
		telephone := "111-111-1111"

		payload := &ghcmessages.OfficeUserUpdate{
			Telephone: &telephone,
		}

		oldMiddleInitials := "H"
		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:      "John",
					LastName:       "Doe",
					MiddleInitials: &oldMiddleInitials,
					Telephone:      "555-555-5555",
					Email:          "johndoe@example.com",
					Active:         true,
				},
			},
		}, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(*payload.Telephone, returnedOfficeUser.Telephone)
	})

	suite.Run("success - only update Telephone", func() {
		telephone := "111-111-1111"
		payload := &ghcmessages.OfficeUserUpdate{
			Telephone: &telephone,
		}

		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Telephone: "555-555-5555",
				},
			},
		}, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(oldOfficeUser.ID, returnedOfficeUser.ID)
		suite.Equal(oldOfficeUser.UserID, returnedOfficeUser.UserID)
		suite.Equal(oldOfficeUser.Email, returnedOfficeUser.Email)
		suite.Equal(oldOfficeUser.FirstName, returnedOfficeUser.FirstName)
		suite.Equal(oldOfficeUser.MiddleInitials, returnedOfficeUser.MiddleInitials)
		suite.Equal(oldOfficeUser.LastName, returnedOfficeUser.LastName)
		suite.Equal(*payload.Telephone, returnedOfficeUser.Telephone)
		suite.Equal(oldOfficeUser.Active, returnedOfficeUser.Active)
	})

	suite.Run("Fields do not update if payload is empty", func() {
		payload := &ghcmessages.OfficeUserUpdate{}

		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		returnedOfficeUser := OfficeUserModelFromUpdate(payload, &oldOfficeUser)

		suite.NotNil(returnedOfficeUser)
		suite.Equal(oldOfficeUser.ID, returnedOfficeUser.ID)
		suite.Equal(oldOfficeUser.UserID, returnedOfficeUser.UserID)
		suite.Equal(oldOfficeUser.Email, returnedOfficeUser.Email)
		suite.Equal(oldOfficeUser.FirstName, returnedOfficeUser.FirstName)
		suite.Equal(oldOfficeUser.MiddleInitials, returnedOfficeUser.MiddleInitials)
		suite.Equal(oldOfficeUser.LastName, returnedOfficeUser.LastName)
		suite.Equal(oldOfficeUser.Telephone, returnedOfficeUser.Telephone)
		suite.Equal(oldOfficeUser.Active, returnedOfficeUser.Active)
	})

	suite.Run("Error - Return Office User if payload is nil", func() {
		oldOfficeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		returnedUser := OfficeUserModelFromUpdate(nil, &oldOfficeUser)

		suite.Equal(&oldOfficeUser, returnedUser)
	})

	suite.Run("Error - Return nil if Office User is nil", func() {
		telephone := "111-111-1111"
		payload := &ghcmessages.OfficeUserUpdate{
			Telephone: &telephone,
		}
		returnedUser := OfficeUserModelFromUpdate(payload, nil)

		suite.Nil(returnedUser)
	})
}

func (suite *PayloadsSuite) TestProGearWeightTicketModelFromUpdate() {
	suite.Run("Success - Complete input", func() {
		weight := int64(100)
		status := ghcmessages.PPMDocumentStatusAPPROVED
		reason := "Valid reason"
		description := "test description"
		hasWeightTickets, belongsToSelf := true, true

		input := &ghcmessages.UpdateProGearWeightTicket{
			Weight:           &weight,
			HasWeightTickets: hasWeightTickets,
			BelongsToSelf:    belongsToSelf,
			Status:           status,
			Reason:           reason,
			Description:      description,
		}

		result := ProgearWeightTicketModelFromUpdate(input)

		suite.IsType(&models.ProgearWeightTicket{}, result)
		suite.Equal(handlers.PoundPtrFromInt64Ptr(&weight), result.Weight)
		suite.Equal(hasWeightTickets, *result.HasWeightTickets)
		suite.Equal(belongsToSelf, *result.BelongsToSelf)
		suite.Equal(reason, *result.Reason)
		suite.Equal(description, *result.Description)
		suite.Equal((*models.PPMDocumentStatus)(handlers.FmtString(string(status))), result.Status)
	})
}
